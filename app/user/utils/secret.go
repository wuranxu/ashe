package utils

import (
	"crypto/md5"
	"fmt"
)

const salt = "ashe_user"

// md5加密密码
func Encode(data string) string {
	finalStr := data + salt
	return fmt.Sprintf("%x", md5.Sum([]byte(finalStr)))
}
