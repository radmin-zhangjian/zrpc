package utils

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type TimePolling struct {
	curIndex  int
	slots     []map[string]*Task
	sTime     time.Time
	mu        sync.Mutex
	next      chan bool
	closed    chan bool
	taskClose chan bool
	timeClose chan bool
	ticker    *time.Ticker
	cycle     int
}

type TaskFunc func(args ...any)

type Task struct {
	cycleNum int
	method   TaskFunc
	params   []any
}

func NewPolling(cycle int) (polling *TimePolling) {
	if cycle == 0 {
		cycle = 3600
	}
	polling = &TimePolling{
		curIndex:  0,
		sTime:     time.Now(),
		next:      make(chan bool, 10),
		closed:    make(chan bool),
		taskClose: make(chan bool),
		timeClose: make(chan bool),
		cycle:     cycle,
	}
	for i := 0; i < cycle; i++ {
		//polling.slots[i] = make(map[string]*Task)
		polling.slots = append(polling.slots, make(map[string]*Task))
	}
	return
}

func (tp *TimePolling) Register(seconds time.Duration, key string, method TaskFunc, args []any) (err error) {
	t := time.Now().Add(time.Second * seconds)
	if tp.sTime.After(t) {
		err = errors.New("time error")
		return
	}
	subSecond := t.Unix() - tp.sTime.Unix()
	cycle := int64(tp.cycle)
	cycleNum := int(subSecond / cycle)
	task := &Task{
		cycleNum: cycleNum,
		method:   method,
		params:   args,
	}
	tp.mu.Lock()
	curIndex := subSecond % cycle
	tasks := tp.slots[curIndex]
	if _, ok := tasks[key]; ok {
		err = errors.New("task key exist")
		return
	}
	tasks[key] = task
	tp.mu.Unlock()
	return nil
}

func (tp *TimePolling) Run() {
	go tp.taskLoop()
	go tp.timeLoop()
	select {
	case <-tp.closed:
		tp.timeClose <- true
		tp.taskClose <- true
	}
	close(tp.taskClose)
	close(tp.timeClose)
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
			tp.mu.Lock()
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
			tp.mu.Unlock()
		}
	}
}

func (tp *TimePolling) timeLoop() {
	defer func() {
		fmt.Println("timeLoop exit")
	}()

	tp.ticker = time.NewTicker(time.Second)
	defer tp.ticker.Stop()

	for {
		select {
		case <-tp.timeClose:
			close(tp.next)
			return
		case <-tp.ticker.C:
			log.Printf(time.Now().Format("2006-01-02 15:04:05"))
			if tp.curIndex >= tp.cycle-1 {
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
	close(tp.closed)
}
