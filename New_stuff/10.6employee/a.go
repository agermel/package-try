package main

import "fmt"

type employee struct {
	salary int
}

func (sl *employee) giveRaise() {
	sl.salary *= 2 //要改变原有值就用指针形式
}

func main() {
	Ben := employee{salary: 1000}
	Ben.giveRaise()
	fmt.Println(Ben.salary)
}
