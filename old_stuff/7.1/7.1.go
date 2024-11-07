package main

import "fmt"

var arr [41]int

func main() {
	//内容缓存加速斐波那契
	for i := 0; i <= 40; i++ {
		fmt.Println(fi(i))
	}
}

func fi(n int) int {
	arr[1] = 1
	if arr[n] != 0 || n == 0 {
		return arr[n]
	}
	arr[n] = fi(n-1) + fi(n-2)
	return arr[n]
}
