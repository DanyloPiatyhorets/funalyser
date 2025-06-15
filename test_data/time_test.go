package main

import "fmt"

var num int = 4
var text string = "hello"

func addNumbers(a, b int) int {
	x := 10
	var y string = "hello"
	fmt.Print(x)
	fmt.Print(y)
	return a + b
}


func countToTen(uselessParam int) {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}

func printItems(items []string) {
	for _, item := range items {
		fmt.Println(item)
	}
}

func nestedLoop(n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Println(i, j)
		}
	}
}

// TODO: add undefined time complexity
func loopForever() {
	for {
		fmt.Println("Running forever")
		break 
	}
}

func labeledBreak(param int) {
outer:
	for i := 0; i < param; i++ {
		for j := 0; j < param; j++ {
			if i == j {
				break outer
			}
			fmt.Println(i, j)
		}
	}
}

func conditionalLoop(doIt bool, n int) {
	if doIt {
		for i := 0; i < n; i++ {
			fmt.Println("DoIt:", i)
		}
	}
}

func loopInSwitch(x int) {
	switch x {
	case 1:
		for i := 0; i < 2; i++ {
			fmt.Println("Case 1:", i)
		}
	case 2:
		fmt.Println("No loop here")
	}
}