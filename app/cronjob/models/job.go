package models

import (
	exp "ashe/exception"
	"ashe/library/cronjob"
	"ashe/library/database"
	"time"
)

var (
	InsertError = exp.ErrString("添加job出错")
	DeleteError = exp.ErrString("删除job出错")
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

// 更新redis缓存
func deleteAllJob() error {
	return cronjob.DeleteAllJobs()
}

// 同步缓存
func Sync() error {
	if err := deleteAllJob(); err != nil {
		return err
	}
	jobs, err := getAllJobFromDb()
	if err != nil {
		return err
	}
	for _, j := range jobs {
		if err = cronjob.SetJobToRedis(&j.Job); err != nil {
			continue
		}
	}
	return err
}

// 获取所有job
func getAllJobFromDb() ([]AsheJob, error) {
	// 从db读取job
	jobs := make([]AsheJob, 0, 100)
	if err := Conn.Find(&jobs, `deleted = false`).Error; err != nil {
		return jobs, err
	}
	return jobs, nil
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
	asheJobs := make([]AsheJob, 0, 10)
	total, err = Conn.FindPaginationAndOrder(page, pageSize, `create_time desc`, &asheJobs, `deleted = false`)
	return transJobs(asheJobs), total, err
}

// 转换job
func transJobs(jbs []AsheJob) []*cronjob.Job {
	result := make([]*cronjob.Job, 0, len(jbs))
	for _, j := range jbs {
		result = append(result, &j.Job)
	}
	return result
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
		return InsertError.New(err)
	}
	// 更新redis
	return job.SyncToRedis()
}

// 删除job
func DelJob(id uint) error {
	job := &AsheJob{Job: cronjob.Job{ID: id}}
	n, err := Conn.Updates(job, database.Columns{"deleted": true})
	if err != nil || n == 0 {
		return DeleteError.New(err)
	}
	err = cronjob.DelJobFromRedis(id)
	return err
}
