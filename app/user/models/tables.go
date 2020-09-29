package models

import (
	"ashe/library/database"
)

var Conn *database.Cursor

var Tables = []interface{}{
	&AsheUser{},
	&TUserLog{},
}