package parser

import (
    "go/ast"
    "go/parser"
    "go/token"
    "os"
)

type FunctionInfo struct {
    Name      string
    StartLine int
    EndLine   int
}

func ExtractFunctions(filePath string) ([]FunctionInfo, error) {
    src, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, filePath, src, parser.AllErrors)
    if err != nil {
        return nil, err
    }

    var funcs []FunctionInfo
    ast.Inspect(node, func(n ast.Node) bool {
        fn, ok := n.(*ast.FuncDecl)
        if ok {
            pos := fset.Position(fn.Pos())
            end := fset.Position(fn.End())
            funcs = append(funcs, FunctionInfo{
                Name:      fn.Name.Name,
                StartLine: pos.Line,
                EndLine:   end.Line,
            })
        }
        return true
    })

    return funcs, nil
}
