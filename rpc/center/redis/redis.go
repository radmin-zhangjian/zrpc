package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

func GetRedis() *redis.Client {
	return rdb
}

func InitRedis(host string, port string, password string, db int64, poolSize int64) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       int(db),
		PoolSize: int(poolSize), // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("redis connect get failed.%v", err)
		return
	}
}

// SetStringExp 存值
func SetStringExp(key string, value any, exp time.Duration) error {
	if key == "" || value == nil {
		return errors.New("无效的参数")
	}
	err := rdb.Set(ctx, key, value, exp).Err()
	if err != nil {
		return errors.New(err.Error())
	}
	//log.Printf("push key: %v, value: %v", key, value)
	return nil
}

// SetString 存值
func SetString(key string, value any) error {
	return SetStringExp(key, value, 0)
}

// GetString 取值
func GetString(key string) (any, error) {
	if key == "" {
		return "", errors.New("无效的参数")
	}
	res, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("数据为空或nil")
	}
	if err != nil {
		return "", errors.New(err.Error())
	}
	return res, nil
}

// GetScan 取值
func GetScan(key string) (any, error) {
	if key == "" {
		return "", errors.New("无效的参数")
	}

	var cursor uint64
	keys, cursor, err := rdb.Scan(ctx, cursor, key+"/*", 100).Result()
	if err != nil {
		fmt.Println("scan keys failed err:", err)
		return "", errors.New(err.Error())
	}
	//fmt.Println("scan keys:", keys)

	var vals []string
	for _, key := range keys {
		//fmt.Println("key:", key)
		sType, err := rdb.Type(ctx, key).Result()
		if err != nil {
			fmt.Println("get type failed :", err)
			continue
		}
		//fmt.Printf("key :%v ,type is %v\n", key, sType)
		var val string
		if sType == "string" {
			val, err = rdb.Get(ctx, key).Result()
			if err != nil {
				fmt.Println("get key values failed err:", err)
				continue
			}
			//fmt.Printf("key :%v ,value :%v\n", key, val)
		} else if sType == "list" {
			val, err = rdb.LPop(ctx, key).Result()
			if err != nil {
				fmt.Println("get list value failed :", err)
				continue
			}
			//fmt.Printf("key:%v value:%v\n", key, val)
		}
		if val != "" {
			vals = append(vals, val)
		}
	}

	return vals, nil
}

// HSet 哈希存值
func HSet(key string, value any) error {
	if key == "" || value == nil {
		return errors.New("无效的参数")
	}
	err := rdb.HSet(ctx, key, value).Err()
	if err != nil {
		return errors.New(err.Error())
	}
	//log.Printf("push key: %v, value: %v", key, value)
	return nil
}

// HGetAll 哈希取值
func HGetAll(key string) (any, error) {
	if key == "" {
		return "", errors.New("无效的参数")
	}
	res, err := rdb.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("数据为空或nil")
	}
	if err != nil {
		return "", errors.New(err.Error())
	}
	return res, nil
}

// HDel 删除字段
func HDel(key string, field string) error {
	if key == "" || field == "" {
		return errors.New("无效的参数")
	}
	err := rdb.HDel(ctx, key, field).Err()
	if err != nil {
		return errors.New(err.Error())
	}
	//log.Printf("push key: %v, value: %v", key, value)
	return nil
}
