package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

var data string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzE3NjU3MjEsInN1YiI6ImNvZGUifQ.ZbwoMu7y_prHmYc_JLHuTMiTCsOhEb-VVXlAISQEE2Y"

func checkpoint1() {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://muxithief.muxixyz.com/api/v1/login", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("code", "code")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res))
}

func checkpoint2() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://muxithief.muxixyz.com/api/v1/getStrategy", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", data)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res))
}

func checkpoint3() string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://muxithief.muxixyz.com/api/v1/eyes", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", data)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	res1 := strings.Split(string(res), "\"")

	image, err := base64.StdEncoding.DecodeString(res1[3])

	return string(image)
}

func main() {
	file, err := os.Open("src/test/5/code4.jpg")
	if err != nil {
		log.Fatal("打开图片文件失败:", err)
	}
	defer file.Close()

	client := &http.Client{}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fileWriter, err := writer.CreateFormFile("file", "code4.jpg")
	if err != nil {
		log.Fatal("创建文件字段失败:", err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Fatal("写入文件内容失败:", err)
	}

	// 关闭 multipart writer 以设置结束边界
	writer.Close()

	req, err := http.NewRequest("PUT", "http://muxithief.muxixyz.com/api/v1/attack", &body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", data)
	go func() {
		for i := 0; i < 10; i++ {
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			res, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(res))

			time.Sleep(time.Millisecond * 300)
		}
	}()
	time.Sleep(time.Second * 2)
	client1 := &http.Client{}

	req1, err := http.NewRequest("GET", "http://muxithief.muxixyz.com/api/v1/paper", nil)
	if err != nil {
		log.Fatal(err)
	}

	req1.Header.Add("Authorization", data)

	resp1, err := client1.Do(req1)
	if err != nil {
		log.Fatal(err)
	}
	defer resp1.Body.Close()

	res, err := io.ReadAll(resp1.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res))

}
