package main

import "fmt"

func jiec(n uint64) (res uint64) {
	if n > 1 {
		res = n * jiec(n-1)
		return res
	}
	return 1
}

func main() {
	fmt.Println(jiec(15))
}
