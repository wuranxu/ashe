package handler

import (
	"ashe/library/auth"
	"ashe/protocol"
	"encoding/json"
	"errors"
	"github.com/kataras/iris"
)

const (
	ArgsParseFailed = 102 + iota
	LoginRequired 
)

var (
	ParamsError        = errors.New("抱歉, 网络似乎开小差了")
	//ServiceMethodError = errors.New("抱歉, 网络似乎开小差了")
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

func (s *res) toApi(resp protocol.Response) *res {
	s.Code, s.Msg, s.Data = resp.Code, resp.Msg, resp.ResultJson
	var data interface{}
	if err := json.Unmarshal([]byte(resp.ResultJson), &data); err == nil {
		s.Data = data
	}
	return s
}

type Params map[string]interface{}

func (p *Params) Marshal() (*protocol.Request, error) {
	bt, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return &protocol.Request{RequestJson: string(bt)}, nil
}

func response(ctx iris.Context, r *res) {
	ctx.JSON(r)
}

func parse(params map[string]interface{}) (service string, err error) {
	action, ok := params["action"]
	var (
		act string
		suc bool
	)
	if act, suc = action.(string); !ok || !suc {
		err = ParamsError
		return
	}
	service = act
	return service, nil
}

func CallRpc(ctx iris.Context) {
	result := new(res)
	params := make(Params)
	var userInfo *auth.CustomClaims
	if err := ctx.ReadJSON(&params); err != nil {
		result.Code = 40000
		result.Msg = "抱歉，网络似乎开小差了"
		response(ctx, result)
		return
	}
	service := ctx.Params().Get("service")
	method := ctx.Params().Get("method")
	client, err := protocol.NewGrpcClient(service, method)
	defer client.Close()
	if err != nil {
		result.Msg = err.Error()
		response(ctx, result)
		return
	}
	// 新增请求ip地址
	params["remote_ip"] = ctx.RemoteAddr()
	requestData, err := params.Marshal()
	if err != nil {
		result.Code = ArgsParseFailed
		result.Msg = err.Error()
		response(ctx, result)
		return
	}
	if client.Auth() {
		// 需要解析token
		if userInfo, err = auth.Authrozation(ctx); err != nil {
			result.Code = LoginRequired
			result.Msg = err.Error()
			response(ctx, result)
			return
		}
	}
	resp, err := client.Invoke(requestData, userInfo)
	if err != nil {
		result.Msg = err.Error()
		response(ctx, result)
		return
	}
	response(ctx, result.toApi(*resp))
}
