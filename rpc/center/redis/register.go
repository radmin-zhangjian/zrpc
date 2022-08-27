package redis

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
	"sync"
	"time"
)

type ServiceDiscovery struct {
	regClient *redis.Client
	basePath  string
	mu        sync.Mutex
}

func NewRedis(basePath string, host string, port string, password string, db int64, poolSize int64) *ServiceDiscovery {
	if GetRedis() == nil {
		InitRedis(host, port, password, db, poolSize)
	}
	path := strings.Trim(basePath, "/")
	return &ServiceDiscovery{basePath: path, regClient: rdb}
}

// ServeRegister redis服务注册
func (sd *ServiceDiscovery) ServeRegister(addr string) {
	// hash 存储
	//sd.serveHSet(addr)
	// set 存储
	sd.serveSet(addr)
}

// ServeDiscovery redis服务发现
func (sd *ServiceDiscovery) ServeDiscovery() []map[string]any {
	// hash get 方式
	//return sd.serveHGetAll()
	// scan get 方式
	return sd.serveScan()
}

// TickerHeartbeatTask 心跳监测
func (sd *ServiceDiscovery) TickerHeartbeatTask(addr string) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		sd.ServeRegister(addr)
	}
}

func (sd *ServiceDiscovery) Close(addr string) {
	i := strings.Index(addr, ":")
	serverIp := addr[:i]
	serverPort := addr[i+1:]
	key := sd.basePath
	nodeId := []byte(addr + serverIp + serverPort)
	newSig := md5.Sum(nodeId) //转成加密编码
	newArr := fmt.Sprintf("%x", newSig)
	nodeKey := strings.ToTitle(newArr)
	err := HDel(key, nodeKey)
	if err != nil {
		log.Printf("redis discovery close: %v", err)
	}
}

// serveHSet redis服务注册
func (sd *ServiceDiscovery) serveHSet(addr string) {
	i := strings.Index(addr, ":")
	serverIp := addr[:i]
	serverPort := addr[i+1:]
	key := sd.basePath
	nodeId := []byte(addr + serverIp + serverPort)
	newSig := md5.Sum(nodeId) //转成加密编码
	newArr := fmt.Sprintf("%x", newSig)
	nodeKey := strings.ToTitle(newArr)
	sd.mu.Lock()
	mapNode := map[string]any{
		"nodeId":         nodeKey,
		"serverIp":       serverIp,
		"serverPort":     serverPort,
		"serviceVersion": "1.0",
		"serviceName":    sd.basePath,
		"lastHeartBeat":  time.Now().UnixNano()/1e9 + 30,
	}
	value, _ := json.Marshal(mapNode)
	mapValue := map[string]any{
		nodeKey: value,
	}
	sd.mu.Unlock()
	err := HSet(key, mapValue)
	if err != nil {
		log.Printf("redis discovery register: %v", err)
	}
}

// serveHGetAll redis服务发现
func (sd *ServiceDiscovery) serveHGetAll() []map[string]any {
	key := sd.basePath
	val, err := HGetAll(key)
	if err != nil {
		log.Printf("redis discovery get: %v", err)
	}

	sd.mu.Lock()
	defer sd.mu.Unlock()

	var i int
	disArrs := []map[string]any{}
	for _, v := range val.(map[string]string) {
		mapVal := make(map[string]any)
		b := []byte(v)
		json.Unmarshal(b, &mapVal)
		timeNow := time.Now().UnixNano() / 1e9
		lastHeartBeat := int64(mapVal["lastHeartBeat"].(float64))

		//tm := time.Unix(lastHeartBeat, 0)
		//timeStr := tm.Format("2006-01-02 15:04:05")
		//log.Printf("RedisFind - serverPort: %#v, lastHeartBeat: %#v, timeStr: %#v", mapVal["serverPort"], lastHeartBeat, timeStr)
		if lastHeartBeat < timeNow {
			continue
		}
		mapVal["lastHeartBeat"] = lastHeartBeat
		disArrs = append(disArrs, mapVal)
		i++
	}
	return disArrs
}

// serveSet redis服务注册
func (sd *ServiceDiscovery) serveSet(addr string) {
	key := sd.basePath + "/" + addr
	i := strings.Index(addr, ":")
	serverIp := addr[:i]
	serverPort := addr[i+1:]
	nodeId := []byte(addr + serverIp + serverPort)
	newSig := md5.Sum(nodeId) //转成加密编码
	newArr := fmt.Sprintf("%x", newSig)
	nodeKey := strings.ToTitle(newArr)
	sd.mu.Lock()
	mapNode := map[string]any{
		"nodeId":         nodeKey,
		"serverIp":       serverIp,
		"serverPort":     serverPort,
		"serviceVersion": "1.0",
		"serviceName":    sd.basePath,
		"lastHeartBeat":  time.Now().UnixNano()/1e9 + 30,
	}
	value, _ := json.Marshal(mapNode)
	sd.mu.Unlock()
	err := SetStringExp(key, value, 30*time.Second)
	if err != nil {
		log.Printf("redis discovery register: %v", err)
	}
}

// serveScan redis服务发现
func (sd *ServiceDiscovery) serveScan() []map[string]any {
	key := sd.basePath
	vals, err := GetScan(key)
	if err != nil {
		log.Printf("redis discovery get: %v", err)
	}

	sd.mu.Lock()
	defer sd.mu.Unlock()

	var i int
	disArrs := []map[string]any{}
	for _, v := range vals.([]string) {
		mapVal := make(map[string]any)
		b := []byte(v)
		json.Unmarshal(b, &mapVal)
		timeNow := time.Now().UnixNano() / 1e9
		lastHeartBeat := int64(mapVal["lastHeartBeat"].(float64))

		//tm := time.Unix(lastHeartBeat, 0)
		//timeStr := tm.Format("2006-01-02 15:04:05")
		//log.Printf("RedisFind - serverPort: %#v, lastHeartBeat: %#v, timeStr: %#v", mapVal["serverPort"], lastHeartBeat, timeStr)
		if lastHeartBeat < timeNow {
			continue
		}
		mapVal["lastHeartBeat"] = lastHeartBeat
		disArrs = append(disArrs, mapVal)
		i++
	}
	return disArrs
}
