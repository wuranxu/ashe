package api

import (
	"ashe/app/cronjob/code"
	"ashe/app/cronjob/models"
	"ashe/library/check"
	"ashe/protocol"
	"context"
)

type Job struct {
}

func (j *Job) unmarshalData(in *protocol.Request, data interface{}) error {
	if err := protocol.Unmarshal(in, data); err != nil {
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
	// 校验参数
	if err := check.Check(jb, code.ParamsCheckError); err != nil {
		res.Code = code.JobParseFail
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
	if pg.Page == 0 {
		pg.Page = 1
	}
	if pg.PageSize == 0 {
		pg.PageSize = 8
	}
	jbs, total, err := models.GetJobList(pg.Page, pg.PageSize)
	if err != nil {
		res.Msg, res.Code = err.Error(), code.PageError
		return res, nil
	}
	protocol.Marshal(res, getRes(jbs, total))
	res.Msg = code.GetListSuccess
	return res, nil
}

func (j *Job) TestAssert(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	return protocol.Call("assert", "equal", &protocol.Request{
	})
}

func (j *Job) Sync(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	res := new(protocol.Response)
	if _, err := models.Sync(); err != nil {
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

func getRes(jbs interface{}, total int) map[string]interface{} {
	mp := map[string]interface{}{
		"jobs": jbs, "total": total,
	}
	return mp
}
