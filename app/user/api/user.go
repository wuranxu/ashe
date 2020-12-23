package api

import (
	"ashe/app/user/models"
	"ashe/app/user/utils"
	"ashe/exception"
	"ashe/library/check"
	"ashe/protocol"
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
)

const (
	Success = 0

	ParamsError = iota + 10010
	RegisterError
	LoginParamsError
	LoginFailed
	TokenParseFailed
	EditUserFailed
)

var (
	OperateSuccess = "操作成功"
	LoginSuccess   = "登录成功"

	ParamsInValid     = errors.New("非法json数据")
	ParamsCheckError  = errors.New("请求参数不合法, 请检查")
	ParamsCheckFailed = exception.ErrString(ParamsCheckError.Error())

	//log = logging.NewLog("userService")
)

type UserApi struct {
}

type LoginForm struct {
	Username string `json:"username" validate:"gt=0"`
	Password string `json:"password" validate:"gt=0"`
}

type EditForm struct {
	Nickname string `json:"nickname" validate:"gt=0"`
	Email    string `json:"email" validate:"gt=0"`
}

func (*UserApi) Register(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	usr := new(models.AsheUser)
	res := new(protocol.Response)
	if err := protocol.Unmarshal(in, usr); err != nil {
		return res.Build(ParamsError, ParamsInValid), nil
	}
	// 校验参数
	if err := check.Check(usr, ParamsCheckFailed); err != nil {
		return res.Build(ParamsError, err), nil
	}
	if err := usr.Register(); err != nil {
		return res.Build(RegisterError, err), nil
	}
	res.Msg = "注册成功"
	return res, nil
}

func (*UserApi) Login(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	metadata.FromIncomingContext(ctx)
	var form LoginForm
	if err := protocol.Unmarshal(in, &form); err != nil {
		return res.Build(LoginParamsError, ParamsInValid), nil
	}
	// 校验参数
	if err := check.Check(&form, ParamsCheckFailed); err != nil {
		return res.Build(ParamsError, err), nil
	}
	pwd := utils.Encode(form.Password)
	user, token, err := models.LoginVerify(form.Username, pwd)
	if err != nil {
		return res.Build(LoginFailed, err), nil
	}
	return res.Build(Success, LoginSuccess, map[string]interface{}{
		"token": token,
		"user":  user.AsheUserJson(),
	}), nil
}

func (*UserApi) Edit(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	user, err := protocol.FetchUserInfo(ctx)
	if err != nil {
		return res.Build(TokenParseFailed, err), nil
	}
	data := new(EditForm)
	if err = protocol.Unmarshal(in, data); err != nil {
		return res.Build(ParamsError, ParamsInValid), nil
	}
	// 校验参数
	if err := check.Check(data, ParamsCheckFailed); err != nil {
		return res.Build(ParamsError, err), nil
	}
	if err = models.Edit(data.Nickname, data.Email, user.ID); err != nil {
		res.Code = EditUserFailed
		res.Msg = err.Error()
		return res, nil
	}
	res.Msg = "修改成功"
	return res, nil
}

func (*UserApi) InsertUserLog(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	data := new(models.TUserLog)
	if err := protocol.Unmarshal(in, data); err != nil {
		return res.Build(ParamsError, ParamsCheckError), nil
	}
	if err := models.Insert(data); err != nil {
		res.Msg = "失败"
		res.Code = ParamsError
	} else {
		protocol.Marshal(res, data)
		res.Msg = "修改成功"
	}
	return res, nil
}
