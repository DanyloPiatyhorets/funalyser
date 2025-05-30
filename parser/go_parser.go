package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
)

type FunctionInfo struct {
	Name           string
	TimeComplexity int
	StartLine      int
	EndLine        int
}

func Process(filePath string) ([]FunctionInfo, error) {
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


	ast.Inspect(node, func(node ast.Node) bool {
		currentFunction, isFunction := node.(*ast.FuncDecl) // this is a current node
		if isFunction {
            
			pos := fset.Position(currentFunction.Pos())
			end := fset.Position(currentFunction.End())
			
            var funcStmtList []ast.Stmt = currentFunction.Body.List
            var timeComplexityIndex int = getMaxLoopDepth(funcStmtList, 0, 0)

            funcs = append(funcs, FunctionInfo{
				Name:      currentFunction.Name.Name,
                TimeComplexity: timeComplexityIndex,
				StartLine: pos.Line,
				EndLine:   end.Line,
			})
		}
		return true
	})

	return funcs, nil
}

func getMaxLoopDepth(currentBodyList []ast.Stmt, maxDepth int , currentDepth int) int {
    maxDepth = int(math.Max(float64(currentDepth), float64(maxDepth)))
	fset := token.NewFileSet()
    fmt.Print("For current List currentDepth: ", currentDepth, " maxDepth: ", maxDepth)
    for _, stmt := range currentBodyList {
	    ast.Print(fset, stmt)
    }
    for _, forStmt := range currentBodyList{
        currentLoop, isLoop := forStmt.(*ast.ForStmt)
        if isLoop {
            return getMaxLoopDepth(currentLoop.Body.List, maxDepth, currentDepth+1)
        }
    }
    return maxDepth
}
