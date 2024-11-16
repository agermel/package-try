package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	url := "http://119.3.2.168:8080/database/info"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)
}
