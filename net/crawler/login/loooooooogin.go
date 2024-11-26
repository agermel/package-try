package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type User struct {
	Name      string `json:"name"`
	StudentId string `json:"Pid"`
	Grade     string
}

var maxConcurrentRequests = 20 // 限制最大并发请求数
var wg sync.WaitGroup
var mutex sync.Mutex

func main() {
	// 初始化 cookie jar
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
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	// 模拟浏览器请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取HTML内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
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
	data.Set("username", "2024214744")  // 请根据需求修改为实际的用户名
	data.Set("password", "WEEdru1351.") // 请根据需求修改为实际的密码
	data.Set("lt", lt)
	data.Set("execution", execution)
	data.Set("_eventId", eventID)
	data.Set("submit", "登录")

	payload := strings.NewReader(data.Encode())

	// 发送POST请求进行登录
	reqPost, err := http.NewRequest("POST", loginURL, payload)
	if err != nil {
		fmt.Println("请求错误:", err)
		return
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
		return
	}
	defer respPost.Body.Close()

	// 控制最大并发数量
	ch := make(chan bool, maxConcurrentRequests)

	// 存储爬取的用户数据
	var allUsers []User

	// 生成学号范围
	var studentIds []string
	for i := 2024214000; i <= 2024215000; i++ {
		studentIds = append(studentIds, fmt.Sprintf("%d", i))
	}

	// 分批处理每200个学号
	for i := 0; i < len(studentIds); i += 200 {
		end := i + 200
		if end > len(studentIds) {
			end = len(studentIds)
		}
		// 获取每批次学号
		batchIds := studentIds[i:end]

		wg.Add(1)
		ch <- true // 控制并发

		go func(batchIds []string) {
			defer wg.Done()

			// 发起请求进行数据爬取
			for _, studentId := range batchIds {
				// 构建 GET 请求，替换学生学号参数
				searchURL := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx"
				params := "?type=logonname&ReservaApply=ReservaApply&term=" + studentId + "&_=" + fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

				// 发起GET请求搜索
				reqSearch, err := http.NewRequest("GET", searchURL+params, nil)
				if err != nil {
					fmt.Println("请求错误:", err)
					continue
				}

				respSearch, err := client.Do(reqSearch)
				if err != nil {
					fmt.Println("运行错误:", err)
					continue
				}

				bodySearch, err := ioutil.ReadAll(respSearch.Body)
				if err != nil {
					fmt.Println("运行错误:", err)
					continue
				}

				var Users []User
				// 解析 JSON 数据
				err = json.Unmarshal(bodySearch, &Users)
				if err != nil {
					fmt.Println("解析 JSON 失败:", err)
					continue
				}

				// 提取并打印每个学生的姓名
				for _, user := range Users {
					if strings.HasPrefix(user.StudentId, "2024") {
						user.Grade = "小登"
					}
					mutex.Lock()
					allUsers = append(allUsers, user)
					mutex.Unlock()
				}
			}
			fmt.Println("over")
			<-ch // 完成后释放信号
		}(batchIds)
	}

	// 等待所有请求完成
	wg.Wait()

	// 保存数据到 JSON 文件
	saveToJSON(allUsers)

	// 保存数据到 CSV 文件
	saveToCSV(allUsers)
}

func saveToJSON(users []User) {
	file, err := os.Create("users.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	file.Write(data)
	fmt.Println("数据已保存到 users.json")
}

func saveToCSV(users []User) {
	file, err := os.Create("users.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入 CSV 头
	writer.Write([]string{"Name", "StudentId", "Grade"})

	// 写入数据
	for _, user := range users {
		writer.Write([]string{user.Name, user.StudentId, user.Grade})
	}
	fmt.Println("数据已保存到 users.csv")
}
