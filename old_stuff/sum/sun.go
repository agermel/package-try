package main

import "fmt"

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func main() {
	c := []int{1, 3, 4, 1, 5, -9}
	ch := make(chan int)
	go sum(c[:len(c)/2], ch)
	go sum(c[:], ch)
	x, y := <-ch, <-ch
	fmt.Println(x, y)
}
