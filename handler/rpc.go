package handler

import (
	"ashe/library/auth"
	"ashe/protocol"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris"
)

const (
	// 参数解析失败 10001
	ArgsParseFailed = 10001 + iota
	// 用户未登录 10002
	LoginRequired
	// 方法未找到 10003
	MethodNotFound
	// rpc调用失败 10004
	RemoteCallFailed
	// 服务出错 10005
	IntervalServerError
)

var (
	InnerError  = errors.New("系统内部错误")
	SystemError = errors.New("抱歉, 网络似乎开小差了")
)

type Response interface {
	toJson() []byte
}

type res struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (s *res) Fill(code int32, msg interface{}, data ...interface{}) *res {
	if len(data) > 0 {
		return s.SetCode(code).SetMsg(msg).SetData(data[0])
	}
	return s.SetCode(code).SetMsg(msg)
}

func (s *res) SetCode(code int32) *res {
	s.Code = code
	return s
}

func (s *res) SetData(data interface{}) *res {
	s.Data = data
	return s
}

func (s *res) SetMsg(msg interface{}) *res {
	switch msg.(type) {
	case string:
		s.Msg = msg.(string)
	default:
		s.Msg = fmt.Sprintf("%v", msg)
	}
	return s
}

func (s *res) toApi(resp *protocol.Response) *res {
	if resp.ResultJson != nil {
		if err := json.Unmarshal(resp.ResultJson, &s.Data); err != nil {
			return s.Fill(IntervalServerError, InnerError)
		}
	}
	return s.Fill(resp.Code, resp.Msg)
}

type Params map[string]interface{}

func (p *Params) Marshal() (*protocol.Request, error) {
	return protocol.MarshalRequest(p)
}

func response(ctx iris.Context, r *res) {
	ctx.JSON(r)
}

// rpc调用接口
func CallRpc(ctx iris.Context) {
	result := new(res)
	params := make(Params)
	var userInfo *auth.CustomClaims
	if err := ctx.ReadJSON(&params); err != nil {
		response(ctx, result.Fill(ArgsParseFailed, SystemError))
		return
	}
	// 获取url中版本/APP/方法名(首字母小写, 与其他语言服务保持一致)
	version := ctx.Params().Get("version")
	service := ctx.Params().Get("service")
	method := ctx.Params().Get("method")
	client, err := protocol.NewGrpcClient(version, service, method)
	defer client.Close()
	if err != nil {
		response(ctx, result.Fill(MethodNotFound, err))
		return
	}
	requestData, err := params.Marshal()
	if err != nil {
		response(ctx, result.Fill(ArgsParseFailed, err))
		return
	}
	if !client.NoAuth() {
		// 需要解析token
		if userInfo, err = auth.Authrozation(ctx); err != nil {
			response(ctx, result.Fill(LoginRequired, err))
			return
		}
	}
	resp, err := client.Invoke(requestData, ctx.RemoteAddr(), userInfo)
	if err != nil {
		response(ctx, result.Fill(RemoteCallFailed, err))
		return
	}
	response(ctx, result.toApi(resp))
}
