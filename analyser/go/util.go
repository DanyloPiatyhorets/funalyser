package analyser

import (
	"go/ast"
	"go/token"
)

type FileContext struct {
	Globals []string
}

type FunctionContext struct {
	Name          	string
	SymbolTable   	SymbolTable
	CurrentDepth  	float32
	MaxDepth      	float32
	CurrentMalloc 	float32
    MaxMalloc     	float32
	RecursiveFanOut int
}

func NewFunctionContext() *FunctionContext {
	return &FunctionContext{}
}

type FunctionInfo struct {
	Name                 string
	TimeComplexityIndex  float32
	SpaceComplexityIndex float32
	SymbolTable          SymbolTable
	FanOut				 int
}

func (functionContext *FunctionContext) GetFunctionInfo() FunctionInfo {
	return FunctionInfo{
		Name:                 functionContext.Name,
		TimeComplexityIndex:  functionContext.MaxDepth,
		SpaceComplexityIndex: functionContext.MaxMalloc,
		SymbolTable:          functionContext.SymbolTable,
		FanOut:				  functionContext.RecursiveFanOut,
	}
}

type SymbolTable struct {
	Locals  []string
	Params  []string
	Globals []string
}

func (st *SymbolTable) IsParam(name string) bool {
	for _, param := range st.Params {
		if param == name {
			return true
		}
	}
	return false
}

func (analyser *FileContext) getSymbolTableForFunction(function *ast.FuncDecl) SymbolTable {
	// var symbolTable SymbolTable
	symbolTable := SymbolTable{
		Globals: analyser.Globals,
	}

	// add parameters from the function AST
	for _, params := range function.Type.Params.List {
		for _, param := range params.Names {
			symbolTable.Params = append(symbolTable.Params, param.Name)
		}
	}
	// add short variable declarations (assignments) from the function AST
	for _, stmt := range function.Body.List {
		if assingStmt, ok := stmt.(*ast.AssignStmt); ok && assingStmt.Tok == token.DEFINE {
			for _, exp := range assingStmt.Lhs {
				if identifier, ok := exp.(*ast.Ident); ok {
					symbolTable.Locals = append(symbolTable.Locals, identifier.Name)
				}
			}
		}
	}
	// add variable declarations from the function AST
	for _, stmt := range function.Body.List {
		if declStmt, ok := stmt.(*ast.DeclStmt); ok {
			if genDecl, ok := declStmt.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, identifier := range valueSpec.Names {
							symbolTable.Locals = append(symbolTable.Locals, identifier.Name)
						}
					}
				}
			}
		}
	}

	return symbolTable
}

func ExpressionContainsParam(expr ast.Expr, symbolTable *SymbolTable) bool {
    switch exp := expr.(type) {
	case *ast.Ident:
		return symbolTable.IsParam(exp.Name)
    case *ast.CallExpr:
        if funIdent, ok := exp.Fun.(*ast.Ident); ok && funIdent.Name == "len" {
            if len(exp.Args) == 1 {
                if argIdent, ok := exp.Args[0].(*ast.Ident); ok {
                    return symbolTable.IsParam(argIdent.Name)
                }
            }
        }
    case *ast.BinaryExpr:
        return ExpressionContainsParam(exp.X, symbolTable) || ExpressionContainsParam(exp.Y, symbolTable)
	case *ast.IndexExpr:
        return ExpressionContainsParam(exp.X, symbolTable) || ExpressionContainsParam(exp.Index, symbolTable)
    case *ast.ParenExpr:
        return ExpressionContainsParam(exp.X, symbolTable)
    }
    return false
}

func GetRecursiveComplexity(expr *ast.CallExpr) (float32, float32) {
	var time float32 = 0
	var space float32 = 0
	if exp, ok := expr.Args[0].(*ast.BinaryExpr); ok  {
		switch exp.Op.String() {
		case "+", "-":
			time = 1
			space = 1
		case "*", "/":
			time = 0.5
			space = 0.5
		}
	}
	return time, space
}