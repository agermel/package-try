package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	res, _ := reader.ReadString('\n')
	res = strings.ToLower(res)
	words := strings.Fields(res)
	newcnt := 0
	newchar := 'z'
	newlen := 9999999
	var haoword string
	for _, word := range words { //遍历句子
		if word == "I'm" {
			word = "i am"
		}
		if word == "I'll" {
			word = "i will"
		}
		var alpha []rune
		for _, b := range word { //遍历字母
			alpha = append(alpha, b)
		}
		cnt := 1
		for i := len(alpha) - 1; i > 0; i-- {
			if alpha[i] == alpha[i-1] {
				cnt++
			} else {
				if cnt > newcnt || (cnt == newcnt && alpha[i] < newchar || (alpha[i] == newchar && len(alpha) < newlen) || (len(alpha) == newlen)) {
					haoword = word
					newcnt = cnt
					newlen = len(alpha)
					newchar = alpha[i]
				}
				break
			}
		}
	}
	fmt.Print(haoword)
}

//第四点做不到
