package common

import (
	"ashe/library/cache/redis"
	"ashe/library/conf"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	Conf = new(Config)
	Env = "dev"
)

type EtcdConfig struct {
	Endpoints   []string      `json:"endpoints"`
	DialTimeout time.Duration `json:"dial-timeout"`
}

type Config struct {
	Etcd     EtcdConfig         `json:"etcd"`
	Database conf.SqlConfig     `json:"database"`
	Redis    redis.RedisCliInfo `json:"redis"`
	Scheme   string             `json:"scheme"`
}

func Init(file string) {
	fmt.Println("本机环境: ", Env)
	var (
		once sync.Once
		err  error
	)
	once.Do(func() {
		err = conf.ParseConfig(file, Conf, Env)
		if err != nil {
			log.Fatalf("获取配置出错, error: %v", err)
		}
	})
}
