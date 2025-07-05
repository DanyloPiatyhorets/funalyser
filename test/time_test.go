package test

import (
	analyser "funalyser/analyser/go"
	"testing"
)

func TestTimeComplexity(t *testing.T) {
	file := "test_data/time_samples.go"
	funcs, err := analyser.Process(file)

	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]float32{
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

	for _, fn := range funcs {
		got := fn.Complexity.TimeIndex
		want, ok := expected[fn.Name]
		if !ok {
			t.Errorf("No expected result for %s", fn.Name)
		} else if got != want {
			t.Errorf("time for %s: expected %f, got %f", fn.Name, want, got)
		}
	}
}
