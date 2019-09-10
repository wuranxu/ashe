package api

import (
	"ashe/app/cronjob/code"
	"ashe/library/cronjob"
	"ashe/models"
	"ashe/protocol"
	"context"
	"encoding/json"
)

type Job struct {
}

func (j *Job) unmarshal(in *protocol.Request) (*cronjob.Job, error) {
	jb := new(cronjob.Job)
	if err := json.Unmarshal([]byte(in.RequestJson), jb); err != nil {
		return nil, err
	}
	return jb, nil
}

func (j *Job) Add(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	jb, err := j.unmarshal(in)
	if err != nil {
		res.Code = code.JobMarshalFail
		res.Msg = err.Error()
		return nil, err
	}
	if err = models.NewAsheJob(jb.Name, jb.Command, jb.IP, "洗澡狗", 22); err != nil {
		res.Code = code.JobAddFail
		res.Msg = err.Error()
		return res, nil
	}
	res.Msg = code.InsertSuccess
	return res, nil

}

func (j *Job) Del(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {	res := new(protocol.Response)
	jb, err := j.unmarshal(in)
	if err != nil {
		res.Code = code.JobMarshalFail
		res.Msg = err.Error()
		return res, err
	}
	if err = models.DelJob(jb.ID); err != nil {
		res.Code, res.Msg = code.JobDeleteFail, err.Error()
	}
	res.Msg = code.DeleteSuccess
	return res, nil

}
//func (j *Job) Search(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
//
//}
//func (j *Job) Del(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
//
//}
