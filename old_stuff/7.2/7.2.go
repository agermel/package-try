package main

import "fmt"

var arr [41]int

func main() {
	//1.先分配内存（创建代表数组）在填入
	//2.利用append的随时随地添加分配内存存入数字
	var slice []int = make([]int, 41)
	for i := 0; i <= 40; i++ {
		slice[i] = fi(i)
	}
	fmt.Println(slice)
}

func fi(n int) int {
	arr[1] = 1
	if arr[n] != 0 || n == 0 {
		return arr[n]
	}
	arr[n] = fi(n-1) + fi(n-2)
	return arr[n]
}
