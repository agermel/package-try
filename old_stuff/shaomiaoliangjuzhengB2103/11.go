package main

import (
	"fmt"
)

func main() {
	var a, b [][]int
	var m, n int
	fmt.Scan(&m, &n)
	a = make([][]int, m)
	for p := range a {
		a[p] = make([]int, n)
	}
	b = make([][]int, m)
	for p := range a {
		b[p] = make([]int, n)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			fmt.Scan(&a[i][j])
		}
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			fmt.Scan(&b[i][j])
		}
	}
	d := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if a[i][j] == b[i][j] {
				d++
			}
		}
	}
	fmt.Printf("%.2f", float64(d)/float64(m*n)*100)
}
