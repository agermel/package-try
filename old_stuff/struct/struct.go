package main

import "fmt"

type Books struct {
	title   string
	author  string
	subject string
	book_id int
}

func main() {
	book1 := Books{"Go 语言", "www.runoob.com", "Go 语言教程", 6495407}
	fmt.Println(book1)
}
