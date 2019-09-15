package auth

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"net/http"
	"strings"
	"time"
)

type JWT struct {
	SigningKey []byte
}

// http response
type Res struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewRes(code int, msg string) *Res {
	return &Res{code, msg}
}

var (
	TokenExpired     = errors.New("身份认证已过期, 请重新登录")
	TokenNotValidYet = errors.New("身份认证未生效")
	TokenMalformed   = errors.New("登录信息已失效, 请重新登录")
	TokenInvalid     = errors.New("身份认证不合法")
	signKey          = "ashe"
	AuthFailCode     = 103
)

type CustomClaims struct {
	ID                 int    `json:"id"`    // userId
	Email              string `json:"email"` // user_email
	Name               string `json:"name"`  // username
	jwt.StandardClaims `json:"jwt,omitempty"`
}

func (c *CustomClaims) Marshal() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}
func GetSignKey() string {
	return signKey
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(t string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(t, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

func (j *JWT) RefreshToken(tokenStr string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

func getUserInfo(ctx iris.Context) (*CustomClaims, error) {
	token := ctx.GetHeader("Authorization")
	if s := strings.Split(token, " "); len(s) == 2 {
		token = s[1]
	}
	j := NewJWT()
	return j.ParseToken(token)
}

func Auth(handler func(ctx iris.Context, userInfo *CustomClaims)) iris.Handler {
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
