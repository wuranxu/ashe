package models

import (
	exp "ashe/exception"
	"ashe/library/cronjob"
	"ashe/library/database"
	"time"
)

var (
	InsertError = exp.ErrString("添加job出错")
)

type AsheJob struct {
	cronjob.Job
	Pid     uint `gorm:"type:int(8)" json:"pid"`               // 测试计划id
	Deleted bool `gorm:"type:boolean;not null" json:"deleted"` // 是否被删除
}

func (a *AsheJob) TableName() string { return "ashe_job" }

// 同步数据到redis
func (a *AsheJob) SyncToRedis() error {
	conn := cronjob.Pool.Get()
	defer conn.Close()
	return cronjob.SetJobToRedis(&a.Job)
}

// 获取job列表
func GetJobList(page, pageSize int) ([]*cronjob.Job, int, error) {
	jobs, total, err := cronjob.GetJobList(page, pageSize)
	if err == cronjob.PageOutOfRange || err == cronjob.PageSizeTooLong {
		// 并非redis故障
		return jobs, total, err
	}
	if err == nil {
		return jobs, total, nil
	}
	// 从db读取job
	total, err = Conn.FindPaginationAndOrder(page, pageSize, `create_time desc`, &jobs, `deleted = false`)
	return jobs, total, err
}

// 添加job
func NewAsheJob(name, command, ip, user string, userId uint, pid ...uint) error {
	var planId uint
	if len(pid) > 0 {
		planId = pid[0]
	}
	job := &AsheJob{
		Job: cronjob.Job{
			Name:       name,
			Command:    command,
			IP:         ip,
			Editor:     user,
			Creator:    user,
			Uid:        userId,
			EditorId:   userId,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		},
		Pid:     planId,
		Deleted: false,
	}
	if err := Conn.Insert(job); err != nil {
		return InsertError.Error(err)
	}
	// 更新redis
	return job.SyncToRedis()
}

// 删除job
func DelJob(id uint) error {
	job := &AsheJob{Job: cronjob.Job{ID: id}}
	_, err := Conn.Updates(job, database.Columns{"name": "吴冉旭不爱"})
	if err != nil {
		return err
	}
	err = cronjob.DelJobFromRedis(id)
	return err
}

//// 批量更新数据到redis
//func UpdateJobList(conn redis.RedisCon, data []*cronjob.Job) error {
//
//}
