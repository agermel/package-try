package main

import "fmt"

func main() {
	s1 := []string{"1", "12", "31", "4"}
	s2 := []string{"a", "s", "s", "g"}
	fmt.Println(Insert(s1, s2, 3))
}

func Insert(s1, s2 []string, n int) []string {
	//Newslice := s1[:n+1] + s2[:] + s1[n+1:]
	var newslice []string
	newslice = append(newslice, s1[:n]...)
	newslice = append(newslice, s2[:]...)
	newslice = append(newslice, s1[n:]...)
	return newslice
}
