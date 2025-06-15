package test

import (
	"testing"
    "funalyser/analyser/go"
)

func TestSpaceComplexity(t *testing.T) {
    file := "../test_data/space_test.go"
	funcs, err := analyser.Process(file)

    if err != nil {
        t.Fatal(err)
    }

    expected := map[string]int{
        "constantSpace":     0,
        "linearSpace":       1,
        "linearAppend": 1,
        "quadraticSpace":    2,
        "allocationPerIteration":       1,
        "recursiveStack":  1,
        "tailRecursive": 1,
        "fixedLoop": 0,
        "multiInputAllocation": 1,
        "mapSpace": 1,
        "reuseBuffer": 1,
        "conditionalAlloc": 1,
        "fixedAlloc": 0,
        "recurAlloc": 2,
    }

    for _, fn := range funcs {
        got := fn.SpaceComplexityIndex
        want, ok := expected[fn.Name]
        if !ok {
            t.Errorf("No expected result for %s", fn.Name)
        } else if got != want {
            t.Errorf("space for %s: expected %d, got %d", fn.Name, want, got)
        }
    }
}
