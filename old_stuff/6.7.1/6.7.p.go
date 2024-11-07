package main

import (
	"fmt"
	"unicode"
)

func main() {
	saying := "啊啊啊anno你带我走吧！！！！！"
	fmt.Print(nmap(modnonAscii, saying))
}

func modnonAscii(a rune) rune {
	if a > unicode.MaxASCII {
		a = '?'
	}
	return a
}

func nmap(f func(a rune) rune, s string) string {
	var new []rune
	for _, k := range s {
		new = append(new, f(k))
	}
	return string(new)
}
