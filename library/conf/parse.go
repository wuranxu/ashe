package conf

import (
	"encoding/json"
	"errors"
	sp "github.com/bitly/go-simplejson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	ConfigError = errors.New("config is not found")
)

type SqlConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  uint   `json:"timeout"`
	LogMode  bool   `json:"log_mode"`
}

func ParseConfig(filepath string, object interface{}, env string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	if env == "" {
		// 未配置环境变量, 默认读取dev配置
		env = "dev"
	}
	result, err := sp.NewJson(data)
	if err != nil {
		return err
	}
	if result.Get(env).Interface() == nil {
		return ConfigError
	}
	if data, err = result.Get(env).Encode(); err != nil {
		return err
	}
	err = json.Unmarshal(data, object)
	return err
}

func ParseYaml(filepath string, out interface{}) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}
