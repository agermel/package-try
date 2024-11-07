package main

import (
	"fmt"
)

func test() (a error) {
	defer func() {
		if p := recover(); p != nil {
			a = fmt.Errorf("null : %s", p)
			fmt.Println(p)
			fmt.Println(a)
		}
	}()
	panic("error")
}

func main() {
	a := test()
	fmt.Println(a)
}
