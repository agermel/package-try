package main

import (
	"fmt"
)

func main() {
	var r, c int
	fmt.Scan(&r, &c)
	square := 0
	a := make([][]int, r) //先前再后
	for i := range a {
		a[i] = make([]int, c) //正确定义二元切片的方法
	}
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Scan(&a[i][j])
		}
	}
	for i := 0; i < r-1; i++ {
		for j := 0; j < c-1; j++ {
			if a[i][j] == a[i+1][j] {
				square++
			}
		}
		fmt.Println()
	}
}
