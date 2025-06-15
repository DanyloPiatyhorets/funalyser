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

	spaceExpected := map[string]int{
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
	// timeExpected := map[string]int{
	// 	"addNumbers":      0,
	// 	"countToTen":      0,
	// 	"printItems":      1,
	// 	"nestedLoop":      2,
	// 	"loopForever":     1,
	// 	"labeledBreak":    2,
	// 	"conditionalLoop": 1,
	// 	"loopInSwitch":    0,
	// }
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


	// expectedTC := timeExpected[fn.Name]
	expectedSC := spaceExpected[fn.Name]

	// tcCorrect := fn.TimeComplexityIndex == expectedTC
	scCorrect := fn.SpaceComplexityIndex == expectedSC

	// tcSymbol := "❌"
	// if tcCorrect {
	// 	tcSymbol = "✅"
	// }
	scSymbol := "❌"
	if scCorrect {
		scSymbol = "✅"
	}
	

    fmt.Printf(
		// | Time: %-7s %s
    "Func: %-15s | Space: %-2d %s | Params: %v | Locals: %v | Globals: %v\n",
    fn.Name,
    // parseIndexToTimeComplexity(fn.TimeComplexityIndex),
    // tcSymbol,
    fn.SpaceComplexityIndex,
    scSymbol,
    fn.SymbolTable.Params,
    fn.SymbolTable.Locals,
    fn.SymbolTable.Globals,
)   
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
