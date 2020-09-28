package redis

import (
	"fmt"
	rds "github.com/go-redis/redis"
	"log"
	"sync"
	"time"
)

var (
	Pool *rds.Client
)

func NewClient(cfg RedisCliInfo) *rds.Client {
	var once sync.Once
	if Pool == nil {
		once.Do(func() {
			Pool = rds.NewClient(&rds.Options{
				Addr:        cfg.Server,
				DB:          cfg.Db,
				PoolSize:    cfg.PoolSize,
				Password:    cfg.Password,
				IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Second,
			})
			if _, err := Pool.Ping().Result(); err != nil {
				log.Fatal("redis连接失败, error: ", err)
			}
		})
	}
	return Pool
}

func NewClusterClient() error {
	client := rds.NewClusterClient(&rds.ClusterOptions{
		Addrs: []string{"106.13.173.14:7000", "106.13.173.14:7001", "106.13.173.14:7002",
			"106.13.173.14:7004", "106.13.173.14:7005", "106.13.173.14:7003",
		},
	})
	if err := client.Ping().Err(); err != nil {
		return err
	}
	fmt.Println(client.Set("sp", "nmsl", 0).Err())
	fmt.Println(client.Get("sp").Result())
	return nil
}
