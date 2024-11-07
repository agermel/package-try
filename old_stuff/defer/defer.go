package main

import (
	"fmt"
	"time"
)

func main() {
	var x [5]struct{}
	for i := range x {
		defer func(y int) {
			time.Sleep(1000000000)
			fmt.Println(y)
			time.Sleep(1000000000)
		}(i)
	}
}
