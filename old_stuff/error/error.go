package main

import (
	"fmt"
)

func say(v int) {
	if v <= 0 {
		panic("For no reason")
	}
	fmt.Println("v = ", v)
}

func main() {
	var a, b int = -10, 100
	say(a)
	say(b)
}
