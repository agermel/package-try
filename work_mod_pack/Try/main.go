package main

import (
	"package1/examine"
	"package1/greetings"
)

func main() {
	a := "morning"
	if greetings.ISmorning(a) {
		greetings.Morning()
	}
	for i := 0; i <= 50; i++ {
		examine.Oddeven(i)
	}
}
