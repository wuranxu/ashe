package handler

import (
	"ashe/protocol"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"github.com/wuranxu/library/auth"
	"io/ioutil"
	"strings"
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

func (s *res) Build(code int32, msg interface{}, data ...interface{}) *res {
	s.SetMsg(msg).Code = code
	if len(data) > 0 {
		s.Data = data[0]
	}
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
			return s.Build(IntervalServerError, InnerError)
		}
	}
	return s.Build(resp.Code, resp.Msg)
}

type Params map[string]interface{}

func (p Params) Marshal() (*protocol.Request, error) {
	return protocol.MarshalRequest(p)
}

func (p Params) MakeFile(ctx iris.Context) {
	list := fileNameList(ctx)
	if list == nil {
		return
	}
	fileList := make([]map[string]interface{}, 0, len(list))
	for _, l := range list {
		file, header, err := ctx.FormFile(l)
		if err != nil {
			continue
		}
		buf, err := ioutil.ReadAll(file)
		if err != nil {
			continue
		}
		file.Close()
		fileList = append(fileList, map[string]interface{}{
			"filename": header.Filename,
			"size":     header.Size,
			"content":  buf,
		})
	}
	p["fileList"] = fileList
}

func response(ctx iris.Context, r *res) {
	ctx.JSON(r)
}

func fileNameList(ctx iris.Context) []string {
	fileList := ctx.URLParam("files")
	if fileList == "" {
		return nil
	}
	return strings.Split(fileList, ";")
}

// rpc调用接口
func CallRpc(ctx iris.Context) {
	result := new(res)
	params := make(Params)
	var userInfo *auth.CustomClaims
	// 如果是form
	if strings.Contains(ctx.GetHeader("Content-Type"), "form") {
		values := ctx.FormValues()
		params.MakeFile(ctx)
		for k, v := range values {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	} else {
		if err := ctx.ReadJSON(&params); err != nil {
			response(ctx, result.Build(ArgsParseFailed, SystemError))
			return
		}
	}
	// 获取url中版本/APP/方法名(首字母小写, 与其他语言服务保持一致)
	version := ctx.Params().Get("version")
	service := ctx.Params().Get("service")
	method := ctx.Params().Get("method")
	client, err := protocol.NewGrpcClient(version, service, method)
	defer client.Close()
	if err != nil {
		response(ctx, result.Build(MethodNotFound, err))
		return
	}
	requestData, err := params.Marshal()
	if err != nil {
		response(ctx, result.Build(ArgsParseFailed, err))
		return
	}
	if !client.NoAuth() {
		// 需要解析token
		if userInfo, err = Authorization(ctx); err != nil {
			response(ctx, result.Build(LoginRequired, err))
			return
		}
	}
	resp, err := client.Invoke(requestData, ctx.RemoteAddr(), userInfo)
	if err != nil {
		response(ctx, result.Build(RemoteCallFailed, err))
		return
	}
	response(ctx, result.toApi(resp))
}
