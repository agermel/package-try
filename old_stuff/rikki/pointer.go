package main

import "fmt"

const max int = 3

var ptr [max]*int

func main() {
	a := []int{100, 102, 103}
	for i := 0; i < len(a); i++ {
		fmt.Printf("a[%d] = [%d] ", i, a[i])
		ptr[i] = &a[i]
	}
	fmt.Println(ptr)
}
