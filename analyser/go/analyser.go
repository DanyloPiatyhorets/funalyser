package analyser

import (
	// "fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
)

type Analyser interface {
	Visit(node ast.Node, functionContext *FunctionContext)
}

type TimeComplexityAnalyser struct {
}

type SpaceComplexityAnalyser struct {
}

func Process(filePath string) ([]FunctionInfo, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	// node, err := parser.ParseFile(fset, filePath, src, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	ast.Print(fset, file)

	var funcs []FunctionInfo
	var fileContext FileContext

	for _, declaration := range file.Decls {
		switch decl := declaration.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				if valSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, identifier := range valSpec.Names {
						fileContext.Globals = append(fileContext.Globals, identifier.Name)
					}
				}
			}
		case *ast.FuncDecl:

			functionContext := NewFunctionContext()
			functionContext.SymbolTable = fileContext.getSymbolTableForFunction(decl)
			functionContext.Name = decl.Name.Name

			analysers := []Analyser{&TimeComplexityAnalyser{}, &SpaceComplexityAnalyser{}}

            WalkAST(decl.Body, analysers, functionContext)

			funcs = append(funcs, functionContext.GetFunctionInfo())
		}
	}

	return funcs, nil
}

func WalkAST(block *ast.BlockStmt, analysers []Analyser, functionContext *FunctionContext) {
	for _, stmt := range block.List {
		for _, analyser := range analysers {
			analyser.Visit(stmt, functionContext)
		}
	}
}

func (tcAnalyser *TimeComplexityAnalyser) Visit(node ast.Node, functionContext *FunctionContext) {
	functionContext.MaxDepth = int(math.Max(float64(functionContext.CurrentDepth), float64(functionContext.MaxDepth)))

	switch stmt := node.(type) {
	case *ast.BlockStmt:
		for _, inner := range stmt.List {
			tcAnalyser.Visit(inner, functionContext)
		}
	case *ast.LabeledStmt:
		tcAnalyser.Visit(stmt.Stmt, functionContext)

	case *ast.ForStmt:
		condExpr, ok := stmt.Cond.(*ast.BinaryExpr)
		if !ok {
			return
		}
		switch iterator := condExpr.Y.(type) {
		case *ast.BasicLit:
			for _, inner := range stmt.Body.List {
				tcAnalyser.Visit(inner, functionContext)
			}
		case *ast.Ident:
			if functionContext.SymbolTable.IsParam(iterator.Name) {
				functionContext.CurrentDepth++
				for _, inner := range stmt.Body.List {
					tcAnalyser.Visit(inner, functionContext)
				}
				functionContext.CurrentDepth--
			}
		}

    case *ast.ReturnStmt:
        for _, inner := range stmt.Results{
            tcAnalyser.Visit(inner, functionContext)
        }
	case *ast.RangeStmt:
		functionContext.CurrentDepth++
		for _, inner := range stmt.Body.List {
			tcAnalyser.Visit(inner, functionContext)
		}
		functionContext.CurrentDepth--

	case *ast.IfStmt:
		for _, inner := range stmt.Body.List {
			tcAnalyser.Visit(inner, functionContext)
		}
		if stmt.Else != nil {
			tcAnalyser.Visit(stmt.Else, functionContext)
		}

		// TODO: add recursion and logarithmic complexity here. Check the arguments passed during
		// the recursion to identify what is the complexity
	}
}

func (scAnalyser *SpaceComplexityAnalyser) Visit(node ast.Node, functionContext *FunctionContext) {

	functionContext.MaxDepth = int(math.Max(float64(functionContext.CurrentDepth), float64(functionContext.MaxDepth)))
	functionContext.MaxMalloc = int(math.Max(float64(functionContext.CurrentMalloc), float64(functionContext.MaxMalloc)))

	switch stmt := node.(type) {
	case *ast.BlockStmt:
		for _, inner := range stmt.List {
			scAnalyser.Visit(inner, functionContext)
		}
	case *ast.LabeledStmt:
		scAnalyser.Visit(stmt.Stmt, functionContext)

	case *ast.ForStmt:
		condExpr, ok := stmt.Cond.(*ast.BinaryExpr)
		if !ok {
			return
		}
		switch iterator := condExpr.Y.(type) {
		case *ast.BasicLit:
			for _, inner := range stmt.Body.List {
				scAnalyser.Visit(inner, functionContext)
			}
		case *ast.Ident:
			if functionContext.SymbolTable.IsParam(iterator.Name) {
				functionContext.CurrentDepth++
				for _, inner := range stmt.Body.List {
					scAnalyser.Visit(inner, functionContext)
				}
				functionContext.CurrentDepth--
			}
		}

	case *ast.RangeStmt:
		functionContext.CurrentDepth++
		for _, inner := range stmt.Body.List {
			scAnalyser.Visit(inner, functionContext)
		}
		functionContext.CurrentDepth--

	case *ast.IfStmt:
		for _, inner := range stmt.Body.List {
			scAnalyser.Visit(inner, functionContext)
		}
		if stmt.Else != nil {
			scAnalyser.Visit(stmt.Else, functionContext)
		}
    case *ast.AssignStmt:
        scAnalyser.Visit(stmt.Rhs[0], functionContext)
    case *ast.ReturnStmt:
        for _, inner := range stmt.Results{
            scAnalyser.Visit(inner, functionContext)
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
                    if functionContext.SymbolTable.IsParam(size.Name) {
                        functionContext.CurrentMalloc = 1 + functionContext.CurrentDepth
                    }
            }
            case *ast.MapType:
                functionContext.CurrentMalloc = 1 + functionContext.CurrentDepth
            }
    
        case "append":
            functionContext.CurrentMalloc = functionContext.CurrentDepth
        case functionContext.Name:
            functionContext.CurrentMalloc += 1
        }


		// TODO: add recursion and logarithmic complexity here. Check the arguments passed during
		// the recursion to identify what is the complexity
	}
}