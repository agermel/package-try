// 检查数字的特性
package examine

import "fmt"

//检查数字是奇是偶
func Oddeven(n int) {
	if n%2 == 0 {
		fmt.Println("Even!!!")
	} else {
		fmt.Println("Odd!!!")
	}
}
