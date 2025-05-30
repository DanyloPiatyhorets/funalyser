package main

import (
    "fmt"
    "go/parser"
    "go/token"
    "go/ast"
)

func main() {
    // src := `
    // package main
    // import "fmt"

    // func greet(name string) string {
    //     return "Hello, " + name
    // }

    // func haveLoop() {
    //     for i := 0; i < 10; i++ {
    //         fmt.Println(i)
    //         for j := 0; j < 10; j++ {
    //             fmt.Println(j)
    //         }
    //     }
    // }    
    // `
    src := `
    package main

import "fmt"

    func greet(name string) string {
        return "Hello, " + name
    }

    func haveLoop() {
        for i := 0; i < 10; i++ {
            fmt.Println(i)
            for j := 0; j < 10; j++ {
                fmt.Println(j)
            }
        }
    }    
    `
    

    fset := token.NewFileSet()
    file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Print the entire AST
    ast.Print(fset, file)
}
