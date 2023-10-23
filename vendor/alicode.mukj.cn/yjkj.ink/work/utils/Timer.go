package utils

import (
	"time"
)

type Timer struct {
	loop bool //是否维循环
	hour int
	min int
	second int

	userData interface{}
	exit chan bool
	process func(interface{})
}

func (t *Timer)timer(seconds time.Duration) {

	defer func() {
		if !t.loop {
			t.exit <- true
		}
	}()
	timer := time.NewTicker(seconds * time.Second)
	for {
		select {
		case <-timer.C:
			t.process(t.userData)

			if t.loop {
				t.startTimer()
			}

			return
		}
	}

	return
}

func (t *Timer)startTimer() {
	now := time.Now()
	// 计算下一个零点
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), t.hour, t.min, t.second, 0, next.Location())

	second := time.Duration(next.Sub(now).Seconds())
	//var index float64 = 1
	//second := time.Duration(index)
	t.timer(second)
}

//定时在每天的某个时间点重复执行，比如每天早8点，每天晚8点（20）
func StartTimerLoopEveryDayTiming(hour ,min,second int,process func(interface{}),userdata interface{}) {

	t := &Timer{hour:hour,min:min,second:second,process:process,loop:true,exit:make(chan bool),userData:userdata}

	now := time.Now()
	// 计算下一个零点
	next := now
	next = time.Date(next.Year(), next.Month(), next.Day(), hour, min, second, 0, next.Location())

	if now.After(next) {
		next = now.Add(24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, min, second,0, next.Location())
	}

	second1 := time.Duration(next.Sub(now).Seconds())

	go func() {
		go t.timer(second1)
		select {
		
		}
	}()
}
//定时在下一个某个时间点重复执行，比如下一个早8点，下一个晚8点（20），只执行一次
func StartTimerNoLoopTiming(hour ,min,second int,process func(interface{}),userdata interface{}) {

	t := &Timer{hour:hour,min:min,second:second,process:process,loop:false,exit:make(chan bool),userData:userdata}

	now := time.Now()
	// 计算下一个零点
	next := now
	next = time.Date(next.Year(), next.Month(), next.Day(), hour, min, second, 0, next.Location())

	if now.After(next) {
		next = now.Add(24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, min, second,0, next.Location())
	}

	second1 := time.Duration(next.Sub(now).Seconds())

	go func() {
		go t.timer(second1)
		<-t.exit
	}()
}

//延迟多少秒，多少秒后循环执行，
func StartTimerLoopBySecond(delay,loop int,process func(interface{}),userdata interface{}) {

	time.Sleep(time.Duration(delay))

	timer := time.NewTicker(time.Duration(loop))
	for {
		select {
		case <-timer.C:
			process(userdata)
			return
		}
	}
}

//延迟多少秒执行一次，
func StartTimerNoLoopBySecond(delay int,process func(interface{}),userdata interface{}) {

	//time.Sleep(time.Duration(delay))

	timer := time.NewTicker(time.Duration(delay))
	for {
		select {
		case <-timer.C:
			process(userdata)
			return
		}
	}
}

