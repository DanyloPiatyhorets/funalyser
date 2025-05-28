package main

import (
    "fmt"
    "go/parser"
    "go/token"
    "go/ast"
)

func main() {
    src := `
    package main
    import "fmt"

    func greet(name string) string {
        return "Hello, " + name
    }`

    fset := token.NewFileSet()
    file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Print the entire AST
    ast.Print(fset, file)
}
