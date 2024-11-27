package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
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

func format_time(startTime_Uncoded, overTime_Uncoded string) (formattedtime Formattedtime, err error) {
	currentTime := time.Now()

	// 格式化为 "YYYY-MM-DD" 格式
	formattedtime.formattedDate = currentTime.Format("2006-01-02")

	// 格式化为 19%3A50 格式
	// timeString := "21:00" 应输入的格式

	formattedtime.startTime_coded_url = url.QueryEscape(startTime_Uncoded)

	parsedTime_start, err := time.Parse("15:04", startTime_Uncoded)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return
	}

	// 格式化为2000格式
	formattedtime.startTime_coded_string = parsedTime_start.Format("1504") // 输出: 2100

	formattedtime.overTime_coded_url = url.QueryEscape(overTime_Uncoded)

	parsedTime_over, err := time.Parse("15:04", startTime_Uncoded)
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

func Get_info(client *http.Client, startTime_Uncoded, overTime_Uncoded string, i int) (reserve_need Reserve_need) {
	sitUrl := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx"

	room_id := []string{"101699187", "101699179", "101699189", "101699191"}
	//room_id := "101699187" //南湖分馆一楼中庭开敞座位区     // N1001 ~ N1172 devID:101699605~101699776
	//room_id := "101699179" //南湖分馆一楼开敞座位区         // N1173 ~ N1328 devID:101699777~101699932
	//room_id := "101699189" //南湖分馆二楼开敞座位区		  // N2001 ~ N2148 devID:101699933~101700080
	//room_id := "101699191" //南湖分馆二楼卡座区		      // K2001 ~ K2096 devID:101700081~101700176

	formattedtime, err := format_time(startTime_Uncoded, overTime_Uncoded)
	if err != nil {
		log.Fatal(err)
	}

	// real_url http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx?byType=devcls&classkind=8&display=fp&md=d&room_id=101699187&purpose=&selectOpenAty=&cld_name=default&date=2024-11-26&fr_start=14%3A50&fr_end=17%3A50&act=get_rsv_sta&_=1732603671177

	params := "?byType=devcls&classkind=8&display=fp&md=d&room_id=" + room_id[i] + "&purpose=&selectOpenAty=&cld_name=default&date=" + formattedtime.formattedDate +
		"&fr_start=" + formattedtime.startTime_coded_url + "&fr_end=" + formattedtime.overTime_coded_url + "&act=get_rsv_sta&_=" +
		fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	var devID string
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
	defer respSit.Body.Close()

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
				if !(formattedtime.over_to_compare.Before(start) || formattedtime.start_to_compare.After(end)) {
					continue
				} else {
					devID = seat.DevID
					break
				}
			}
		}
	}
	reserve_need.client = client
	reserve_need.DevID = devID
	reserve_need.formattedtime = formattedtime
	return
}

// 一直找下去直到找到id为止
func All_rooms_reserve_find(client http.Client, startTime_Uncoded, overTime_Uncoded string) (reserve_need Reserve_need) {
	for i := 0; i < 4; i++ {
		reserve_need = Get_info(&client, startTime_Uncoded, overTime_Uncoded, i)
		if reserve_need.DevID != "" {
			break
		}
	}
	if reserve_need.DevID == "" {
		time.Sleep(time.Minute * 15)
		All_rooms_reserve_find(client, startTime_Uncoded, overTime_Uncoded)
	}
	return
}
