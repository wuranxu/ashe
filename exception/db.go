package exception

var (
	UnSupportedDatabase = ErrString("database is not supported now")
	DangrousDelete      = ErrString("delete operation without condition is not allowed")
)
