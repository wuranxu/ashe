package etcd

import (
	"fmt"
)

func RegisterMethod(client *Client, service, method string) error {
	_, err := client.cli.Put(client.cli.Ctx(), fmt.Sprintf("%s.%s", service, method), fmt.Sprintf("/%s/%s", service, method))
	if err != nil {
		fmt.Println("注册方法失败, error: ", err)
		return err
	}
	client.cli.Get(client.cli.Ctx(), fmt.Sprintf("%s.%s", service, method))
	return nil
}
