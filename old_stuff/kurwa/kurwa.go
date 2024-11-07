package main

import "fmt"

func get_average(arr []int, size int) float32 {
	var sum float32
	for i := 0; i < size; i++ {
		sum += float32(arr[i])

	}
	avg := sum / float32(size)
	return avg
}

func main() {
	arr := []int{1, 3, 4, 5, 3}
	avg := get_average(arr, 5)
	fmt.Println(avg)
}
