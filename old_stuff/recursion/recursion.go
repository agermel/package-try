package main

import "fmt"

func fff(n int) int {
	if n <= 2 {
		return n
	}
	return fff(n-1) + fff(n-2)
}

func main() {
	for i := 1; i <= 10; i++ {
		fmt.Printf("First:%d ", fff(i))
	}
}
