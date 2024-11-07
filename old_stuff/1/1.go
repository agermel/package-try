package main

import (
	"fmt"
)

func jump(nums []int) bool {
	var ok bool
	used := make([]bool, len(nums))
	dist := 0 //第几个格子
	ok = false
	for i := 0; !used[dist]; i++ {
		used[dist] = true //标记走过的点
		dist += nums[dist]
		if dist >= len(nums)-1 {
			ok = true
			break
		}
	}
	return ok
}

func main() {
	var block []int
	var temp, n int
	var ok bool
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		fmt.Scan(&temp)
		block = append(block, temp) //读入格子数据
	}
	ok = jump(block)
	if ok {
		fmt.Print("Yes")
	} else {
		fmt.Print("No")
	}
}
