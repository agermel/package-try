package main

import "fmt"

func main() {
	callback(1, add)
}

func add(a, b int) {
	fmt.Print(a + b)
}

func callback(y int, x func(int, int)) {
	x(y, 2)
}
