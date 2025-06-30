package cmd

import (
	"fmt"
	"funalyser/analyser/go"
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
	functions, err := analyser.Process(args[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, fn := range functions {
		printFunctionReport(fn)
	}
}

func printFunctionReport(fn analyser.FunctionInfo) {

	spaceExpected := map[string]float32{
    "constantSpace":         0,
    "linearSpace":           1,
    "linearAppend":          1,
    "quadraticSpace":        2,
    "allocationPerIteration": 1,
    "recursiveStack":        1,
    "tailRecursive":         1,
    "fixedLoop":             0,
    "multiInputAllocation":  1,
    "mapSpace":              1,
    "reuseBuffer":           1,
    "conditionalAlloc":      1,
    "fixedAlloc": 0,
	"recurAlloc": 2,

	}	

	// Hardcoded timeExpected values
	timeExpected := map[string]float32{
		"addNumbers":      0,
		"countToTen":      0,
		"printItems":      1,
		"nestedLoop":      2,
		"loopForever":     0,
		"labeledBreak":    2,
		"conditionalLoop": 1,
		"loopInSwitch":    1,
		"recursion":       1,
	}
	// spaceExpected := map[string]int{
	// 	"addNumbers":       0,
	// 	"countToTen":       0,
	// 	"printItems":       1,
	// 	"nestedLoop":       2,
	// 	"loopForever":      0,
	// 	"labeledBreak":     1,
	// 	"conditionalLoop":  0,
	// 	"loopInSwitch":     1,
	// 	"hybridAlloc":      2,
	// 	"appendDynamic":    2,
	// }


	expectedTC := timeExpected[fn.Name]
	expectedSC := spaceExpected[fn.Name]

	tcCorrect := fn.TimeComplexityIndex == expectedTC
	scCorrect := fn.SpaceComplexityIndex == expectedSC

	tcSymbol := "❌"
	if tcCorrect {
		tcSymbol = "✅"
	}
	scSymbol := "❌"
	if scCorrect {
		scSymbol = "✅"
	}
	

    fmt.Printf(
    "Func: %-15s | Time: %-7s %s | Space: %-7s %s | FanOut=%d | Params: %v | Locals: %v | Globals: %v\n",
    fn.Name,
    parseIndexToTimeComplexity(fn.TimeComplexityIndex),
    tcSymbol,
    parseIndexToTimeComplexity(fn.SpaceComplexityIndex),
    scSymbol,
	fn.FanOut,
    fn.SymbolTable.Params,
    fn.SymbolTable.Locals,
    fn.SymbolTable.Globals,
)   
}

func parseIndexToTimeComplexity(maxLoopDepth float32) string {
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
