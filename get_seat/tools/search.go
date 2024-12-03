package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Formattedtime struct {
	startTime_Uncoded      string
	overTime_Uncoded       string
	startTime_coded_url    string
	overTime_coded_url     string
	startTime_coded_string string
	overTime_coded_string  string
	formattedDate          string
	start_to_compare       time.Time
	over_to_compare        time.Time
}

type Reserve_need struct {
	DevID         string
	formattedtime Formattedtime
	client        *http.Client
	SeatID        string
}

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

type History struct {
	Message string `json:"msg"`
}

func format_time(startTime_Uncoded, overTime_Uncoded string) (formattedtime Formattedtime, err error) {
	currentTime := time.Now()

	// 获取今天的星期几
	weekday := currentTime.Weekday()

	// 开放时间：周一到周四，和周六到周日为 08:00 - 22:00
	openStart := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 8, 0, 0, 0, currentTime.Location()) // 08:00
	openEnd := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 22, 0, 0, 0, currentTime.Location())  // 22:00

	// 如果今天是星期五，开放时间改为 08:00 - 14:00
	if weekday == time.Friday {
		openEnd = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 14, 0, 0, 0, currentTime.Location()) // 14:00
	}

	// 如果当前时间不在开放时间内，返回错误
	if currentTime.Before(openStart) || currentTime.After(openEnd) {
		return formattedtime, fmt.Errorf("还没开馆哥们")
	}

	// 格式化为 "YYYY-MM-DD" 格式
	formattedtime.formattedDate = currentTime.Format("2006-01-02")

	// 格式化为 19%3A50 格式
	formattedtime.startTime_coded_url = url.QueryEscape(startTime_Uncoded)

	parsedTime_start, err := time.Parse("15:04", startTime_Uncoded)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return
	}

	// 格式化为2000格式
	formattedtime.startTime_coded_string = parsedTime_start.Format("1504") // 输出: 2100

	formattedtime.overTime_coded_url = url.QueryEscape(overTime_Uncoded)

	parsedTime_over, err := time.Parse("15:04", overTime_Uncoded)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return
	}

	// 格式化为2000格式
	formattedtime.overTime_coded_string = parsedTime_over.Format("1504") // 输出: 2100

	layout := "2006-01-02 15:04"
	formattedtime.start_to_compare, err = time.Parse(layout, formattedtime.formattedDate+" "+startTime_Uncoded)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return
	}

	formattedtime.over_to_compare, err = time.Parse(layout, formattedtime.formattedDate+" "+overTime_Uncoded)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return
	}

	return
}

