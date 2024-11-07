package main

import (
	"fmt"
	"strconv"
)

func main() {
	var a string = "1"
	var b string = "jj"
	aa, err1 := strconv.Atoi(a)
	bb, err2 := strconv.Atoi(b)
	fmt.Println(aa, err1, bb, err2)
}
