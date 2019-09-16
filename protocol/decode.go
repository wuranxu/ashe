package protocol

import (
	"ashe/library/auth"
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/grpc/metadata"
)

var (
	UserInfoNotFound   = errors.New("未获取到用户信息")
	UserInfoParseError = errors.New("用户信息解析失败")
)

func Unmarshal(in *Request, data interface{}) error {
	if err := json.Unmarshal([]byte(in.RequestJson), data); err != nil {
		return err
	}
	return nil
}

func Marshal(out *Response, data interface{}) {
	bt, err := json.Marshal(data)
	if err != nil {
		out.ResultJson = ""
		return
	}
	out.ResultJson = string(bt)
}

func GetHeader(ctx context.Context) map[string][]string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md
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