func get_info(client *http.Client, startTime_Uncoded, overTime_Uncoded string, i int) (Reserve_need, error) {

	var reserve_need Reserve_need
	reserve_need.client = client

	// real_url http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx?byType=devcls&classkind=8&display=fp&md=d&room_id=101699187&purpose=&selectOpenAty=&cld_name=default&date=2024-11-26&fr_start=14%3A50&fr_end=17%3A50&act=get_rsv_sta&_=1732603671177
	sitUrl := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx"

	room_id := []string{"101699187", "101699179", "101699189", "101699191"}
	//room_id := "101699187" //南湖分馆一楼中庭开敞座位区     // N1001 ~ N1172 devID:101699605~101699776
	//room_id := "101699179" //南湖分馆一楼开敞座位区         // N1173 ~ N1328 devID:101699777~101699932
	//room_id := "101699189" //南湖分馆二楼开敞座位区		  // N2001 ~ N2148 devID:101699933~101700080
	//room_id := "101699191" //南湖分馆二楼卡座区		      // K2001 ~ K2096 devID:101700081~101700176

	// 格式化时间
	formattedtime, err := format_time(startTime_Uncoded, overTime_Uncoded)
	if err != nil {
		return reserve_need, err
	}

	// 构建查询参数
	params := "?byType=devcls&classkind=8&display=fp&md=d&room_id=" + room_id[i] + "&purpose=&selectOpenAty=&cld_name=default&date=" + formattedtime.formattedDate +
		"&fr_start=" + formattedtime.startTime_coded_url + "&fr_end=" + formattedtime.overTime_coded_url + "&act=get_rsv_sta&_=" +
		fmt.Sprintf("%d", time.Now().UnixMilli())

	var devID string
	// 查询请求
	reqSit, err := http.NewRequest("GET", sitUrl+params, nil)
	if err != nil {
		return reserve_need, fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送请求并处理响应
	respSit, err := client.Do(reqSit)
	if err != nil {
		return reserve_need, fmt.Errorf("请求失败: %w", err)
	}
	defer respSit.Body.Close()

	// 读取响应内容
	bodySir, err := ioutil.ReadAll(respSit.Body)
	if err != nil {
		return reserve_need, fmt.Errorf("读取响应失败: %w", err)
	}

	var response Response
	// 解析 JSON 数据
	err = json.Unmarshal(bodySir, &response)
	if err != nil {
		return reserve_need, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	var seatName string

	// 遍历座位信息
	for _, seat := range response.Data {
		if len(seat.TS) == 0 { // 如果 TS 为空
			devID = seat.DevID
			break
		} else {
			for _, ts := range seat.TS {
				layout := "2006-01-02 15:04"
				start, err := time.Parse(layout, ts.Start)
				if err != nil {
					return Reserve_need{}, fmt.Errorf("时间解析失败，开始时间解析错误: %w", err)
				}
				end, err := time.Parse(layout, ts.End)
				if err != nil {
					return Reserve_need{}, fmt.Errorf("时间解析失败，结束时间解析错误: %w", err)
				}
				// 检查时间区间是否符合
				if !(formattedtime.over_to_compare.Before(start) || formattedtime.start_to_compare.After(end)) {
					continue
				} else {
					devID = seat.DevID
					seatName = seat.DevName
					break
				}
			}
		}
	}

	reserve_need.formattedtime = formattedtime
	reserve_need.DevID = devID
	reserve_need.SeatID = seatName
	return reserve_need, nil
}

// 找啊找啊找座位
func All_rooms_reserve_find(client *http.Client, startTime_Uncoded, overTime_Uncoded string) (Reserve_need, error) {
	var reserve_need Reserve_need
	var err error
	for i := 0; i < 4; i++ {
		reserve_need, err = get_info(client, startTime_Uncoded, overTime_Uncoded, i)
		if reserve_need.DevID != "" || err != nil {
			return reserve_need, err
		}
	}
	return reserve_need, err
}

func Search_reserve_history(client *http.Client) (rsvID, seat string, err error) {
	urlHistory := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/center.aspx"
	paramsHistory := "?act=get_History_resv&strat=90&StatFlag=New&_=" + fmt.Sprintf("%d", time.Now().UnixMilli())

	reqHistory, err := http.NewRequest("GET", urlHistory+paramsHistory, nil)
	if err != nil {
		return
	}

	respHistory, err := client.Do(reqHistory)
	if err != nil {
		return
	}
	defer respHistory.Body.Close()

	bodyHistory, err := ioutil.ReadAll(respHistory.Body)
	if err != nil {
		return
	}

	var history History
	err = json.Unmarshal(bodyHistory, &history)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return
	}

	htmlContent := history.Message
	fmt.Println("完整 HTML 内容:", htmlContent) // 输出调试信息
	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return
	}

	// 查找含有 onclick 属性的 <a class="click"> 标签
	// 这里提取 onclick 属性的值
	onClick := doc.Find("a.click").AttrOr("onclick", "")

	// 使用正则表达式提取其中的 revID（即括号中的数字）
	re := regexp.MustCompile(`finish\((\d+)\)`)
	matches := re.FindStringSubmatch(onClick)

	seat = doc.Find(".box a").First().Text()

	if len(matches) < 2 {
		err = fmt.Errorf("哥们，你还没预定")
		return
	}

	if matches[1] != "" {
		// 提取到的 revID 是匹配到的第一个数字
		rsvID = matches[1]
		return
	} else {
		err = fmt.Errorf("哥们，你还没预定")
		return
	}

}
