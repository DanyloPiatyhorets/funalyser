package analyser

import (
	"go/ast"
	"go/token"
	"strings"
)

type FileContext struct {
	Globals []string
}

type FunctionContext struct {
	Name            string
	SymbolTable     SymbolTable
	CurrentDepth    float32
	MaxDepth        float32
	CurrentMalloc   float32
	MaxMalloc       float32
	RecursiveFanOut int
}

type FunctionInfo struct {
	Name        string
	Complexity  Complexity
	SymbolTable SymbolTable
	FanOut      int
}

type SymbolTable struct {
	Locals  []string
	Params  []string
	Globals []string
}

type Complexity struct {
	Time  float32
	Space float32
}


func ParseContextToInfo(functionContext *FunctionContext) FunctionInfo {
	return FunctionInfo{
		Name: functionContext.Name,
		Complexity: Complexity {
			Time:  functionContext.MaxDepth,
			Space: functionContext.MaxMalloc,
		},
		SymbolTable: functionContext.SymbolTable,
		FanOut:      functionContext.RecursiveFanOut,
	}
}

func GetFileContext(file *ast.File) FileContext {
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
		}
	}
	return fileContext

}

func GetFunctionContext(decl *ast.FuncDecl, fileContext *FileContext) *FunctionContext {
	functionContext := &FunctionContext{}
	functionContext.Name = decl.Name.Name

	functionContext.SymbolTable.Globals = fileContext.Globals

	// add parameters
	for _, params := range decl.Type.Params.List {
		for _, param := range params.Names {
			functionContext.SymbolTable.Params = append(functionContext.SymbolTable.Params, param.Name)
		}
	}
	// add short variable declarations (assignments)
	for _, stmt := range decl.Body.List {
		if assingStmt, ok := stmt.(*ast.AssignStmt); ok && assingStmt.Tok == token.DEFINE {
			for _, exp := range assingStmt.Lhs {
				if identifier, ok := exp.(*ast.Ident); ok {
					functionContext.SymbolTable.Locals = append(functionContext.SymbolTable.Locals, identifier.Name)
				}
			}
		}
	}
	// add variable declarations
	for _, stmt := range decl.Body.List {
		if declStmt, ok := stmt.(*ast.DeclStmt); ok {
			if genDecl, ok := declStmt.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, identifier := range valueSpec.Names {
							functionContext.SymbolTable.Locals = append(functionContext.SymbolTable.Locals, identifier.Name)
						}
					}
				}
			}
		}
	}

	return functionContext
}

func GetRecursiveComplexity(expr *ast.CallExpr) (float32, float32) {
	var time float32 = 0
	var space float32 = 0
	if exp, ok := expr.Args[0].(*ast.BinaryExpr); ok {
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

func isFunctionName(funcDecl *ast.FuncDecl, funcName string) bool {
	return strings.ToLower(funcName) == strings.ToLower(funcDecl.Name.Name)
}

func IsParam(name string, symbolTable *SymbolTable) bool {
	for _, param := range symbolTable.Params {
		if param == name {
			return true
		}
	}
	return false
}

func ExprContainsParam(expr ast.Expr, symbolTable *SymbolTable) bool {
	switch exp := expr.(type) {
	case *ast.Ident:
		return IsParam(exp.Name, symbolTable)
	case *ast.CallExpr:
		if funIdent, ok := exp.Fun.(*ast.Ident); ok && funIdent.Name == "len" {
			if len(exp.Args) == 1 {
				if argIdent, ok := exp.Args[0].(*ast.Ident); ok {
					return IsParam(argIdent.Name, symbolTable)
				}
			}
		}
	case *ast.BinaryExpr:
		return ExprContainsParam(exp.X, symbolTable) || ExprContainsParam(exp.Y, symbolTable)
	case *ast.IndexExpr:
		return ExprContainsParam(exp.X, symbolTable) || ExprContainsParam(exp.Index, symbolTable)
	case *ast.ParenExpr:
		return ExprContainsParam(exp.X, symbolTable)
	}
	return false
}
