package main

import "fmt"

type Car struct {
	wheelCount int
}

type Mercedes struct {
	Car
}

func (car Car) count() {
	fmt.Println(car.wheelCount)
}

func (nicecar Mercedes) hello() {
	fmt.Println("hello")
}

func main() {
	nicecar := Mercedes{Car{4}}
	nicecar.Car.count() // nicecar.count 嵌套使子结构体的方法整合到父结构体来
	nicecar.hello()
}
