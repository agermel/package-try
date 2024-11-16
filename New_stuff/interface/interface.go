package main

import "fmt"

type Simpler interface {
	Get() int
	Set(a int)
}

type Simple struct {
	number int
}

func (s Simple) Get() int {
	return 1
}

func (s Simple) Set(a int) {
	s.number = a
}

func main() {
	var Intf Simpler = Simple{}
	Intf.Set(2)
	fmt.Println(Intf.Get())
}
