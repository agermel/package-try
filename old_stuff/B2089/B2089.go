package main

import (
	"fmt"
)

func main() {
	var n, temp int
	aa, bb, cc, dd := 0, 0, 0, 0
	var a []int
	fmt.Scanf("%d\n", &n)
	for i := 0; i < n; i++ {
		fmt.Scanf("%d", &temp)
		a = append(a, temp)
	}
	for i := 0; i < n; i++ {
		if a[i] <= 18 {
			aa++
		} else if a[i] <= 35 {
			bb++
		} else if a[i] <= 60 {
			cc++
		} else {
			dd++
		}
	}
	fmt.Printf("%.2f%%\n%.2f%%\n%.2f%%\n%.2f%%", float64(aa)/float64(n)*100, float64(bb)/float64(n)*100, float64(cc)/float64(n)*100, float64(dd)/float64(n)*100)
}
