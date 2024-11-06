package fibo

// 包级别的全局变量，存储最后输入的值!!! 不被函数所束缚
var LastInput int

// 计算斐波那契数列的函数
func Fibonacci(n int) int {
	LastInput = n
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return Fibonacci(n-1) + Fibonacci(n-2)
	}
}

// 输出最后输入的值
func GetLastInput() int {
	return LastInput
}
