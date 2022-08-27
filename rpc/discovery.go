package rpc

import (
	"fmt"
	"strings"
	"zrpc/rpc/center"
)

// CreateServiceDiscovery 注册服务
func CreateServiceDiscovery(basePath string, registry string, password string, db int64, poolSize int64) (center.ServeDiscovery, error) {
	registryAddr := registry
	i := strings.Index(registryAddr, "://")
	regType := registryAddr[:i]
	regAddr := registryAddr[i+3:]

	switch regType {
	case "redis":
		return center.NewRedisDiscovery(basePath, regAddr, password, db, poolSize), nil
	default:
		return nil, fmt.Errorf("wrong registry type %s. only support etcd, consul and redis", regType)
	}
}

// ServiceDiscovery 发现服务
func ServiceDiscovery(basePath string, registry string, password string, db int64, poolSize int64) (center.ServeDiscovery, error) {
	registryAddr := registry
	i := strings.Index(registryAddr, "://")
	regType := registryAddr[:i]
	regAddr := registryAddr[i+3:]

	switch regType {
	case "redis":
		return center.NewRedisDiscovery(basePath, regAddr, password, db, poolSize), nil
	default:
		return nil, fmt.Errorf("wrong registry type %s. only support etcd, consul and redis", regType)
	}
}
