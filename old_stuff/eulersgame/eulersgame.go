package main

import "fmt"

func euler(n int) int {
	var prime []int     //记录素数
	var st [400000]bool //判断合数
	st[1] = true
	cnt := 0
	for i := 2; i <= 400000; i++ {
		if !st[i] {
			prime = append(prime, i)
			cnt++
		}
		for _, k := range prime {
			if k*i > 400000 {
				break
			}
			st[k*i] = true
			if i%k == 0 {
				break
			}
		}
		if cnt == n {
			return prime[n-1]
		}
	} //素数全填进了prime
	return 0
}

func main() {
	var n, m int
	fmt.Scanf("%d", &n)
	m = euler(n)
	fmt.Printf("%d", m)
}
