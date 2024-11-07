package main

import (
	"fmt"
)

func main() {
	var r, c int
	sum := 0
	fmt.Scan(&r, &c)
	a := make([][]int, r) //先前再后
	for i := range a {
		a[i] = make([]int, c) //正确定义二元切片的方法
	}
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Scan(&a[i][j])
			if i == 0 || j == 0 || i == r-1 || j == c-1 {
				sum += a[i][j]
			} //无重复因为每个数只扫描一次
		}
	}
	fmt.Print(sum)
}
