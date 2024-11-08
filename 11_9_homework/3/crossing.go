package main

import (
	"fmt"
)

func main() {
	str1 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str2 := "0123456789"
	ch1 := make(chan rune)
	ch2 := make(chan rune)
	len1 := len(str1)
	len2 := len(str2)

	go func() {
		for _, a := range str1 {
			ch1 <- a
		}
	}()

	go func() {
		for _, a := range str2 {
			ch2 <- a
		}
	}()

	count1 := 0
	count2 := 0

	for i := 0; i < len1+len2; i++ {
		b, ok1 := <-ch1
		if ok1 {
			fmt.Print(string(b))
		}
		count1++

		c, ok2 := <-ch2
		if ok2 {
			fmt.Print(string(c))
		}
		count2++

		if count1 == len1 {
			close(ch1)
		}

		if count2 == len2 {
			close(ch2)
		}
	}
}
