package exception

import "errors"

var (
	UnSupportedDatabase = errors.New("database is not supported now")
	DangrousDelete  = errors.New("delete operation without condition is not allowed")
)