package main

import (
	"fmt"
)

func countWays(n, m int) int {
	f := make([]int, m+1) //传m次在1手上
	g := make([]int, m+1) //传m次不在1手上
	f[0] = 1
	g[0] = 0
	for i := 1; i <= m; i++ {
		f[i] = g[i-1] * 2
		g[i] = f[i-1] + g[i-1]*(n-2)
	}
	return f[m]
}

func main() {
	var n, m int
	fmt.Scan(&n, &m)
	result := countWays(n, m)
	fmt.Println(result)
}
