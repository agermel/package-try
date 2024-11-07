package main

import "fmt"

func main() {
	text1 := "Hello," +
		"world"
	fmt.Print(text1, " ")
	for i := 0; i < len(text1); i++ {
		fmt.Print(string(text1[i]))
	}
}
