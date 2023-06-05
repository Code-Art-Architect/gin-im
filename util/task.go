package util

import (
	"time"
)

type TimerFunc func(any) bool

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
