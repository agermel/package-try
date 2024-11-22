package main

import (
	"fmt"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count,Count1 int

func counter(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Request Method:", r.Method) // 打印请求方法
    fmt.Println("before", count)
    mu.Lock()
    fmt.Fprintf(w, "Count1 %d\n", count)
    count++
    mu.Unlock()
    fmt.Println("after", count)
}

/*func normal(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Request Method:", r.Method) // 打印请求方法
    fmt.Println("before", Count1)
    mu.Lock()
    fmt.Fprintf(w, "Count2 %d\n", Count1)
    Count1++
    mu.Unlock()
    fmt.Println("after", Count1)
}*/


func main() {
	http.HandleFunc("/",counter)
    //http.HandleFunc("/momoka",normal)
	http.ListenAndServe(":8080",nil)
}