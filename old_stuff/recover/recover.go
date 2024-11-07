package main

import (
	"fmt"
)

/*
func division(a, b float64) float64 {
	if b == 0 {
		panic("why divisor zero")
	}
	return a / b
*/

/*
func division(a, b float64) (c float64, err error) {
	if b == 0 {
		return 0, errors.New("NOT ZERO")
	} else {
		return a / b, nil
	}
}

func main() {
	var err error
	var a, b = 1932.232, 0.0
	res, err := division(a, b)
	if err == nil {
		fmt.Println(res)
	} else {
		fmt.Println(err)
	}
}	*/

func test(e error) {
	defer func() {
		if p := recover(); p != nil { //直接在if语句中声明p
			e = fmt.Errorf("jaha %s", p)

		}
	}()
	panic("error")
}
