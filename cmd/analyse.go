package cmd

import (
    "fmt"
    "funalyser/parser"
    "github.com/spf13/cobra"
)

var analyseCmd = &cobra.Command{
    Use:   "analyse [file.go]",
    Short: "Analyse functions in a Go source file",
    Args:  cobra.ExactArgs(1),
    Run: analyse,
}

func init() {
    rootCmd.AddCommand(analyseCmd)
}

func analyse(cmd *cobra.Command, args []string){
	functions, err := parser.ExtractFunctions(args[0])
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        for _, fn := range functions {
            fmt.Printf("Function: %s (Lines %d–%d)\n", fn.Name, fn.StartLine, fn.EndLine)
        }
}

