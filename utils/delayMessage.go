package utils

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type TimePolling struct {
	curIndex  int
	slots     [3600]map[string]*Task
	sTime     time.Time
	mu        sync.Mutex
	next      chan bool
	closed    chan bool
	taskClose chan bool
	timeClose chan bool
}

type TaskFunc func(args ...any)

type Task struct {
	cycleNum int
	method   TaskFunc
	params   []any
}

func NewPolling() (polling *TimePolling) {
	polling = &TimePolling{
		curIndex:  0,
		sTime:     time.Now(),
		next:      make(chan bool),
		closed:    make(chan bool),
		taskClose: make(chan bool),
		timeClose: make(chan bool),
	}
	for i := 0; i < 3600; i++ {
		polling.slots[i] = make(map[string]*Task)
	}
	return
}

func (tp *TimePolling) Register(seconds time.Duration, key string, method TaskFunc, args []any) {
	t := time.Now().Add(time.Second * seconds)
	if tp.sTime.After(t) {
		log.Printf("Time error")
	}
	subSecond := t.Unix() - tp.sTime.Unix()
	cycleNum := int(subSecond / 3600)
	task := &Task{
		cycleNum: cycleNum,
		method:   method,
		params:   args,
	}
	tp.mu.Lock()
	curIndex := subSecond % 3600
	tasks := tp.slots[curIndex]
	if _, ok := tasks[key]; ok {
		log.Printf("task key exist")
	}
	tasks[key] = task
	tp.mu.Unlock()
}

func (tp *TimePolling) Run() {
	go tp.taskLoop()
	go tp.timeLoop()
	select {
	case <-tp.closed:
		tp.taskClose <- true
		tp.timeClose <- true
	}
}

func (tp *TimePolling) taskLoop() {
	defer func() {
		fmt.Println("taskLoop exit")
	}()
	for {
		select {
		case <-tp.taskClose:
			return
		case <-tp.next:
			tasks := tp.slots[tp.curIndex]
			if len(tasks) > 0 {
				for k, v := range tasks {
					if v.cycleNum == 0 {
						go v.method(v.params...)
						delete(tasks, k)
					} else {
						v.cycleNum--
					}
				}
			}
		}
	}
}

func (tp *TimePolling) timeLoop() {
	defer func() {
		fmt.Println("timeLoop exit")
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-tp.timeClose:
			return
		case <-ticker.C:
			log.Printf(time.Now().Format("2006-01-02 15:04:05"))
			if tp.curIndex >= 3599 {
				tp.curIndex = 0
			} else {
				tp.curIndex++
			}
			tp.next <- true
		}
	}
}

func (tp *TimePolling) Close() {
	tp.closed <- true
}
