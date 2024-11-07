package main

import "fmt"

// 函数接受一个数组作为参数
func modifyArray(arr [5]int) {
	for i := 0; i < len(arr); i++ {
		arr[i] = arr[i] * 2
	}
	fmt.Println(arr)
}

// 函数接受一个数组的指针作为参数
func modifyArrayWithPointer(arr *[5]int) {
	for i := 0; i < len(*arr); i++ {
		(*arr)[i] = (*arr)[i] * 2
	}
}

func main() {
	// 创建一个包含5个元素的整数数组
	myArray := [5]int{1, 2, 3, 4, 5}

	fmt.Println("Original Array:", myArray)

	// 传递数组给函数，但不会修改原始数组的值
	modifyArray(myArray)
	fmt.Println("Array after modifyArray:", myArray)

	// 传递数组的指针给函数，可以修改原始数组的值
	modifyArrayWithPointer(&myArray)
	fmt.Println("Array after modifyArrayWithPointer:", myArray)
}
