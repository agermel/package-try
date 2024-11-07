package main

import (
	"fmt"
)

/*思路：桶式标记*/
func main() {
	var a [5][5]int
	var n, m, temp int
	for i := 0; i < 5; i++ { //r
		for j := 0; j < 5; j++ { //c
			fmt.Scan(&a[i][j])
		}
	}
	fmt.Scan(&n, &m)
	for i := 0; i < 5; i++ {
		temp = a[m-1][i]
		a[m-1][i] = a[n-1][i]
		a[n-1][i] = temp
	}
	for i := 0; i < 5; i++ { //c
		for j := 0; j < 5; j++ { //r
			fmt.Print(a[i][j], " ")
		}
		fmt.Println()
	}
}
