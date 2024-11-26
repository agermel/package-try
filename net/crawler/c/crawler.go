package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	requestUrl := "https://www.liaoxuefeng.com/"
	// 发送Get请求
	rsp, err := http.Get(requestUrl)
	if err != nil {
		log.Println(err.Error())
		return
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	content := string(body)
	defer rsp.Body.Close()

	fmt.Println(content)
}
