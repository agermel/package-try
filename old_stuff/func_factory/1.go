package main

import "fmt"

func mult(a, b, c int) int {
	return a * b * c
}
func sum(a, b, c int) int {
	return a + b + c
}

type fff func(int, int, int) int

func house(a string, b []fff /*类型为函数的切片*/) fff {
	var x fff
	switch a {
	case "sum":
		x = b[0]
	case "mult":
		x = b[1]
	}
	return x
}

func main() {
	w := make([]fff, 2)
	w[0] = sum
	w[1] = mult
	a, b, c := 1, 2, 3
	fmt.Println(house("sum", w)(a, b, c))

}
