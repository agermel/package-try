package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Login(username, password string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// 获取登录页面HTML
	loginURL := "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		fmt.Println("请求错误:", err)
		return client, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	// 模拟浏览器请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return client, err
	}
	defer resp.Body.Close()

	// 读取HTML内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return client, err
	}

	// 使用正则表达式提取动态参数lt、execution、_eventId
	ltRegexp := regexp.MustCompile(`name="lt" value="([^"]+)"`)
	executionRegexp := regexp.MustCompile(`name="execution" value="([^"]+)"`)
	eventIDRegexp := regexp.MustCompile(`name="_eventId" value="([^"]+)"`)

	ltMatches := ltRegexp.FindStringSubmatch(string(body))
	executionMatches := executionRegexp.FindStringSubmatch(string(body))
	eventIDMatches := eventIDRegexp.FindStringSubmatch(string(body))

	// 提取到的动态参数
	lt := ltMatches[1]
	execution := executionMatches[1]
	eventID := eventIDMatches[1]

	// 准备POST请求数据
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("lt", lt)
	data.Set("execution", execution)
	data.Set("_eventId", eventID)
	data.Set("submit", "登录")

	payload := strings.NewReader(data.Encode())

	// 发送POST请求进行登录
	reqPost, err := http.NewRequest("POST", loginURL, payload)
	if err != nil {
		fmt.Println("请求错误:", err)
		return client, err
	}

	// 设置请求头
	reqPost.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	reqPost.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	reqPost.Header.Add("Origin", "https://account.ccnu.edu.cn")
	reqPost.Header.Add("Referer", "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page=")
	reqPost.Header.Add("Sec-Fetch-Mode", "navigate")
	reqPost.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	// 发起POST请求
	respPost, err := client.Do(reqPost)
	if err != nil {
		fmt.Println("请求失败:", err)
		return client, err
	}
	defer respPost.Body.Close()

	bodyPost, _ := ioutil.ReadAll(respPost.Body)

	fmt.Println(string(bodyPost))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyPost)))
	if err != nil {
		log.Fatal(err)
	}

	msg := doc.Find("#msg").Text()
	if msg == "您输入的用户名或密码有误。" {
		return client, fmt.Errorf("%s", "您输入的用户名或密码有误")
	}
	return client, nil
}
