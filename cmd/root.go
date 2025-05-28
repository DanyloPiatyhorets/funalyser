package cmd

import (
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "funanalyser",
    Short: "A CLI tool to analyze functions in Go source code",
    Long:  `funanalyser analyzes Go files to extract function-level insights.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
