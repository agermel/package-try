package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Response struct {
	Ret  int    `json:"ret"`
	Act  string `json:"act"`
	Msg  string `json:"msg"`
	Data []Seat `json:"data"`
}

type Seat struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	DevID    string `json:"devId"`
	DevName  string `json:"devName"`
	KindName string `json:"kindName"`
	RoomName string `json:"roomName"`
	FreeTime int    `json:"freeTime"`
	TS       []TS   `json:"ts"`
	Ops      []Ops  `json:"ops"`
}

type TS struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	State  string `json:"state"`
	Title  string `json:"title"`
	Owner  string `json:"owner"`
	AccNo  string `json:"accno"`
	Occupy bool   `json:"occupy"`
}

type Ops struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	State  string `json:"state"`
	Occupy bool   `json:"occupy"`
}

func main() {
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
	fmt.Printf("%+v\n", respPost.Body)

	//至此登录过程结束
	sitUrl := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx"

	//room_id := "101699187" //南湖分馆一楼中庭开敞座位区     // N1001 ~ N1172 devID:101699605~101699776
	room_id := "101699179" //南湖分馆一楼开敞座位区         // N1173 ~ N1328 devID:101699777~101699932
	//room_id := "101699189" //南湖分馆二楼开敞座位区		  // N2001 ~ N2148 devID:101699933~101700080
	//room_id := "101699191" //南湖分馆二楼卡座区		      // K2001 ~ K2096 devID:101700081~101700176

	currentTime := time.Now()

	// 格式化为 "YYYY-MM-DD" 格式
	formattedDate := currentTime.Format("2006-01-02")

	// 格式化为 19%3A50 格式
	// timeString := "21:00" 应输入的格式
	startTime_Uncoded := "20:00"
	startTime_coded_url := url.QueryEscape(startTime_Uncoded)

	parsedTime_start, err := time.Parse("15:04", startTime_Uncoded)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return
	}

	// 格式化为2000格式
	startTime_coded_string := parsedTime_start.Format("1504") // 输出: 2100

	overTime_Uncoded := "21:00"
	overTime_coded_url := url.QueryEscape(overTime_Uncoded)

	parsedTime_over, err := time.Parse("15:04", startTime_Uncoded)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return
	}

	// 格式化为2000格式
	overTime_coded_string := parsedTime_over.Format("1504") // 输出: 2100

	layout := "2006-01-02 15:04"
	start_to_compare, err := time.Parse(layout, formattedDate+" "+startTime_Uncoded)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return
	}

	over_to_compare, err := time.Parse(layout, formattedDate+" "+overTime_Uncoded)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return
	}

	// real_url http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx?byType=devcls&classkind=8&display=fp&md=d&room_id=101699187&purpose=&selectOpenAty=&cld_name=default&date=2024-11-26&fr_start=14%3A50&fr_end=17%3A50&act=get_rsv_sta&_=1732603671177

	params := "?byType=devcls&classkind=8&display=fp&md=d&room_id=" + room_id + "&purpose=&selectOpenAty=&cld_name=default&date=" + formattedDate +
		"&fr_start=" + startTime_coded_url + "&fr_end=" + overTime_coded_url + "&act=get_rsv_sta&_=" +
		fmt.Sprintf("%d", time.Now().UnixMilli())

	//查询过程
	reqSit, err := http.NewRequest("GET", sitUrl+params, nil)
	if err != nil {
		fmt.Println("请求错误:", err)
		return
	}

	respSit, err := client.Do(reqSit)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取HTML内容
	bodySir, err := ioutil.ReadAll(respSit.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	var response Response
	//解析 JSON 数据
	err = json.Unmarshal(bodySir, &response)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return
	}

	var devID string

	for _, seat := range response.Data {
		if len(seat.TS) == 0 { // 如果 TS 为空
			devID = seat.DevID
			break
		} else {
			for _, ts := range seat.TS {
				layout := "2006-01-02 15:04"
				start, err := time.Parse(layout, ts.Start)
				if err != nil {
					fmt.Println("时间解析失败:", err)
					return
				}
				end, err := time.Parse(layout, ts.End)
				if err != nil {
					fmt.Println("时间解析失败:", err)
					return
				}
				if !(over_to_compare.Before(start) || start_to_compare.After(end)) {
					continue
				} else {
					devID = seat.DevID
					break
				}
			}
		}
	}

	//预约过程

	// rese_url :http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?dialogid=&dev_id=101699827&lab_id=&kind_id=&room_id=&type=dev&prop=&test_id=&term=&Vnumber=&classkind=&test_name=&start=2024-11-26+21%3A00&end=2024-11-26+22%3A00&start_time=2100&end_time=2200&up_file=&memo=&act=set_resv&_=1732623358349

	urlRese := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx"
	paramsRese := "?dialogid=&dev_id=" + devID + "&lab_id=&kind_id=&room_id=&type=dev&prop=&test_id=&term=&Vnumber=&classkind=&test_name=&start=" + formattedDate +
		"+" + startTime_coded_url + "&end=" + formattedDate + "+" + overTime_coded_url + "&start_time=" + startTime_coded_string + "&end_time=" + overTime_coded_string +
		"&up_file=&memo=&act=set_resv&_=" + fmt.Sprintf("%d", time.Now().UnixMilli())

	reqRese, err := http.NewRequest("GET", urlRese+paramsRese, nil)
	if err != nil {
		log.Fatal(err)
	}

	respRese, err := client.Do(reqRese)
	if err != nil {
		log.Fatal(err)
	}
	defer respRese.Body.Close()

	bodyRese, err := ioutil.ReadAll(respRese.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bodyRese))
}
