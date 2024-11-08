package main

import (
	"fmt"
	"sync"
	"time"
)

var goroutine_times int = 0

type message struct {
	Topic     string
	Partition int32
	Offset    int64
}

type FeedEventDM struct {
	Type    string
	UserID  int
	Title   string
	Content string
}

type MSG struct {
	ms        message
	feedEvent FeedEventDM
}

const ConsumeNum = 5

func main() {
	var consumeMSG []MSG
	var lastConsumeTime time.Time // 记录上次消费的时间
	var wait sync.WaitGroup
	msgs := make(chan MSG)

	//这里源源不断的生产信息
	go func() {
		for i := 0; ; i++ {
			msgs <- MSG{
				ms: message{
					Topic:     "消费主题",
					Partition: 0,
					Offset:    0,
				},
				feedEvent: FeedEventDM{
					Type:    "grade",
					UserID:  i,
					Title:   "成绩提醒",
					Content: "您的成绩是xxx",
				},
			}
			//每次发送信息会停止0.01秒以模拟真实的场景
			time.Sleep(100 * time.Millisecond)
		}
	}()

	//不断接受消息进行消费
	for msg := range msgs {
		// 添加新的值到events中
		consumeMSG = append(consumeMSG, msg)
		// 如果数量达到额定值就批量消费
		if len(consumeMSG) >= ConsumeNum {
			//进行异步消费
			wait.Add(1)
			go func() {
				defer wait.Done()
				m := consumeMSG[:ConsumeNum]
				fn(m)
				goroutine_times++
			}()
			// 更新上次消费时间
			lastConsumeTime = time.Now()
			wait.Wait()
			// 清除插入的数据
			consumeMSG = consumeMSG[ConsumeNum:]
		} else if !lastConsumeTime.IsZero() && time.Since(lastConsumeTime) > 5*time.Minute {
			// 如果距离上次消费已经超过5分钟且有未处理的消息
			if len(consumeMSG) > 0 {
				//进行异步消费
				wait.Add(1)
				go func() {
					defer wait.Done()
					m := consumeMSG[:ConsumeNum]
					fn(m)
					goroutine_times++
				}()
				// 更新上次消费时间
				lastConsumeTime = time.Now()
				wait.Wait()
				// 清空插入的数据
				consumeMSG = consumeMSG[ConsumeNum:]
			}
		}
		if goroutine_times == 10 {
			fmt.Println("Exit!")
			return
		}
	}
}

func fn(m []MSG) {
	fmt.Printf("本次消费了%d条消息\n", len(m))
}
