package protocol

import (
	"ashe/library/auth"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var (
	UserInfoNotFound   = errors.New("未获取到用户信息")
	UserInfoParseError = errors.New("用户信息解析失败")
	DecodeError        = "解析返回数据失败"
)

func Unmarshal(in *Request, data interface{}) error {
	if err := json.Unmarshal(in.RequestJson.GetValue(), data); err != nil {
		return err
	}
	return nil
}

func MarshalRequest(out *Request, data interface{}) error {
	var result any.Any
	var msg proto.Message
	bt, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(bt, msg)
	if err != nil {
		return err
	}
	err = result.MarshalFrom(msg)
	if err != nil {
		return err
	}
	out.RequestJson = &result
	return nil
}

func Marshal(out *Response, data interface{}) {
	var result any.Any
	var msg proto.Message
	bt, err := json.Marshal(data)
	if err != nil {
		out.ResultJson = nil
		return
	}
	err = proto.Unmarshal(bt, msg)
	if err != nil {
		out.ResultJson = nil
		out.Msg = DecodeError
		return
	}
	err = result.MarshalFrom(msg)
	if err != nil {
		out.ResultJson = nil
		out.Msg = DecodeError
		return
	}
	out.ResultJson = &result
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
