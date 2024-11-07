package main

import "fmt"

func main() {
	a, b := 10, 203
	ans := test(a, b)
	fmt.Println(ans)
}

func test(a, b int) (ans int) {
	ans = a + b
	defer fmt.Println("defer")
	fmt.Println("First")
	return
}
