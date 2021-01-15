package common

import (
	rds "github.com/go-redis/redis"
	"github.com/wuranxu/library/conf"
	logger "github.com/wuranxu/library/log"
	"sync"
	"time"
)

var (
	Conf = new(Config)
	Env  = "DEV"
	log  = logger.InitLogger("logs/config.log")
	once sync.Once
)

type EtcdConfig struct {
	Endpoints   []string      `json:"endpoints"`
	DialTimeout time.Duration `json:"dial-timeout"`
}

type Config struct {
	Etcd     EtcdConfig     `json:"etcd"`
	Database conf.SqlConfig `json:"database"`
	Redis    rds.Options    `json:"redis"`
	Scheme   string         `json:"scheme"`
}

type YamlConfig struct {
	Service string        `yaml:"service"`
	Version string        `yaml:"version"`
	Port    int           `yaml:"port"`
	Method  map[string]Md `yaml:"method"`
}

type Md struct {
	NoAuth bool `yaml:"no_auth"`
}

func Init(file string) {
	log.Info("本机环境: ", Env)
	var err error
	once.Do(func() {
		err = conf.ParseConfig(file, Conf, Env)
		if err != nil {
			log.Fatalf("获取配置出错, error: %v", err)
		}
	})
}
