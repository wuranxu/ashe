package net

import (
	"fmt"
	"testing"
)

func TestGetLocalIp(t *testing.T) {
	fmt.Println(GetLocalIp())
}
