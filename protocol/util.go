package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wuranxu/library/auth"
	"google.golang.org/grpc/metadata"
)

var (
	UserInfoNotFound   = errors.New("未获取到用户信息")
	UserInfoParseError = errors.New("用户信息解析失败")
	DecodeError        = "解析返回数据失败"
)

func (m *Response) Build(code int32, msg interface{}, data ...interface{}) *Response {
	m.Code = code
	if msg != nil {
		switch msg.(type) {
		case string:
			m.Msg = msg.(string)
		case error:
			m.Msg = msg.(error).Error()
		default:
			m.Msg = fmt.Sprintf("%v", msg)
		}
	}
	if len(data) > 0 {
		bt, err := json.Marshal(data)
		if err != nil {
			m.Msg = err.Error()
		}
		m.ResultJson = bt
	}
	return m
}

func Unmarshal(in *Request, data interface{}) error {
	if err := json.Unmarshal(in.RequestJson, data); err != nil {
		return err
	}
	return nil
}

func MarshalRequest(data interface{}) (*Request, error) {
	in := new(Request)
	result, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	in.RequestJson = result
	return in, nil
}

func Marshal(out *Response, data interface{}) {
	bt, err := json.Marshal(data)
	if err != nil {
		out.ResultJson = nil
		out.Msg = DecodeError
		return
	}
	out.ResultJson = bt
}

func GetHeader(ctx context.Context) map[string][]string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md
}

func GetHeaderKey(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if get := md.Get(key); len(get) > 0 {
			return get[0]
		}
	}
	return ""
}

func FetchUserInfo(ctx context.Context) (*auth.CustomClaims, error) {
	headers := GetHeader(ctx)
	if result, ok := headers["user"]; ok {
		if len(result) > 0 {
			var claims auth.CustomClaims
			if err := json.Unmarshal([]byte(result[0]), &claims); err != nil {
				return nil, UserInfoParseError
			}
			return &claims, nil
		}
	}
	return nil, UserInfoNotFound
}
