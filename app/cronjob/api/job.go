package api

import (
	"ashe/app/cronjob/code"
	"ashe/app/cronjob/models"
	"ashe/library/cronjob"
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
func (j *Job) unmarshalData(in *protocol.Request, data interface{}) error {
	if err := json.Unmarshal([]byte(in.RequestJson), data); err != nil {
		return err
	}
	return nil
}

func (j *Job) Add(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	jb := new(models.AsheJob)
	if err := j.unmarshalData(in, jb); err != nil {
		res.Code = code.JobMarshalFail
		res.Msg = err.Error()
		return res, nil
	}
	if err := models.NewAsheJob(jb.Name, jb.Command, jb.IP, "洗澡狗", 22); err != nil {
		res.Code = code.JobAddFail
		res.Msg = err.Error()
		return res, nil
	}
	res.Msg = code.InsertSuccess
	return res, nil

}

func (j *Job) Del(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	jb := new(models.AsheJob)
	if err := j.unmarshalData(in, &jb); err != nil {
		res.Code = code.JobMarshalFail
		res.Msg = err.Error()
		return res, nil
	}
	if err := models.DelJob(jb.ID); err != nil {
		res.Code, res.Msg = code.JobDeleteFail, err.Error()
		return res, nil
	}
	res.Msg = code.DeleteSuccess
	return res, nil

}

func (j *Job) List(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	pg := &struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}{}
	if err := j.unmarshalData(in, pg); err != nil {
		res.Msg, res.Code = err.Error(), code.PageError
		return res, nil
	}
	jbs, total, err := models.GetJobList(pg.Page, pg.PageSize)
	if err != nil {
		res.Msg, res.Code = err.Error(), code.PageError
		return res, nil
	}
	res.ResultJson, res.Msg = getRes(jbs, total), code.GetListSuccess
	return res, nil
}

func (j *Job) TestAssert(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	return protocol.Call("assert", "equal", &protocol.Request{
		RequestJson: `{"exp": 2, "act": 3, "msg": "呀屎啦你"}`,
	})
}

func (j *Job) Sync(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	if err := models.Sync(); err != nil {
		res.Code, res.Msg = code.SyncError, err.Error()
		return res, nil
	}
	res.Msg = code.SyncSuccess
	return res, nil
}

//func (j *Job) Search(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
//
//}
//func (j *Job) Del(ctx context.Context, in *protocol.Request, opts ...grpc.CallOption) (*protocol.Response, error) {
//
//}

func getRes(jbs interface{}, total int) string {
	mp := map[string]interface{}{
		"jobs": jbs, "total": total,
	}
	b, _ := json.Marshal(mp)
	return string(b)
}
