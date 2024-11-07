package main

import "fmt"

func main() {
	n := [][]int{}
	row1 := []int{3, 4, 5}
	row2 := []int{6, 6, 6}
	n = append(n, row1) //这里在赋值给n
	n = append(n, row2)
	for i := 0; i <= 1; i++ {
		for j := 0; j <= 2; j++ {
			fmt.Print(n[i][j], " ")
		}
	}
}
