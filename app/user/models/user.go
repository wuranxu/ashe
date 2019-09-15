package models

import (
	"ashe/app/user/utils"
	"ashe/library/auth"
	"ashe/library/database"
	tm "ashe/library/time"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	UserExistsError = errors.New("用户名或邮箱已被注册")
	LoginFailed     = errors.New("用户名或密码不正确")
	UserInvalid     = errors.New("用户已被删除, 如有疑问请联系管理员")
)

type AsheUser struct {
	ID            int         `gorm:"primary_key;unique" json:"id"`
	Username      string      `gorm:"type:varchar(16);not null;unique" json:"username" validate:"gt=0"` // 用户名
	Password      string      `gorm:"type:varchar(32);not null" json:"password" validate:"gt=0"`
	Nickname      string      `gorm:"type:varchar(32);not null" json:"nickname" validate:"gt=0"` // 昵称
	Email         string      `gorm:"type:varchar(64);not null;unique" json:"email" validate:"gt=0"`
	LastLoginIp   string      `gorm:"type:varchar(15)" json:"remote_ip"`     // 上次登录ip
	LastLoginTime tm.JSONTime `gorm:"type:timestamp" json:"last_login_time"` // 上次登录时间
	CreateTime    tm.JSONTime `gorm:"type:timestamp" json:"create_time"`
	Deleted       bool        `gorm:"type:boolean;default false" json:"deleted"`
}

func (u *AsheUser) TableName() string { return "ashe_user" }

func (u *AsheUser) Register() error {
	var temp AsheUser
	if err := Conn.Find(&temp, "email = ? or username = ?", u.Email, u.Username).Error; err == nil {
		return UserExistsError
	}
	// 加密password
	u.Password = utils.Encode(u.Password)
	u.CreateTime = tm.JSONTime{time.Now()}
	u.LastLoginTime = tm.JSONTime{time.Now()}
	return Conn.Insert(u)
}

func LoginVerify(username, password string) (*AsheUser, string, error) {
	var (
		user  AsheUser
		token string
	)
	if err := Conn.Find(&user, `username = ? and password = ?`, username, password).DB.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &user, token, LoginFailed
		}
		return &user, token, err
	}
	if user.Deleted {
		// 用户已被删除
		return nil, token, UserInvalid
	}
	// 根据用户名和Email以及姓名生成token
	jt := auth.NewJWT()
	token, err := jt.CreateToken(auth.CustomClaims{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.Nickname,
		StandardClaims: jwt.StandardClaims{},
	})
	return &user, token, err
}

func (u *AsheUser) AsheUserJson() interface{}{
	return map[string]interface{}{
		"nickname": u.Nickname,
		"email": u.Email,
		"last_login_ip": u.LastLoginIp,
		"last_login_time": u.LastLoginTime,
		"user_id": u.ID,
	}
}

//用户修改自己的信息
func Edit(nickname, email string, userId int) error {
	_, err := Conn.Updates(&AsheUser{ID: userId}, database.Columns{
		"nickname": nickname, "email": email,
	})
	return err
}
