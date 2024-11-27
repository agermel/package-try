package main

import (
	"get_seat/tools"
	"log"
)

func main() {
	username := "2024214744"
	password := "WEEdru1351."

	client, err := tools.Login(username, password)
	if err != nil {
		log.Fatal(err)
	}

	startTime_Uncoded := "20:00"
	overTime_Uncoded := "21:00"

	reserve_need := tools.All_rooms_reserve_find(*client, startTime_Uncoded, overTime_Uncoded)

	tools.Reserve(reserve_need)

}
