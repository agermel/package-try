package main

import (
	"fmt"
	"regexp"
)

func main() {
	str := "请联系 support@example.com 或 sales@example.org 获取更多信息。另一个联系人的邮箱是 john.doe123@example.co.uk。"
	pat := "[0-9a-zA-Z.]+@[0-9a-zA-Z.]+"
	if ok, _ := regexp.MatchString(pat, str); ok {
		fmt.Println("Find!")
	}
	re, _ := regexp.Compile(pat) //翻译pat代表的搜索模式？是的// 编译正则表达式
	matches := re.FindAllString(str, -1)
	for _, email := range matches {
		fmt.Println(email)
	}
}
