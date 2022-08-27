package center

import (
	"strings"
	"zrpc/rpc/center/redis"
)

type ServeDiscovery interface {
	ServeRegister(addr string)
	ServeDiscovery() []map[string]any
	TickerHeartbeatTask(addr string)
	Close(addr string)
}

// NewRedisDiscovery redis服务发现
func NewRedisDiscovery(basePath string, regAddr string, password string, db int64, poolSize int64) ServeDiscovery {
	i := strings.Index(regAddr, ":")
	host := regAddr[:i]
	port := regAddr[i+1:]
	return redis.NewRedis(basePath, host, port, password, db, poolSize)
}

// NewConsulDiscovery consul服务发现
func NewConsulDiscovery(basePath string, regAddr string) ServeDiscovery {
	//i := strings.Index(regAddr, ":")
	//host := regAddr[:i]
	//port := regAddr[i+1:]

	return nil
}
