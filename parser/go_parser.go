package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
)

type FunctionInfo struct {
	Name                string
	TimeComplexityIndex int
	StartLine           int
	EndLine             int
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
				Name:                currentFunction.Name.Name,
				TimeComplexityIndex: timeComplexityIndex,
				StartLine:           pos.Line,
				EndLine:             end.Line,
			})
		}
		return true
	})

	return funcs, nil
}

func getMaxLoopDepth(currentBodyList []ast.Stmt, maxDepth int, currentDepth int) int {
	maxDepth = int(math.Max(float64(currentDepth), float64(maxDepth)))
	for _, stmt := range currentBodyList {
		switch stmtType := stmt.(type) {
		case *ast.ForStmt:
			maxDepth = getMaxLoopDepth(stmtType.Body.List, maxDepth, currentDepth+1)
		case *ast.RangeStmt:
			maxDepth = getMaxLoopDepth(stmtType.Body.List, maxDepth, currentDepth+1)
		case *ast.LabeledStmt:
			maxDepth = getMaxLoopDepth([]ast.Stmt{stmtType.Stmt}, maxDepth, currentDepth)
		case *ast.IfStmt:
			maxDepth = getMaxLoopDepth(stmtType.Body.List, maxDepth, currentDepth) // go deeper without incrementing
			if stmtType.Else != nil {
				if elseBlock, ok := stmtType.Else.(*ast.BlockStmt); ok {
					maxDepth = getMaxLoopDepth(elseBlock.List, maxDepth, currentDepth)
				}
			}
		case *ast.SwitchStmt:
			for _, stmt := range stmtType.Body.List {
				if caseClause, isCaseClause := stmt.(*ast.CaseClause); isCaseClause {
					maxDepth = getMaxLoopDepth(caseClause.Body, maxDepth, currentDepth)
				}
			}
		}
	}
	return maxDepth
}
