package handler

import (
	"github.com/kataras/iris"
	"github.com/wuranxu/library/auth"
	"net/http"
	"strings"
)

const (
	SignKey      = "ASHE"
	AuthFailCode = 103
)

// http response
type Res struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewRes(code int, msg string) *Res {
	return &Res{code, msg}
}

func getUserInfo(ctx iris.Context) (*auth.CustomClaims, error) {
	token := ctx.GetHeader("Authorization")
	if s := strings.Split(token, " "); len(s) == 2 {
		token = s[1]
	}
	j := auth.NewJWT(SignKey)
	return j.ParseToken(token)
}

func Auth(handler func(ctx iris.Context, userInfo *auth.CustomClaims)) iris.Handler {
	return func(ctx iris.Context) {
		claims, err := getUserInfo(ctx)
		if err != nil {
			ctx.StatusCode(http.StatusOK)
			ctx.JSON(NewRes(103, err.Error()))
			return
		}
		handler(ctx, claims)
	}
}

func Authorization(ctx iris.Context) (*auth.CustomClaims, error) {
	return getUserInfo(ctx)
}

func AuthMail(handler func(ctx iris.Context, email string, userId int)) iris.Handler {
	return func(ctx iris.Context) {
		claims, err := getUserInfo(ctx)
		if err != nil {
			ctx.StatusCode(http.StatusOK)
			ctx.JSON(NewRes(AuthFailCode, err.Error()))
			return
		}
		handler(ctx, claims.Email, claims.ID)
	}
}

func AuthName(handler func(ctx iris.Context, name string, userId int)) iris.Handler {
	return func(ctx iris.Context) {
		claims, err := getUserInfo(ctx)
		if err != nil {
			ctx.StatusCode(http.StatusOK)
			ctx.JSON(NewRes(AuthFailCode, err.Error()))
			return
		}
		handler(ctx, claims.Name, claims.ID)
	}
}
