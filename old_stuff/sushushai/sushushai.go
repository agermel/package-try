package main

import "fmt"

func euler(n int) int {
	var prime [30001]int //记录素数
	var st [30001]bool   //判断合数
	st[1] = true
	cnt := 0
	for i := 2; i <= 30001; i++ {
		if !st[i] {
			prime[cnt] = i
			cnt++
		}
		for _, k := range prime {
			st[k*i] = true
			if i%k == 0 {
				break
			}
		}
	} //素数全填进了prime
	return prime[n-1]
}

func main() {
	var n, m int
	fmt.Scanf("%d", &n)
	m = euler(n)
	fmt.Printf("%d", m)
}
