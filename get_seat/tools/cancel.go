package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Cancel(client *http.Client, rsvID string) {
	//eg_url = http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?act=del_resv&id=159768688&_=1732847236390
	urlCancel := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx"
	paramsCancel := "?act=del_resv&id=" + rsvID + "&_=" + fmt.Sprintf("%d", time.Now().UnixMilli())

	reqCancel, err := http.NewRequest("GET", urlCancel+paramsCancel, nil)
	if err != nil {
		log.Fatal(err)
	}

	respCancel, err := client.Do(reqCancel)
	if err != nil {
		log.Fatal(err)
	}
	defer respCancel.Body.Close()

	bodyCancel, err := ioutil.ReadAll(respCancel.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bodyCancel))
}
