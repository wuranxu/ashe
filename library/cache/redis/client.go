package redis

import (
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
