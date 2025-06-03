package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type Analyser interface {
    Visit(node ast.Node, functionContext *FunctionContext)
}

type TimeComplexityAnalyser struct {
    // MaxDepth int
}

type SpaceComplexityAnalyser struct {

}

type FunctionContext struct {
    Name string
    SymbolTable SymbolTable
    Depth int
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
        TimeComplexityIndex: functionContext.Depth,
        SpaceComplexityIndex: functionContext.Malloc,
        SymbolTable: functionContext.SymbolTable,
    }
}

type SymbolTable struct {
    Locals []string
    Params []string
    Globals []string
}

type FileContext struct {
    Globals []string
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
    // ast.Print(fset, file)

	var funcs []FunctionInfo
    var fileContext FileContext

    for _, declaration := range file.Decls {
        switch decl := declaration.(type){
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

            analysers := []Analyser{&TimeComplexityAnalyser{}, &SpaceComplexityAnalyser{},}

            WalkAST(decl.Body, analysers, functionContext)

            funcs = append(funcs, functionContext.GetFunctionInfo())
        }
    }

        return funcs, nil
    }

func WalkAST(block *ast.BlockStmt, analysers []Analyser, functionContext *FunctionContext){
    ast.Inspect(block, func(node ast.Node) bool {
        for _, analyser := range analysers {
            analyser.Visit(node, functionContext)
        }
        return true
    })
}

func (tcAnalyser *TimeComplexityAnalyser) Visit(node ast.Node, functionContext *FunctionContext){
    // tcAnalyser.MaxDepth = int(math.Max(float64(tcAnalyser.MaxDepth), float64(functionContext.Depth)))
    switch stmt := node.(type) {
    case *ast.ForStmt:
            condExpr, ok := stmt.Cond.(*ast.BinaryExpr)
            if !ok {
                return 
            }
            switch condExpr.Y.(type) {
            case *ast.BasicLit:
                // skip
            case *ast.Ident:
                functionContext.Depth++
            }

    case *ast.RangeStmt:
        functionContext.Depth++
    }   
}

func (scAnalyser *SpaceComplexityAnalyser) Visit(node ast.Node, functionContext *FunctionContext){
    switch node.(type){
    case *ast.AssignStmt:

    }
    // functionContext.Malloc++
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
