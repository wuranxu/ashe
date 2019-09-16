package etcd

import (
	"ashe/common"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"reflect"
	"time"
)

func (cl *Client) RegisterService(name, addr string, ttl int64) error {
	ticker := time.NewTicker(time.Second * time.Duration(ttl))

	go func() {
		for {
			getResp, err := cl.cli.Get(context.Background(), "/"+cl.scheme+"/"+name+"/"+addr)
			if err != nil {
				fmt.Println(err)
			} else if getResp.Count == 0 {
				if err = cl.withAlive(name, addr, ttl); err != nil {
					fmt.Println(err)
				}
			} else {
			}
			<-ticker.C
		}
	}()
	return nil
}

func (cl *Client) withAlive(name, addr string, ttl int64) error {
	leaseResp, err := cl.cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	fmt.Printf("key:%v\n", "/"+cl.scheme+"/"+name+"/"+addr)
	if _, err := cl.cli.Put(context.Background(), "/"+cl.scheme+"/"+name+"/"+addr, addr, clientv3.WithLease(leaseResp.ID)); err != nil {
		return err
	}

	if _, err := cl.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil {
		return err
	}
	return nil
}

func (cl *Client) UnRegister(name, addr string) error {
	if cl.cli != nil {
		_, err := cl.cli.Delete(context.Background(), "/"+cl.scheme+"/"+name+"/"+addr)
		return err
	}
	return nil
}

// 传入yaml配置
func (cl *Client) RegisterApi(name string, data interface{}, config common.YamlConfig) error {
	inf := reflect.ValueOf(data)
	for i := 0; i < inf.NumMethod(); i++ {
		methodName := inf.Type().Method(i).Name
		md, ok := config.Method[methodName]
		if !ok {
			// 说明配置文件没有包含此方法
			log.Fatal("注册Api失败, service.yaml文件未包含此方法: ", methodName)
		}
		err := RegisterMethod(cl, name, methodName, md.Auth)
		if err != nil {
			return err
		}
	}
	return nil
}
