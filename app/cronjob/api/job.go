package api

import (
	"ashe/library/cronjob"
	"ashe/models"
	"ashe/protocol"
	"context"
	"encoding/json"
	"fmt"
)

type Job struct {
}

func (j *Job) unmarshal(in *protocol.Request) (*cronjob.Job, error) {
	jb := new(cronjob.Job)
	if err := json.Unmarshal([]byte(in.RequestJson), jb); err != nil {
		return nil, err
	}
	fmt.Println(jb)
	return jb, nil
}

func (j *Job) Add(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	jb, err := j.unmarshal(in)
	if err != nil {
		return nil, err
	}
	if err = models.NewAsheJob(jb.Name, jb.Command, jb.IP, "洗澡狗", 22); err != nil {
		res.Code = 10001
		res.Msg = err.Error()
		return res, nil
	}
	res.Code = 0
	res.Msg = "插入成功"
	return res, nil

}

//func (j *Job) Edit(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
//
//}
//func (j *Job) Search(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
//
//}
//func (j *Job) Del(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
//
//}
