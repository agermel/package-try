package main

import "fmt"

func sum(x ...int) int {
	summ := 0
	for _, valve := range x {
		summ += valve
	}
	return summ
}

func main() {
	j := []int{1, 2, 3, 4, 5}
	fmt.Println(sum(j...))
}
