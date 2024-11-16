package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.bilibili.com/")
	if err != nil {
		fmt.Println("Wrong:", err)
	}
	defer resp.Body.Close()

	resp1, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Wrong:", err)
	}

	resp2 := string(resp1)
	fmt.Println(resp2)
}
