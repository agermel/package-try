package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 替换为你提供的网址
	url := "http://gtainmuxi.muxixyz.com/api/v1/organization/code"

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 输出响应的所有头部信息
	fmt.Println("响应头:")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
			// 检查是否有名为 passport 的字段
			if key == "passport" {
				fmt.Printf("\n找到passport: %s\n", value)
			}
		}
	}

	// 如果没有直接找到 passport 字段，可以检查常见的认证相关字段
	fmt.Println("\n如果没有找到passport，检查其他常见的认证相关字段:")
	// 比如 Authorization、Set-Cookie、X-Request-Id 等字段
	authorization := resp.Header.Get("Authorization")
	if authorization != "" {
		fmt.Printf("Authorization: %s\n", authorization)
	}

	// 或者检查是否有 Set-Cookie 字段
	cookies := resp.Header["Set-Cookie"]
	if len(cookies) > 0 {
		fmt.Println("Set-Cookie:")
		for _, cookie := range cookies {
			fmt.Println(cookie)
		}
	}
}
