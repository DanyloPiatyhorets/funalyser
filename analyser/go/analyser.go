package analyser

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
)

type Analyser interface {
	Visit(node ast.Node, functionContext *FunctionContext)
}

type TimeAndSpaceComplexityAnalyser struct {
}

func Analyse(filePath string, functionName string) ([]FunctionInfo, error) {
	src, IOError := os.ReadFile(filePath)
	if IOError != nil {
		return nil, IOError
	}
	fset := token.NewFileSet()
	file, IOError := parser.ParseFile(fset, "", src, parser.AllErrors)
	if IOError != nil {
		return nil, IOError
	}
	// ast.Print(fset, file)

	var funcsInfo []FunctionInfo
	var fileContext FileContext = GetFileContext(file)

	var functionInfoError error
	if functionName == "" {
		funcsInfo, functionInfoError = getAllFunctionInfo(file, &fileContext)
	} else {
		var funcInfo FunctionInfo
		funcInfo, functionInfoError = getFunctionInfoByName(file, &fileContext, functionName)
		if functionInfoError == nil {
			funcsInfo = append(funcsInfo, funcInfo)
		}
	}

	return funcsInfo, functionInfoError
}

func getAllFunctionInfo(file *ast.File, fileContext *FileContext) ([]FunctionInfo, error) {
	var funcs []FunctionInfo

	for _, declaration := range file.Decls {
		switch decl := declaration.(type) {
		case *ast.FuncDecl:
			functionContext := GetFunctionContext(decl, fileContext)
			analyser := &TimeAndSpaceComplexityAnalyser{}
			for _, stmt := range decl.Body.List {
				analyser.Visit(stmt, functionContext)
			}
			funcs = append(funcs, ParseContextToInfo(functionContext))
		}
	}

	return funcs, nil
}

func getFunctionInfoByName(file *ast.File, fileContext *FileContext, functionName string) (FunctionInfo, error) {
	for _, declaration := range file.Decls {
		switch decl := declaration.(type) {
		case *ast.FuncDecl:
			if isFunctionName(decl, functionName) {
				functionContext := GetFunctionContext(decl, fileContext)
				analyser := &TimeAndSpaceComplexityAnalyser{}
				for _, stmt := range decl.Body.List {
					analyser.Visit(stmt, functionContext)
				}
				return ParseContextToInfo(functionContext), nil
			}
		}
	}
	return FunctionInfo{}, errors.New("no such function in this file")
}

func (tscAnalyser *TimeAndSpaceComplexityAnalyser) Visit(node ast.Node, functionContext *FunctionContext) {

	switch stmt := node.(type) {
	case *ast.AssignStmt:
		tscAnalyser.Visit(stmt.Rhs[0], functionContext)

	case *ast.BinaryExpr:
		tscAnalyser.Visit(stmt.X, functionContext)
		tscAnalyser.Visit(stmt.Y, functionContext)

	case *ast.BlockStmt:
		for _, inner := range stmt.List {
			tscAnalyser.Visit(inner, functionContext)
		}

	case *ast.CallExpr:
		funIdent, ok := stmt.Fun.(*ast.Ident)
		if !ok {
			return
		}

		switch funIdent.Name {
		case "make":
			switch stmt.Args[0].(type) {
			case *ast.ArrayType:
				if size, ok := stmt.Args[1].(*ast.Ident); ok {
					if IsParam(size.Name, &functionContext.SymbolTable) {
						functionContext.CurrentMalloc = 1 + functionContext.CurrentDepth
					}
				}
			case *ast.MapType:
				functionContext.CurrentMalloc = 1 + functionContext.CurrentDepth
			}

		case "append":
			functionContext.CurrentMalloc = functionContext.CurrentDepth

		case functionContext.Name:
			time, space := GetRecursiveComplexity(stmt)
			functionContext.CurrentDepth += time
			functionContext.CurrentMalloc += space
			functionContext.RecursiveFanOut++
		}

	case *ast.CaseClause:
		for _, inner := range stmt.Body {
			tscAnalyser.Visit(inner, functionContext)
		}

	case *ast.ExprStmt:
		tscAnalyser.Visit(stmt.X, functionContext)

	case *ast.ForStmt:
		condExpr, ok := stmt.Cond.(*ast.BinaryExpr)
		if !ok {
			return
		}
		switch iterator := condExpr.Y.(type) {
		case *ast.BasicLit:
			for _, inner := range stmt.Body.List {
				tscAnalyser.Visit(inner, functionContext)
			}
		case *ast.Ident:
			if IsParam(iterator.Name, &functionContext.SymbolTable) {
				functionContext.CurrentDepth++
				for _, inner := range stmt.Body.List {
					tscAnalyser.Visit(inner, functionContext)
				}
				functionContext.CurrentDepth--
			}

		default:
			if ExprContainsParam(condExpr, &functionContext.SymbolTable) {
				functionContext.CurrentDepth++
				for _, inner := range stmt.Body.List {
					tscAnalyser.Visit(inner, functionContext)
				}
				functionContext.CurrentDepth--
			} else {
				for _, inner := range stmt.Body.List {
					tscAnalyser.Visit(inner, functionContext)
				}
			}
		}

	case *ast.IfStmt:
		for _, inner := range stmt.Body.List {
			tscAnalyser.Visit(inner, functionContext)
		}
		if stmt.Else != nil {
			tscAnalyser.Visit(stmt.Else, functionContext)
		}

	case *ast.LabeledStmt:
		tscAnalyser.Visit(stmt.Stmt, functionContext)

	case *ast.RangeStmt:
		functionContext.CurrentDepth++
		for _, inner := range stmt.Body.List {
			tscAnalyser.Visit(inner, functionContext)
		}
		functionContext.CurrentDepth--

	case *ast.ReturnStmt:
		for _, inner := range stmt.Results {
			tscAnalyser.Visit(inner, functionContext)
		}

	case *ast.SwitchStmt:
		for _, cases := range stmt.Body.List {
			tscAnalyser.Visit(cases, functionContext)
		}
	}

	functionContext.MaxDepth = float32(math.Max(float64(functionContext.CurrentDepth), float64(functionContext.MaxDepth)))
	functionContext.MaxMalloc = float32(math.Max(float64(functionContext.CurrentMalloc), float64(functionContext.MaxMalloc)))

}
