package main

import "fmt"

func main() {
	hello := "Hello World!"
	slice := []byte(hello)
	s1, s2 := divide(slice, 5)
	fmt.Printf("%s %s", s1, s2)
}

func divide(slice []byte, n int) ([]byte, []byte) {
	b1 := slice[:n+1]
	b2 := slice[n+1:]
	return b1, b2
}
