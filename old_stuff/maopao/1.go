package main

import (
	"fmt"
)

func main() {
	var n, temp int
	var a []int
	fmt.Scanf("%d\n", &n)
	for i := 0; i < n; i++ {
		fmt.Scanf("%d", &temp)
		a = append(a, temp)
	}
	for i := n; i > 1; i-- { //几个数暂未定位
		for j := 0; j < i-1; j++ {
			if a[j] > a[j+1] {
				temp = a[j+1]
				a[j+1] = a[j]
				a[j] = temp
			}
		}
	}
	for i := 0; i < n; i++ {
		fmt.Printf("%d ", a[i])
	}
}
