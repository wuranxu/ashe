package models

import "ashe/exception"

type AsheUser struct {
	ID    uint   `gorm:"primary_key";json:"id"`
	Name  string `gorm:"type:varchar(12)";json:"name"`
	Email string `gorm:"type:varchar(64)";json:"email"`
}

func (u *AsheUser) Add() error {
	return Conn.Insert(u)
}

func (u *AsheUser) Remove(where ...interface{}) error {
	if len(where) == 0 || u.ID == 0 {
		return exception.DangrousDelete
	}
	return Conn.Delete(u, where...)
}

func (u *AsheUser) TableName() string {
	return "ashe_user"
}

func NewUser(name, email string) *AsheUser {
	return &AsheUser{
		Name:  name,
		Email: email,
	}
}
