package main

import (
	"fmt"
)

func main() {
	var a [5][5]int
	d := [5]int{0, 0, 0, 0, 0}
	e := [5]int{10000, 10000, 10000, 10000, 10000}
	for i := 0; i < 5; i++ { //注意去掉5，i是行，j是列
		for j := 0; j < 5; j++ {
			fmt.Scan(&a[i][j])
			if a[i][j] >= d[i] {
				d[i] = a[i][j]
			}
		}
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if a[j][i] <= e[i] {
				e[i] = a[j][i]
			}
		}
	}
	p := false
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if d[i] == e[j] {
				fmt.Print(i+1, " ", j+1, " ", e[j])
				p = true
			}
		}
	}
	if !p {
		fmt.Print("not found")
	}
}
