package main

import "fmt"

type porn struct {
	sex   string
	title string
}

func (c porn) nova() {
	fmt.Println("You have watched", c.sex, c.title)
}

func main() {
	var c porn
	c.sex = "man"
	c.title = "horrifying"
	c.nova()
}
