package main

import (
	"fmt"
)

func abs(a int) int {
	if a < 0 {
		a = -a
	}
	return a
}

func main() {
	var n int
	fmt.Scanln(&n)
	arr := make([]int, n)

	for i := 0; i < n; i++ {
		fmt.Scanf("%d", &arr[i])
	}
	fmt.Println(arr)

	sum, count_plus, count_minus := 0, 0, 0
	var temp int
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			temp = abs(arr[i] - arr[j])
			if temp >= 2 {
				count_minus++
			}
			if temp <= 0 {
				count_plus++
			}
			sum += temp
		}
	}

	if count_plus >= count_minus {
		fmt.Println(sum)
	}
	if count_plus < count_minus && sum%2 == 1 {
		fmt.Println("1")
	}
	if count_plus < count_minus && sum%2 == 0 {
		fmt.Println("0")
	}
	fmt.Println(sum)
}
