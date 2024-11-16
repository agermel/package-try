package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup

func main() {
	waitGroup.Add(1)
	// 传入上下文
	go Do(context.WithValue(context.Background(), 1, 2))
	waitGroup.Wait()
}

func Do(ctx context.Context) {
	// 新建定时器
	ticker := time.NewTimer(time.Second)
	defer waitGroup.Done()
	for {
		select {
		case <-ctx.Done(): // 永远也不会执行
		case <-ticker.C:
			fmt.Println("timeout")
			return
		default:
			fmt.Println(ctx.Value(1))
		}
		time.Sleep(time.Millisecond * 100)
	}
}
