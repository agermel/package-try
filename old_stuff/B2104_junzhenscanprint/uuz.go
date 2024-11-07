package main

import (
	"fmt"
)

func main() {
	var a, b, c [][]int
	var n, m, k int

	fmt.Scan(&n, &m, &k)
	a = make([][]int, n)
	for p := range a {
		a[p] = make([]int, m)
	}
	b = make([][]int, m)
	for p := range b {
		b[p] = make([]int, k)
	}
	c = make([][]int, n)
	for p := range c {
		c[p] = make([]int, k)
	}
	fmt.Print(a, b, c)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			fmt.Scan(&a[i][j])
		}
	}
	for i := 0; i < m; i++ {
		for j := 0; j < k; j++ {
			fmt.Scan(&b[i][j])
		}
	}
	fmt.Print(a, b)
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			for h := 0; h < m; h++ {
				c[i][j] += a[i][h] * b[h][j]
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			fmt.Printf("%d ", a[i][j])
		}
		fmt.Println()
	}
}
