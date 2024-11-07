package main

import (
	"fmt"
)

func main() {
	var r, c int
	sum := 0
	fmt.Scan(&r, &c)
	// 正确初始化二维切片
	a := make([][]int, r)
	for i := range a {
		a[i] = make([]int, c)
	}
	// 读取矩阵元素
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Scan(&a[i][j]) // 注意这里需要使用&来获取值
			// 判断边界条件并累加
			if (i == 0 || j == 0 || i == r-1 || j == c-1) && !((i == 0 && j == 0) || (i == r-1 && j == 0) || (i == 0 && j == c-1) || (i == r-1 && j == c-1)) {
				sum += a[i][j]
			}
		}
	}
	fmt.Println(sum) // 使用Println来输出，这样输出后会换行，更整洁
}
