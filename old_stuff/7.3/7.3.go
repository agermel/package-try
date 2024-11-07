package main

import "fmt"

func main() {
	a := []int{1, 24, 2, 4, 2, 5, 4, 645, 23}
	fmt.Println(sum(a...))
}

func sum(numbers ...int) int {
	ans := 0
	for _,number := range numbers {
		ans += number
	}
	return ans
}
