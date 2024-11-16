package main

import (
	"fmt"
	"sort"
)

func main() {
	var sleep_time, quest int
	fmt.Scanf("%d %d", &sleep_time, &quest)
	wake_num := make([]int, sleep_time)
	quest_num := make([]int, quest)

	var num int
	for i := 0; i < sleep_time; i++ {
		fmt.Scan(&num)
		wake_num[i] = num
	}

	for i := 0; i < quest; i++ {
		fmt.Scan(&num)
		quest_num[i] = num
	}

	sort.Ints(wake_num)

	for j := 0; j < quest; j++ {
		sum := 0
		count := 0
		for i := sleep_time - 1; i >= 0; i-- {
			sum += wake_num[i]
			count++
			if sum >= quest_num[j] {
				fmt.Println(count)
				break
			}
			if i == 0 && sum < quest_num[j] {
				fmt.Println("-1")
				break
			}
		}
	}
}
