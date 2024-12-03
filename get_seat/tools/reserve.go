package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Reserve(reserve_need Reserve_need) {
	urlRese := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx"
	paramsRese := "?dialogid=&dev_id=" + reserve_need.DevID + "&lab_id=&kind_id=&room_id=&type=dev&prop=&test_id=&term=&Vnumber=&classkind=&test_name=&start=" +
		reserve_need.formattedtime.formattedDate + "+" + reserve_need.formattedtime.startTime_coded_url + "&end=" + reserve_need.formattedtime.formattedDate + "+" +
		reserve_need.formattedtime.overTime_coded_url + "&start_time=" + reserve_need.formattedtime.startTime_coded_string + "&end_time=" + reserve_need.formattedtime.overTime_coded_string +
		"&up_file=&memo=&act=set_resv&_=" + fmt.Sprintf("%d", time.Now().UnixMilli())

	reqRese, err := http.NewRequest("GET", urlRese+paramsRese, nil)
	if err != nil {
		log.Fatal(err)
	}

	respRese, err := reserve_need.client.Do(reqRese)
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
