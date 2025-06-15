package analyser

import (
	"go/ast"
	"go/token"	
)

type FileContext struct {
    Globals []string
}

type FunctionContext struct {
    Name string
    SymbolTable SymbolTable
    CurrentDepth int
    MaxDepth int
    Malloc int
}

func NewFunctionContext() *FunctionContext {
    return &FunctionContext{}
}

type FunctionInfo struct {
	Name                string
	TimeComplexityIndex int
    SpaceComplexityIndex int
    SymbolTable         SymbolTable
}

func (functionContext *FunctionContext) GetFunctionInfo() FunctionInfo {
    return FunctionInfo {
        Name: functionContext.Name,
        TimeComplexityIndex: functionContext.MaxDepth,
        SpaceComplexityIndex: functionContext.Malloc,
        SymbolTable: functionContext.SymbolTable,
    }
}

type SymbolTable struct {
    Locals []string
    Params []string
    Globals []string
}

func (st * SymbolTable) IsParam(name string) bool {
    for _, param := range st.Params {
        if (param == name){
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
            for _, exp := range assingStmt.Lhs{
                if identifier, ok := exp.(*ast.Ident); ok{
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