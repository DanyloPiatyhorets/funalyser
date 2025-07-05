package main

func constantSpace() []int {
    arr := make([]int, 10)
    for i := 0; i < 10; i++ {
        arr[i] = i
    }
    return arr
}

func linearSpace(n int) []int {
    arr := make([]int, n)
    for i := 0; i < n; i++ {
        arr[i] = i
    }
    return arr
}

func linearAppend(n int) []int {
    var result []int
    for i := 0; i < n; i++ {
        result = append(result, i)
    }
    return result
}

func quadraticSpace(n int) [][]int {
    matrix := make([][]int, n)
    for i := 0; i < n; i++ {
        matrix[i] = make([]int, n)
    }
    return matrix
}

func allocationPerIteration(n int) [][]int {
    result := make([][]int, 0, n)
    for i := 0; i < n; i++ {
        row := make([]int, 10)
        result = append(result, row)
    }
    return result
}

func recursiveStack(n int) int {
    if n == 0 {
        return 0
    }
    return 1 + recursiveStack(n-1)
}

func tailRecursive(n, acc int) int {
    if n == 0 {
        return acc
    }
    return tailRecursive(n-1, acc+n)
}

func fixedLoop(n int) int {
    sum := 0
    for i := 0; i < 10; i++ {
        sum += i
    }
    return sum
}

func multiInputAllocation(n, m int) ([]int, []int) {
    a := make([]int, n)
    b := make([]int, m)
    return a, b
}

func mapSpace(n int) map[int]int {
    m := make(map[int]int)
    for i := 0; i < n; i++ {
        m[i] = i * i
    }
    return m
}

func reuseBuffer(n int) []int {
    buf := make([]int, n)
    for i := 0; i < n; i++ {
        buf[i] = i * 2
    }
    return buf
}

func conditionalAlloc(n int) []int {
    if n > 100 {
        return make([]int, n)
    }
    return nil
}

func fixedAlloc() []int {
	return make([]int, 10)
}

func recurAlloc(n int) int {
	if n == 0 {
		return 0
	}
	arr := make([]int, n)
	return arr[0] + recurAlloc(n-1)
}
