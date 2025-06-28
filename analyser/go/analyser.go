package analyser

import (
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

func Process(filePath string) ([]FunctionInfo, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	// ast.Print(fset, file)

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

            analyser := &TimeAndSpaceComplexityAnalyser{}

            for _, stmt := range decl.Body.List {
                analyser.Visit(stmt, functionContext)
	        }

			funcs = append(funcs, functionContext.GetFunctionInfo())
		}
	}

	return funcs, nil
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
		
	// BranchStmt

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
			time, space := GetRecursiveComplexity(stmt)
			functionContext.CurrentDepth += time
			functionContext.CurrentMalloc += space
			functionContext.RecursiveFanOut++
		}
	
	// CaseClause

	// CommClause

	// DeclStmt

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
				if functionContext.SymbolTable.IsParam(iterator.Name) {
					functionContext.CurrentDepth++
					for _, inner := range stmt.Body.List {
						tscAnalyser.Visit(inner, functionContext)
					}
					functionContext.CurrentDepth--
				}

			default:
				if ExpressionContainsParam(condExpr, &functionContext.SymbolTable) {
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

	// ParenExpr
	
	case *ast.RangeStmt:
		functionContext.CurrentDepth++
		for _, inner := range stmt.Body.List {
			tscAnalyser.Visit(inner, functionContext)
		}
		functionContext.CurrentDepth--

    case *ast.ReturnStmt:
        for _, inner := range stmt.Results{
            tscAnalyser.Visit(inner, functionContext)
        }
	}
	
	// SelectStmt

	// SwitchStmt

	// TypeSwitchStmt
    
	functionContext.MaxDepth = float32(math.Max(float64(functionContext.CurrentDepth), float64(functionContext.MaxDepth)))
	functionContext.MaxMalloc = float32(math.Max(float64(functionContext.CurrentMalloc), float64(functionContext.MaxMalloc)))

}