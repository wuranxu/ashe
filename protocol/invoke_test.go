package protocol

import (
	"fmt"
	"testing"
)

func TestGrpcClient_InvokeWithToken(t *testing.T) {
	rpc, err := NewGrpcClient("user", "Login")
	if err != nil {
		t.Fatal(rpc)
	}
	res, err := rpc.InvokeWithToken(&Request{
		RequestJson:          `{"d": "2"}`,
	}, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZW1haWwiOiI2MTk0MzQxNzZAcXEuY29tIiwibmFtZSI6IuadjumAjemBpSJ9.NFSPADH1lPHbqXl8TKAZVFqGfaT4Ao7q9mFFFp1E5f8")
	fmt.Println(res, err)
}