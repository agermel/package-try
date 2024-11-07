package main

import "fmt"

func main() {
	n := 0
	fmt.Scan(&n)
	fmt.Print(fi(n), " ", n)
}
func fi(n int) (ans int) {
	if n < 3 {
		ans = 1
	} else {
		ans = fi(n-1) + fi(n-2)
	}
	return
}
