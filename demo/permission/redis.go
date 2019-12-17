package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var Client *redis.Client

type RedisClient struct {
}

func initRedis(db []int) {
	dbCode := 0
	if len(db) > 0 {
		dbCode = db[0]
	}

	Client = redis.NewClient(&redis.Options{
		Addr:         "192.168.10.33:6379",
		PoolSize:     1000,
		ReadTimeout:  time.Millisecond * time.Duration(100),
		WriteTimeout: time.Millisecond * time.Duration(100),
		IdleTimeout:  time.Second * time.Duration(60),
		DB:           dbCode,
	})

	_, err := Client.Ping().Result()
	if err != nil {
		panic("init redis error")
	} else {
		fmt.Println("init redis ok")
	}
}

func (this *RedisClient) Get(key string, db ...int) (string, bool) {
	initRedis(db)
	r, err := Client.Get(key).Result()
	if err != nil {
		return "", false
	}
	return r, true
}

func (this *RedisClient) SetExpTime(key string, val interface{}, expTime int32, db ...int) {
	initRedis(db)
	Client.Set(key, val, time.Duration(expTime)*time.Second)
}

func (this *RedisClient) Set(key string, val interface{}) {

}
