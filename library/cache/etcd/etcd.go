package etcd

import (
	"ashe/common"
	"ashe/library/logging"
	"context"
	v3 "github.com/coreos/etcd/clientv3"
	"sync"
	"time"
)

type Client struct {
	kv     v3.KV
	cli    *v3.Client
	scheme string
}

var (
	once       sync.Once
	EtcdClient *Client
	log        = logging.NewLog("etcd")
)

func (cl *Client) Kv() v3.KV {
	return cl.kv
}

func (cl *Client) Cli() *v3.Client {
	return cl.cli
}

func (cl *Client) Set(key, value string) bool {
	res, err := cl.kv.Put(context.TODO(), key, value, v3.WithPrevKV())
	if err != nil {
		return false
	}
	if res.PrevKv != nil {
		log.Infoln("设置key:", string(res.PrevKv.Key), "成功, 上一个值为", string(res.PrevKv.Value))
	}
	return true
}

func (cl *Client) GetSingle(key string) string {
	res, err := cl.kv.Get(context.TODO(), key)
	if err != nil {
		return ""
	}
	if len(res.Kvs) == 0 {
		return ""
	}
	return string(res.Kvs[0].Value)
}

func (cl *Client) GetPattern(key string) (result map[string]string) {
	res, err := cl.kv.Get(context.TODO(), key, v3.WithPrefix())
	result = make(map[string]string)
	if err != nil {
		return
	}
	for _, item := range res.Kvs {
		result[string(item.Key)] = string(item.Value)
	}
	return
}

// 关闭客户端
func (cl *Client) Close() error {
	return cl.cli.Close()
}

func NewClient(cfg common.EtcdConfig) (*Client, error) {
	var err error
	var cli *v3.Client
	var kv v3.KV
	if EtcdClient == nil {
		once.Do(func() {
			cli, err = v3.New(v3.Config{Endpoints: cfg.Endpoints, DialTimeout: time.Second * cfg.DialTimeout})
			if err != nil {
				cli = nil
				log.Panic("连接etcd失败, error: ", err)
			}
			kv = v3.NewKV(cli)

		})
		EtcdClient = &Client{kv: kv, cli: cli, scheme: common.Conf.Scheme}
	}
	return EtcdClient, nil
}
