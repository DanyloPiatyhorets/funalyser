package analyser

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
    switch stmt := node.(type) {
    case *ast.ForStmt:
            condExpr, ok := stmt.Cond.(*ast.BinaryExpr)
            if !ok {
                return 
            }
            switch iterator := condExpr.Y.(type) {
            case *ast.BasicLit:
                // skip
            case *ast.Ident:
                if functionContext.SymbolTable.IsParam(iterator.Name){
                    functionContext.MaxDepth++
                }
            }

    case *ast.RangeStmt:
        functionContext.MaxDepth++
    case *ast.CallExpr:
        funIdent, ok := stmt.Fun.(*ast.Ident)
        if !ok {
            return
        }
        switch funIdent.Name {
            case functionContext.Name:
            functionContext.MaxDepth += 1
        }
        }
    }   

func (scAnalyser *SpaceComplexityAnalyser) Visit(node ast.Node, functionContext *FunctionContext){
    switch stmt:= node.(type){
    case *ast.ForStmt:
            condExpr, ok := stmt.Cond.(*ast.BinaryExpr)
            if !ok {
                return 
            }
            switch iterator := condExpr.Y.(type) {
            case *ast.BasicLit:
                // skip
            case *ast.Ident:
                if functionContext.SymbolTable.IsParam(iterator.Name){
                    functionContext.CurrentDepth++
                }
            }

    case *ast.RangeStmt:
        functionContext.CurrentDepth++
        
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
                        functionContext.Malloc = 1 + functionContext.CurrentDepth
                    }
            }
            case *ast.MapType:
                functionContext.Malloc = 1 + functionContext.CurrentDepth
            }
    
        case "append":
            functionContext.Malloc = functionContext.CurrentDepth
        case functionContext.Name:
            functionContext.Malloc += 1
        }
    }
}