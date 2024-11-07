package main

import "fmt"

type a func(float64, float64) float64

func size(a, b float64) float64 {
	return a * b
}
func area(a, b float64, c a) float64 { //c代表着使用哪种变量
	return c(a, b) //求面积通用方法
}
func main() {
	l, v := 11.2, 23.4
	rect := area(l, v, size)
	fmt.Println(rect)
}
