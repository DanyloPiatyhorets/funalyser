package cmd

import (
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "funalyser",
    Short: "A CLI tool to analyze functions in Go source code",
    Long:  `funalyser analyzes Go files to extract function-level insights.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
