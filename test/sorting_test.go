package test

import (
	"funalyser/analyser/go"
	"testing"
)

func TestSorting(t *testing.T) {
    file := "test_data/sorting_samples.go"
	funcs, err := analyser.Process(file)

    if err != nil {
        t.Fatal(err)
    }

	expected := map[string]analyser.Complexity{
		"BubbleSort": 	 {TimeIndex: 2, SpaceIndex: 0},
		"InsertionSort": {TimeIndex: 2, SpaceIndex: 0}, 
    	"SelectionSort": {TimeIndex: 2, SpaceIndex: 0},
    	// "QuickSort":     {TimeIndex: 1.5, SpaceIndex: 0},
    	// "MergeSort":     {TimeIndex: 1.5, SpaceIndex: 1},

	}
    for _, fn := range funcs {
        var got analyser.Complexity = fn.Complexity
        want, ok := expected[fn.Name]
        if !ok {
            t.Errorf("No expected result for %s", fn.Name)
        } else if got != want {
            t.Errorf("space for %s: expected %f, got %f", fn.Name, want, got)
        }
    }
}
