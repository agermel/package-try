package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	data := url.Values{}
	data.Set("name", "Alice")
	data.Add("name", "Reimu")
	data.Set("age", "25")

	res, err := http.PostForm("http://httpbin.org/post", data)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	fmt.Printf("%s", body)
}
