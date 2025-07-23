package cmd

import (
	"fmt"
	analyser "funalyser/analyser/go"
	"strconv"
	"github.com/spf13/cobra"
)

var fileAnalysis = &cobra.Command{
	Use:   "analyse [file.go]",
	Short: "Analyse functions in a source file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		functionName, _ := cmd.Flags().GetString("func")
		funcsInfo, err := analyser.Analyse(args[0], functionName)
		if err != nil {
			fmt.Println("❌", err)
			return
		}
		for _, fn := range funcsInfo {
			printFunctionReport(fn)
		}
	},
}

var info = &cobra.Command{
	Use:   "info",
	Short: "Information about funalyser functionality and use cases",
	Run: func(cmd *cobra.Command, args []string) {
		printInfo()
	},
}

// TODO: think of a set of commands and flags for the extended first verion functionality
// TODO: think of making a parser speak via json to enable java parsing
// 		- think if I need to do it now or in the future
// TODO: polish everything
// TODO: ask chatgpt what else I could also do

func init() {
	rootCmd.AddCommand(fileAnalysis)
	rootCmd.AddCommand(info)
	rootCmd.PersistentFlags().String("func", "", "Name of the function to analyse")
}

func printFunctionReport(fn analyser.FunctionInfo) {
	fmt.Println()
	fmt.Println("───────────────────────────────────────────")
	fmt.Printf("🔍 Function: %s\n", fn.Name)
	fmt.Println("─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─")

	fmt.Println("📊 Analysis Summary:")
	fmt.Printf("  • Recursive:         %s\n", checkmark(fn.FanOut != 0))
	if fn.FanOut > 0 {
		fmt.Printf("  • Fan-out Factor:    %d %s\n", fn.FanOut, fanOutHint(fn.FanOut))
	}
	fmt.Printf("  • Time Complexity:   %s\n", parseComplexityIndexToString(fn.Complexity.Time))
	fmt.Printf("  • Space Complexity:  %s\n", parseComplexityIndexToString(fn.Complexity.Space))

	if fn.FanOut > 1 {
		fmt.Println("📌 Notes:")
		fmt.Println("  • Multiple recursive calls detected (fan-out > 1).")
		fmt.Println("    ➤ Consider checking if this leads to exponential growth.")
	}

	fmt.Println("───────────────────────────────────────────")
	fmt.Println(" ")
}

func parseComplexityIndexToString(maxLoopDepth float32) string {
	switch maxLoopDepth {
	case 0:
		return "O(1)"
	case 0.5:
		return "O(log n)"
	case 1:
		return "O(n)"
	case 1.5:
		return "O(n*log n)"
	}
	return "O(n^" + strconv.Itoa(int(maxLoopDepth)) + ")"
}

func checkmark(ok bool) string {
	if ok {
		return "Yes"
	}
	return "No"
}

func fanOutHint(fanOut int) string {
	if fanOut > 1 {
		return "⚠️  Potentially exponential"
	}
	return ""
}

func printInfo() {
	fmt.Println(`
Funalyser – Code Complexity Analyser 

This tool is designed with developers in mind to help them have a quick analysis of methods in a specific file

🧠 Features:
• Analyses time and space complexity of functions
• Detects recursive patterns and fan-out factors
• Tracks memory allocation (for example make and append)

🔁 Recognises:
• Linear/logarithmic recursion
• Multiple recursive calls
• Memory-intensive constructs

✅ Currently supported languges:
• Golang 
• poteantial support of Java and other languages in the future

🚀 Usage Examples:
• funalyser analyse ./main.go
	- gives an analysis for each function in the specified file 
• funalyser analyse ./main.go --func MergeSort
	- gives an analysis for a specific function in the file

👤 Author: Danylo Piatyhorets
📚 GitHub: https://github.com/DanyloPiatyhorets/funalyser
	`)
}
