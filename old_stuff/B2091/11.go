package main

import (
	"fmt"
)

/*思路：桶式筛*/
func main() {
	var a []int
	var n, temp, j int
	j = 1
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		fmt.Scan(&temp)
		a = append(a, temp)
	}
	max := 0
	for i := 0; i < n-1; i++ {
		if a[i] == a[i+1] {
			j++
		} else {
			j = 1
		}
		if max < j {
			max = j
		}
	}
	fmt.Println(j)
}
