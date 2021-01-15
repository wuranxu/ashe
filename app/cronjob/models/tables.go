package models

import (
	"github.com/wuranxu/library/dao"
)

var Conn *dao.Cursor

var Tables = []interface{}{
	&AsheJob{},
}