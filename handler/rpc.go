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
	ArgsParseFailed = 10002 + iota
	LoginRequired
	MethodNotFound
	RemoteCallFailed
	IntervalServerError
)

var (
	ParamsError = errors.New("抱歉, 网络似乎开小差了")
	//ServiceMethodError = errors.New("抱歉, 网络似乎开小差了")
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

func (s *res) toJson() []byte {
	result, _ := json.Marshal(s)
	return result
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

func CallRpc(ctx iris.Context) {
	fmt.Println(ctx.Host())
	fmt.Println(ctx.Subdomain())
	fmt.Println(ctx.GetReferrer())
	fmt.Println(ctx.Host())
	result := new(res)
	params := make(Params)
	var userInfo *auth.CustomClaims
	if err := ctx.ReadJSON(&params); err != nil {
		response(ctx, result.Fill(ArgsParseFailed, SystemError))
		return
	}
	version := ctx.Params().Get("version")
	service := ctx.Params().Get("service")
	method := ctx.Params().Get("method")
	client, err := protocol.NewGrpcClient(version, service, method)
	defer client.Close()
	if err != nil {
		response(ctx, result.Fill(MethodNotFound, err))
		return
	}
	// 新增请求ip地址
	requestData, err := params.Marshal()
	if err != nil {
		result.Code = ArgsParseFailed
		result.Msg = err.Error()
		response(ctx, result)
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
