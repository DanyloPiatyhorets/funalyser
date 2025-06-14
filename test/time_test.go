package test

import (
	"testing"
    "funalyser/analyser"
)

func TestTimeComplexity(t *testing.T) {
    file := "../test_data/time_samples.go"
	funcs, err := analyser.Process(file)
	
    if err != nil {
        t.Fatal(err)
    }

    expected := map[string]int{
       "addNumbers":      0,
		"countToTen":      0,
		"printItems":      1,
		"nestedLoop":      2,
		"loopForever":     0,
		"labeledBreak":    2,
		"conditionalLoop": 1,
		"loopInSwitch":    0,
    }

    for _, fn := range funcs {
        got := fn.TimeComplexityIndex
        want, ok := expected[fn.Name]
        if !ok {
            t.Errorf("No expected result for %s", fn.Name)
        } else if got != want {
            t.Errorf("time for %s: expected %d, got %d", fn.Name, want, got)
        }
    }
}
