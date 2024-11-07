package main

import "fmt"

func main() {
	print(10)
}

func print(n int) {
	if n > 0 {
		fmt.Println(n)
		n--
		print(n)
	} else {
		return
	}
}
