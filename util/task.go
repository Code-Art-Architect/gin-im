package util

import (
	"fmt"
	"time"
)

type TimerFunc func(any) bool

func InitTimer() {
	Timer(time.Second*3, time.Second*6, clearConnection, "")
}

func clearConnection(params any) (ans bool) {
	ans = false
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("clean connections err: ", e)
		}
	}()
	fmt.Println("定时任务: 清理超时连接", params)
	ClearTimedOutConnections()
	return ans
}

func ClearTimedOutConnections() {
	fmt.Println("......")
}

// delay: 首次延迟
// tick: 间隔
// execute: 执行的函数
// params: 参数
func Timer(delay, interval time.Duration, execute TimerFunc, params any) {
	go func() {
		if execute == nil {
			return
		}
		timer := time.NewTimer(delay)
		for {
			select {
			case <-timer.C:
				if !execute(params) {
					return
				}
				timer.Reset(interval)
			}
		}
	}()
}
