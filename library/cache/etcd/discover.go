package etcd

import (
	"encoding/json"
	"fmt"
)

type Method struct {
	Auth bool // 是否需要登录
	Path string
}

func (m *Method) Marshal() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func RegisterMethod(client *Client, service, method string, auth bool) error {
	md := &Method{
		Auth: auth,
		Path: fmt.Sprintf("/%s/%s", service, method),
	}
	_, err := client.cli.Put(client.cli.Ctx(), fmt.Sprintf("%s.%s", service, method), md.Marshal())
	if err != nil {
		fmt.Println("注册方法失败, error: ", err)
		return err
	}
	return nil
}
