package cmd

import (
	"fmt"
	"funalyser/parser"
	"strconv"

	"github.com/spf13/cobra"
)

var analyseCmd = &cobra.Command{
	Use:   "analyse [file.go]",
	Short: "Analyse functions in a Go source file",
	Args:  cobra.ExactArgs(1),
	Run:   analyse,
}

func init() {
	rootCmd.AddCommand(analyseCmd)
}

func analyse(cmd *cobra.Command, args []string) {
	functions, err := parser.Process(args[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, fn := range functions {
		printFunctionReport(fn)
	}
}

func printFunctionReport(fn parser.FunctionInfo) {
	// Hardcoded expected values
	expected := map[string]int{
		"addNumbers":      0,
		"countToTen":      1,
		"printItems":      1,
		"nestedLoop":      2,
		"loopForever":     1,
		"labeledBreak":    2,
		"conditionalLoop": 1,
		"loopInSwitch":    1,
	}

	expectedTC, ok := expected[fn.Name]
	correct := ok && expectedTC == fn.TimeComplexityIndex

	symbol := "❌"
	if correct {
		symbol = "✅"
	}

	fmt.Printf("Function: %-18s | Time Complexity: %-s %s\n",
		fn.Name, parseIndexToTimeComplexity(fn.TimeComplexityIndex), symbol)
}

func parseIndexToTimeComplexity(maxLoopDepth int) string {
	switch maxLoopDepth {
	case 0:
		return "O(1)"
	case 1:
		return "O(n)"
	}
	return "O(n^" + strconv.Itoa(maxLoopDepth) + ")"
}
