package main

import (
	"fmt"
	"net/http"
)

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter,r *http.Request){
    if r.URL.Path == "/Mygo" {
        fmt.Fprint(w,"欢迎来到扣扣空间")
    } else {
        fmt.Fprint(w,"错辣！")
    }
}

func main() {
    mymux := &MyMux{}
    http.ListenAndServe(":9090",mymux)
}