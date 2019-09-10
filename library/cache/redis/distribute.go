// 基于redis的分布式锁

package redis

import (
	"errors"
	"fmt"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	RedisConfigError = errors.New("获取redis配置失败, 请检查配置")
	RedisConnError   = errors.New("连接redis失败")
)

func DistributeLock(cfg RedisCliInfo, taskId int) (*redsync.Mutex, error) {
	pool, err := GetPool(cfg)
	if err != nil || pool == nil {
		log.Errorln("获取redis配置失败, error: %v", err)
		return nil, RedisConnError
	}
	poolList := []redsync.Pool{redsync.Pool(pool)}
	lock := redsync.New(poolList)
	name := fmt.Sprintf("distribute_%d", taskId)
	m := lock.NewMutex(name, redsync.SetExpiry(60*time.Second), redsync.SetTries(10))
	err = m.Lock()
	if err != nil {
		log.Infof("有其他机器正在执行该任务, 锁被占用了")
		return nil, err
	}
	return m, nil
}

func GetPool(cfg RedisCliInfo) (pool *redis.Pool, err error) {
	defer func() {
		errInfo := recover()
		if errInfo != nil {
			pool = nil
			err = fmt.Errorf("%v", errInfo)
		}
	}()
	pool = NewCliPool(cfg)
	return
}
