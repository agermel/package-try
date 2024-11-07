package main

import (
	"fmt"
)

func main() {
	var temp, n, t int
	fmt.Scan(&t)
	var a string
	for j := 0; j < t; j++ {
		fmt.Scan(&n)
		var block []int //手链上的数
		for i := 0; i < n; i++ {
			fmt.Scan(&temp)
			if temp%2 == 0 { //双为0
				temp = 0
			} else {
				temp = 1
			}
			block = append(block, temp) //读入格子数据
		}
		del := 0
		for i := 1; i <= n-1; i++ {
			if block[i] == block[i-1] { //记录被删数数量
				del++
			}
		}
		if block[0] == block[n-1] { //首尾
			del++
		}
		if j == 0 {
			if del%2 == 0 {
				a = "cc"
			} else {
				a = "syj"
			} //删到最后者赢
		}
	}
	fmt.Print(a)
	/* a先 del=1
	   b  del=2
	   a  del=3 */

}
