package redis

import (
	"encoding/json"
	"fmt"
	redisCluster_ "github.com/chasex/redis-go-cluster"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

var RedisPool = RedisConnector{cluster: make(map[string]*redisCluster_.Cluster), cli: make(map[string]redis.Conn)}

type RedisConnector struct {
	cluster map[string]*redisCluster_.Cluster
	cli     map[string]redis.Conn
	lock    sync.RWMutex
}

type RedisCliInfo struct {
	Server      string `json:"server"` //链接地址
	Password    string `json:"password"`
	Port        int    `json:"port"`
	Db          int    `json:"db"`
	MaxIdle     int    `json:"max_idle"`
	MaxActive   int    `json:"max_active"`
	IdleTimeout int    `json:"idle_timeout"`
	PoolSize int `json:"pool_size"`
}

func (t RedisCliInfo) TableName() string { return "redis_cli_info" }

type RedisClusterInfo struct {
	StartNodes   string `json:"start_nodes"` //链接地址
	ConnTimeout  int    `json:"conn_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	KeepAlive    int    `json:"keep_alive"`
	AliveTime    int    `json:"alive_time"`
}

func (t RedisClusterInfo) TableName() string { return "redis_cluster_info" }

type RedisCon interface {
	Do(cmd string, args ...interface{}) (interface{}, error)
}

func (r RedisConnector) GetKey(key string) (redis.Conn, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	con, ok := r.cli[key]
	return con, ok
}

func (r RedisConnector) SetKey(key string, v redis.Conn) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.cli[key] = v
}

func (r RedisConnector) SetCKey(key string, v *redisCluster_.Cluster) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.cluster[key] = v
}

func (r RedisConnector) GetCKey(key string) (*redisCluster_.Cluster, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	c, ok := r.cluster[key]
	return c, ok
}

// 删除redis连接
func (r RedisConnector) CloseSingle(env, name string) {
	key := fmt.Sprintf("%s_%s", env, name)
	if pool, ok := r.cluster[key]; ok {
		pool.Close()
		delete(r.cluster, key)
		return
	}
	if pool, ok := r.cli[key]; ok {
		pool.Close()
		delete(r.cli, key)
	}
}

func (r RedisConnector) Close() {
	for _, v := range r.cluster {
		v.Close()
	}
	r.cluster = make(map[string]*redisCluster_.Cluster)
	for _, v := range r.cli {
		v.Close()
	}
	r.cli = make(map[string]redis.Conn)
}

func transNodes(params ...interface{}) []string {
	newList := make([]string, 0)
	for _, p := range params {
		newList = append(newList, p.(string))
	}
	return newList
}

func transNodesV2(nodes string) []string {
	return strings.Split(nodes, ",")
}

func transInt(v interface{}) int64 {
	data, err := v.(json.Number).Int64()
	if err != nil {
		return 5
	}
	return data
}

func NewRedisCluster(conn map[string]interface{}) (*redisCluster_.Cluster, error) {
	cliCluster, err := redisCluster_.NewCluster(
		&redisCluster_.Options{
			StartNodes:   transNodes(conn["StartNodes"].([]interface{})...),
			ConnTimeout:  time.Duration(transInt(conn["ConnTimeout"])) * time.Second,
			ReadTimeout:  time.Duration(transInt(conn["ReadTimeout"])) * time.Second,
			WriteTimeout: time.Duration(transInt(conn["WriteTimeout"])) * time.Second,
			KeepAlive:    int(transInt(conn["KeepAlive"])),
			AliveTime:    time.Duration(transInt(conn["AliveTime"])) * time.Second,
		})
	if err != nil {
		log.Errorf("redis cluster connection error: %v", err)
	}
	return cliCluster, nil
}

func NewCliPool(conn RedisCliInfo) *redis.Pool {
	p := &redis.Pool{
		MaxIdle:     conn.MaxIdle,
		MaxActive:   conn.MaxActive,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conn.Server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", conn.Password); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", int(conn.Db)); err != nil {
				c.Close()
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	c := p.Get()
	if c.Err() != nil {
		log.Errorf("redis链接失败, error: %v", c.Err())
		return nil
	}
	defer c.Close()
	return p
}
