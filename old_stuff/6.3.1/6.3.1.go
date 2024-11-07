package main

import "fmt"

func main() {
	print("aa", "bb", "cc", "dd", "ee")
}

func print(str ...string) {
	for _, a := range str {
		fmt.Println(a)
	}
}
