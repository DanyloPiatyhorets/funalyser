package test

import (
	analyser "funalyser/analyser/go"
	"testing"
)

func TestSpaceComplexity(t *testing.T) {
	file := "test_data/space_samples.go"
	funcs, err := analyser.Analyse(file, "")

	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]float32{
		"constantSpace":          0,
		"linearSpace":            1,
		"linearAppend":           1,
		"quadraticSpace":         2,
		"allocationPerIteration": 1,
		"recursiveStack":         1,
		"tailRecursive":          1,
		"fixedLoop":              0,
		"multiInputAllocation":   1,
		"mapSpace":               1,
		"reuseBuffer":            1,
		"conditionalAlloc":       1,
		"fixedAlloc":             0,
		"recurAlloc":             2,
	}

	for _, fn := range funcs {
		got := fn.Complexity.Space
		want, ok := expected[fn.Name]
		if !ok {
			t.Errorf("No expected result for %s", fn.Name)
		} else if got != want {
			t.Errorf("space for %s: expected %f, got %f", fn.Name, want, got)
		}
	}
}
