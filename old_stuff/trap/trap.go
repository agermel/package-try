package main

import "fmt"

type person struct {
	age  int
	name string
}

func (p *person) chan_age(a int) {
	p.age = a //你实际上是在操作这个指针所指向的 person 实例的内存位置
}

func main() {
	pp := person{11, "Jackson"}
	pp.chan_age(1993)
	fmt.Println(pp.age)
}

/*注意，在Go中，当你有一个类型的方法，并且这个方法的接收者是指针类型时，
你可以直接通过实例变量（而不是它的地址）来调用这个方法，因为Go会自动为你取地址。
这是Go语言设计中的一个方便之处。*/
