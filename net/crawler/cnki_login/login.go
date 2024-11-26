package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	requestUrl := "https://login.cnki.net/TopLoginCore/api/loginapi/LoginPo"

	// 创建请求数据结构体
	payload := map[string]interface{}{
		"isAutoLogin": true,
		"p":           2,
		"pwd":         "wuhanis38..",
		"userName":    "17386253558",
	}

	// 将数据编码为 JSON 格式
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// 设置请求头
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Origin", "https://www.cnki.net")
	req.Header.Add("Referer", "https://www.cnki.net/")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")

	// 发送请求并获取响应
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 打印响应内容
	fmt.Println(string(body))
}
