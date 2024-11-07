package main

import "fmt"

func main() {
	x := maxinum(10, 193, 2913, 232, 1)
	fmt.Println(x)
}

func maxinum(arr ...int) int {
	max := arr[0]
	for _, a := range arr {
		if a > max {
			max = a
		}
	}
	return max
}
