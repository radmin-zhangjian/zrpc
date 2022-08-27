package center

import (
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const (
	RoundRobin = 1
	Random     = 2
)

type SelectAlgorithm interface {
	Algorithm(mapVals []map[string]any) string
}

func SelectMode(mode int) SelectAlgorithm {
	switch mode {
	case RoundRobin:
		return &roundRobinAlgorithm{}
	case Random:
		return &randomAlgorithm{}
	default:
		return &roundRobinAlgorithm{}
	}
}

// SortAsc 正序排序
func SortAsc(source []map[string]any) {
	sort.Slice(source, func(i, j int) bool {
		a := source[i]["serverPort"].(string)
		b := source[j]["serverPort"].(string)
		return a < b
	})
}

// SortDesc 倒叙排序
func SortDesc(source []map[string]any) {
	sort.Slice(source, func(i, j int) bool {
		a := source[i]["serverPort"].(string)
		b := source[j]["serverPort"].(string)
		return a < b
	})
}

// 轮询选择
type roundRobinAlgorithm struct {
	mu sync.Mutex
}

var mode int = 0

// Algorithm 轮询选择
func (sm *roundRobinAlgorithm) Algorithm(mapVals []map[string]any) string {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	SortAsc(mapVals)
	mNum := len(mapVals)
	if mode >= mNum {
		mode = 0
	}
	// 轮询取出服务
	serviceVal := make(map[string]any)
	serviceVal = mapVals[mode]
	mode = (mode + 1) % mNum
	lastHeartBeat := serviceVal["lastHeartBeat"].(int64)
	tm := time.Unix(lastHeartBeat, 0)
	timeStr := tm.Format("2006-01-02 15:04:05")
	log.Printf("RedisFind - serverPort: %#v, lastHeartBeat: %#v, timeStr: %#v", serviceVal["serverPort"], serviceVal["lastHeartBeat"], timeStr)

	serverIp := serviceVal["serverIp"].(string)
	serverPort := serviceVal["serverPort"].(string)
	addr := serverIp + ":" + serverPort
	return addr
}

// 随机选择
type randomAlgorithm struct {
	mu sync.Mutex
}

// Algorithm 随机选择
func (sm *randomAlgorithm) Algorithm(mapVals []map[string]any) string {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	SortAsc(mapVals)
	// 随机取出服务
	serviceVal := make(map[string]any)
	rand.Seed(time.Now().UnixNano())
	kr := rand.Intn(len(mapVals))
	serviceVal = mapVals[kr]
	lastHeartBeat := serviceVal["lastHeartBeat"].(int64)
	tm := time.Unix(lastHeartBeat, 0)
	timeStr := tm.Format("2006-01-02 15:04:05")
	log.Printf("RedisFind - serverPort: %#v, lastHeartBeat: %#v, timeStr: %#v", serviceVal["serverPort"], serviceVal["lastHeartBeat"], timeStr)

	serverIp := serviceVal["serverIp"].(string)
	serverPort := serviceVal["serverPort"].(string)
	addr := serverIp + ":" + serverPort
	return addr
}
