package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type Seat struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	DevID    string `json:"devId"`
	DevName  string `json:"devName"`
	RoomName string `json:"roomName"`
	Ts       []TS   `json:"ts"`
	Ops      []Ops  `json:"ops"`
}

type TS struct {
	ID    string `json:"id"`
	Start string `json:"start"`
	End   string `json:"end"`
	State string `json:"state"`
	Name  string `json:"name"`
}

type Ops struct {
	Start string `json:"start"`
	End   string `json:"end"`
	State string `json:"state"`
}

func main() {
	// 初始化 CookieJar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("初始化 CookieJar 失败: %v", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// 登录URL
	loginURL := "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		log.Fatalf("创建 GET 请求失败: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	// 获取登录页面
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("获取登录页面失败: %v", err)
	}
	defer resp.Body.Close()

	// 提取动态参数
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}
	lt, execution, eventID := extractLoginParams(string(body))

	// 构造 POST 数据
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	if username == "" || password == "" {
		log.Fatal("用户名或密码未配置，请通过环境变量设置 USERNAME 和 PASSWORD")
	}

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("lt", lt)
	data.Set("execution", execution)
	data.Set("_eventId", eventID)
	data.Set("submit", "登录")

	// 发送登录 POST 请求
	payload := strings.NewReader(data.Encode())
	reqPost, err := http.NewRequest("POST", loginURL, payload)
	if err != nil {
		log.Fatalf("创建 POST 请求失败: %v", err)
	}
	reqPost.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	reqPost.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	respPost, err := client.Do(reqPost)
	if err != nil {
		log.Fatalf("登录请求失败: %v", err)
	}
	defer respPost.Body.Close()
	log.Println("登录成功")

	// 获取座位信息
	sitURL := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx"
	roomID := "101699179"
	currentDate := time.Now().Format("2006-01-02")

	params := url.Values{}
	params.Set("byType", "devcls")
	params.Set("classkind", "8")
	params.Set("display", "fp")
	params.Set("md", "d")
	params.Set("room_id", roomID)
	params.Set("cld_name", "default")
	params.Set("date", currentDate)
	params.Set("fr_start", "19:50")
	params.Set("fr_end", "20:50")
	params.Set("act", "get_rsv_sta")
	params.Set("_", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))

	finalURL := sitURL + "?" + params.Encode()
	reqSit, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		log.Fatalf("创建座位请求失败: %v", err)
	}
	reqSit.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	reqSit.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	reqSit.Header.Set("X-Requested-With", "XMLHttpRequest")

	respSit, err := client.Do(reqSit)
	if err != nil {
		log.Fatalf("获取座位信息失败: %v", err)
	}
	defer respSit.Body.Close()

	bodySit, err := io.ReadAll(respSit.Body)
	if err != nil {
		log.Fatalf("读取座位信息响应失败: %v", err)
	}

	// 解析座位信息
	var seats []Seat
	err = json.Unmarshal(bodySit, &seats)
	if err != nil {
		log.Fatalf("解析座位 JSON 数据失败: %v", err)
	}
	log.Printf("座位信息: %+v", seats)
}

// 提取登录页面中的动态参数
func extractLoginParams(body string) (string, string, string) {
	ltRegexp := regexp.MustCompile(`name="lt" value="([^"]+)"`)
	executionRegexp := regexp.MustCompile(`name="execution" value="([^"]+)"`)
	eventIDRegexp := regexp.MustCompile(`name="_eventId" value="([^"]+)"`)

	ltMatches := ltRegexp.FindStringSubmatch(body)
	executionMatches := executionRegexp.FindStringSubmatch(body)
	eventIDMatches := eventIDRegexp.FindStringSubmatch(body)

	if len(ltMatches) < 2 || len(executionMatches) < 2 || len(eventIDMatches) < 2 {
		log.Fatal("未找到登录页面的动态参数")
	}

	return ltMatches[1], executionMatches[1], eventIDMatches[1]
}
