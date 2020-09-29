package api

import (
	"ashe/app/user/models"
	"ashe/app/user/utils"
	"ashe/exception"
	"ashe/library/check"
	"ashe/library/logging"
	"ashe/protocol"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

const (
	ParamsError = iota + 10010
	RegisterError
	LoginParamsError
	LoginFailed
	TokenParseFailed
	EditUserFailed
)

var (
	ParamsInValid    = "非法json数据"
	ParamsCheckError = exception.ErrString(ParamsInValid)

	log = logging.NewLog("userService")
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
		res.Msg = ParamsInValid
		res.Code = ParamsError
		return res, nil
	}
	// 校验参数
	if err := check.Check(usr, ParamsCheckError); err != nil {
		res.Code = ParamsError
		res.Msg = err.Error()
		return res, nil
	}
	if err := usr.Register(); err != nil {
		res.Msg = err.Error()
		res.Code = RegisterError
		return res, nil
	}
	res.Msg = "注册成功"
	return res, nil
}

func (*UserApi) Login(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	metadata.FromIncomingContext(ctx)
	var form LoginForm
	if err := protocol.Unmarshal(in, &form); err != nil {
		res.Msg = ParamsInValid
		res.Code = LoginParamsError
		return res, nil
	}
	// 校验参数
	if err := check.Check(&form, ParamsCheckError); err != nil {
		res.Code = ParamsError
		res.Msg = err.Error()
		return res, nil
	}
	pwd := utils.Encode(form.Password)
	user, token, err := models.LoginVerify(form.Username, pwd)
	if err != nil {
		res.Code = LoginFailed
		res.Msg = err.Error()
		return res, nil
	}
	protocol.Marshal(res, map[string]interface{}{
		"token": token,
		"user":  user.AsheUserJson(),
	})
	res.Msg = "登录成功"
	return res, nil
}

func (*UserApi) Edit(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	user, err := protocol.FetchUserInfo(ctx)
	if err != nil {
		res.Code = TokenParseFailed
		res.Msg = err.Error()
		return res, nil
	}
	data := new(EditForm)
	if err = protocol.Unmarshal(in, data); err != nil {
		res.Code = ParamsError
		res.Msg = ParamsInValid
		return res, nil
	}
	// 校验参数
	if err := check.Check(data, ParamsCheckError); err != nil {
		res.Code = ParamsError
		res.Msg = err.Error()
		return res, nil
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
	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Errorf("获取header失败")
	}
	fmt.Println(incomingContext.Get("host"))
	data := new(models.TUserLog)
	if err := protocol.Unmarshal(in, data); err != nil {
		res.Code = ParamsError
		res.Msg = ParamsInValid
		return res, nil
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
