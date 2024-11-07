package main

import (
	"fmt"
	"time"
)

// 利用闭包来写斐波那契
func main() {
	start := time.Now()
	fina := fi()
	fmt.Println(fina(100))
	end := time.Now()
	sub := end.Sub(start)
	fmt.Print(sub)
}

func fi() func(int) int {
	former := 0
	latter := 1
	return func(n int) (ans int) {
		if n == 0 {
			return former
		}
		if n == 1 {
			return latter
		}
		for i := 2; i <= n; i++ {
			ans = former + latter
			former = latter
			latter = ans
		}
		return ans
	}
}

//第二种
/*
package main

import "fmt"

// 主函数
func main() {
	fibonacci := fibonacciGenerator() // 创建一个斐波那契数列的闭包

	// 打印前 10 个斐波那契数
	for i := 0; i < 10; i++ {
		fmt.Printf("Fibonacci(%d) = %d\n", i, fibonacci())
	}
}

// fibonacciGenerator 返回一个生成斐波那契数的闭包
func fibonacciGenerator() func() int {
	former, latter := 0, 1 // 初始化前两个斐波那契数
	return func() int {
		next := former // 保存当前斐波那契数
		// 更新斐波那契数
		former, latter = latter, former+latter
		return next // 返回当前斐波那契数
	}
}
*/

/*func main() {
	addcn := suffix_factory(".cn")
	fmt.Println(addcn("Hello"))
}

// func factory 添加各种各样的后缀
func suffix_factory(suffix string) func(s string) string {
	return func(name string) string {
		if !strings.HasSuffix(name, suffix) {
			return name + suffix
		}
		return name
	}
}
*/
