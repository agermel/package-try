package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/anaskhan96/soup"
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

	// 下面主要是解析标签
	doc := soup.HTMLParse(content)

	subDocs := doc.FindAll("div", "class", "home-book-list-title")
	for _, subDoc := range subDocs {
		fmt.Println(strings.TrimSpace(subDoc.Text()))
	}
}
