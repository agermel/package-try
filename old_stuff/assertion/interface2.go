package main

import "fmt"

func main() {
	var inter interface{} = "haren"
	str, ok := inter.(string)
	fmt.Println(str, ok)
}
