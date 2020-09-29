package etcd

import (
	"encoding/json"
	"fmt"
	"unicode"
)

type Method struct {
	NoAuth bool // 是否需要登录
	Path   string
}

func (m *Method) Marshal() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func capitalize(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func lowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func RegisterMethod(client *Client, version, service, method string, auth bool) error {
	md := &Method{
		NoAuth: auth,
		Path:   fmt.Sprintf("/%s/%s", service, method),
	}
	fullPath := fmt.Sprintf("%s.%s.%s", version, lowerFirst(service), lowerFirst(method))
	_, err := client.cli.Put(client.cli.Ctx(), fullPath, md.Marshal())
	if err != nil {
		log.Errorf("注册方法失败, error: ", err)
		return err
	}
	return nil
}
