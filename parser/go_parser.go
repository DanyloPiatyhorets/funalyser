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
    SymbolTable         SymbolTable
	StartLine           int
	EndLine             int
}

type SymbolTable struct {
    Locals []string
    Params []string
    Globals []string
}

func Process(filePath string) ([]FunctionInfo, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
    // file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	node, err := parser.ParseFile(fset, filePath, src, parser.AllErrors)
	if err != nil {
		return nil, err
	}
    // ast.Print(fset, file)


	var funcs []FunctionInfo
    var funcSymbolTable SymbolTable

	ast.Inspect(node, func(node ast.Node) bool {
		currentFunction, isFunction := node.(*ast.FuncDecl) // this is a current node
		if isFunction {

            funcSymbolTable = getSymbolTableForFunction(currentFunction)
			pos := fset.Position(currentFunction.Pos())
			end := fset.Position(currentFunction.End())

			var funcStmtList []ast.Stmt = currentFunction.Body.List
			var timeComplexityIndex int = getMaxLoopDepth(funcStmtList, funcSymbolTable, 0, 0)

			funcs = append(funcs, FunctionInfo{
				Name:                currentFunction.Name.Name,
                SymbolTable: funcSymbolTable,
				TimeComplexityIndex: timeComplexityIndex,
				StartLine:           pos.Line,
				EndLine:             end.Line,
			})
		}
		return true
	})

	return funcs, nil
}

func getMaxLoopDepth(currentBodyList []ast.Stmt, funcSymbolTable SymbolTable, maxDepth int, currentDepth int) int {
	maxDepth = int(math.Max(float64(currentDepth), float64(maxDepth)))
	for _, stmt := range currentBodyList {
		switch stmtType := stmt.(type) {
		case *ast.ForStmt:
			maxDepth = getMaxLoopDepth(stmtType.Body.List, funcSymbolTable, maxDepth, currentDepth+1)
		case *ast.RangeStmt:
			maxDepth = getMaxLoopDepth(stmtType.Body.List, funcSymbolTable, maxDepth, currentDepth+1)
		case *ast.LabeledStmt:
			maxDepth = getMaxLoopDepth([]ast.Stmt{stmtType.Stmt}, funcSymbolTable, maxDepth, currentDepth)
		case *ast.IfStmt:
			maxDepth = getMaxLoopDepth(stmtType.Body.List, funcSymbolTable, maxDepth, currentDepth) // go deeper without incrementing
			if stmtType.Else != nil {
				if elseBlock, ok := stmtType.Else.(*ast.BlockStmt); ok {
					maxDepth = getMaxLoopDepth(elseBlock.List, funcSymbolTable, maxDepth, currentDepth)
				}
			}
		case *ast.SwitchStmt:
			for _, stmt := range stmtType.Body.List {
				if caseClause, isCaseClause := stmt.(*ast.CaseClause); isCaseClause {
					maxDepth = getMaxLoopDepth(caseClause.Body, funcSymbolTable, maxDepth, currentDepth)
				}
			}
		}
	}
	return maxDepth
}

func getSymbolTableForFunction(function *ast.FuncDecl) SymbolTable {
    var symbolTable SymbolTable

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
    // add short variable declarations (assignments) from the function AST
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