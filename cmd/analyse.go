package cmd

import (
	"fmt"
	"funalyser/analyser/go"
	"strconv"
	"github.com/spf13/cobra"
)

var fileAnalysis = &cobra.Command{
	Use:   "analyse [file.go]",
	Short: "Analyse functions in a Go source file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string){
		functionName, _ := cmd.Flags().GetString("func") 
		if functionName != "" {
			analyseFunction(args, functionName)
		} else{
			analyseFile(args)

		}
	},
}

func init() {
	rootCmd.AddCommand(fileAnalysis)
	rootCmd.PersistentFlags().String("func", "", "Name of the function to analyse")
}

func analyseFile(args []string) {
	functionInfos, err := analyser.AnalyseFile(args[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, fn := range functionInfos {
		printFunctionReport(fn)
	}
}

func analyseFunction(args []string, functionName string) {
	functionInfo, err := analyser.AnalyseFunction(args[0], functionName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	printFunctionReport(functionInfo)
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
	fmt.Printf("  • Time Complexity:   %s\n", parseComplexityIndexToString(fn.Complexity.TimeIndex))
	fmt.Printf("  • Space Complexity:  %s\n", parseComplexityIndexToString(fn.Complexity.SpaceIndex))

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
