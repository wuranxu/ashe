package cronjob

import (
	exp "ashe/exception"
	"encoding/json"
	"fmt"
	rds "github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/wuranxu/library/cache/redis"
	"log"
	"sort"
	"time"
)

const (
	JobPrefix   = "ASHE_CRON_JOB"
	MaxPageSize = 50
)

var (
	PageSizeTooLong = exp.ErrString(fmt.Sprintf("任务列表pageSize不能超过%d条", MaxPageSize))
	PageOutOfRange  = exp.ErrString("job页数超出范围")
	PageError       = exp.ErrString("page/pageSize必须大于0")
)

var (
	Pool *rds.Client
)

func InitRedisConnection(cfg *rds.Options) {
	var err error
	Pool, err = redis.NewClient(cfg)
	if err != nil {
		log.Fatal("connect redis failed, error: ", err)
	}
}

type Job struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"type:varchar(64);not null;unique" json:"name" validate:"gt=0"` // job name
	Command    string    `gorm:"type:text; not null" json:"command" validate:"gt=0"`           // shell command
	IP         string    `gorm:"type:varchar(24)" json:"ip"`                                   // 机器执行ip
	Editor     string    `gorm:"type:varchar(32)" json:"editor"`
	Creator    string    `gorm:"type:varchar(32)" json:"creator"`
	Uid        uint      `gorm:"type:int(8)" json:"uid"`
	EditorId   uint      `gorm:"type:int(8)" json:"editor_id"`
	CreateTime time.Time `gorm:"type:timestamp;not null" json:"create_time"`
	UpdateTime time.Time `gorm:"type:timestamp;not null" json:"update_time"`
}

type sortJob []*Job

func (j sortJob) Len() int {
	return len(j)
}

func (j sortJob) Swap(x, y int) {
	j[x], j[y] = j[y], j[x]
}

func (j sortJob) Less(x, y int) bool {
	// 按照更新时间倒叙排列
	return j[x].UpdateTime.After(j[y].UpdateTime)
}

// single
func SetJobToRedis(j *Job) error {
	result, err := json.Marshal(j)
	if err != nil {
		return err
	}
	if err = Pool.Set(fmt.Sprintf(`%s:%d`, JobPrefix, j.ID), result, 5*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}

// delete job
func DelJobFromRedis(id uint) error {
	if err := Pool.Del(fmt.Sprintf(`%s:%d`, JobPrefix, id)).Err(); err != nil {
		return err
	}
	return nil
}

// batch set
func SetBatchJob(job []*Job) (err error) {
	for _, j := range job {
		if err = SetJobToRedis(j); err != nil {
			return
		}
	}
	return
}

func GetJobFromRedis(id int) (*Job, error) {
	result, err := Pool.Get(fmt.Sprintf(`%s:%d`, JobPrefix, id)).Result()
	if err != nil {
		return nil, err
	}
	return parseJob(result)
}

func parseJob(result string) (*Job, error) {
	jb := new(Job)
	err := json.Unmarshal([]byte(result), jb)
	if err != nil {
		return nil, err
	}
	return jb, nil
}

// 获取所有job keys
func getJobKeys() ([]string, int, error) {
	res, err := Pool.Keys(fmt.Sprintf(`%s:*`, JobPrefix)).Result()
	if err != nil {
		return res, 0, err
	}
	return res, len(res), nil
}

// 获取索引
func getIndex(page, pageSize, total int) (int, int, error) {
	var start, end int
	if page <= 0 || pageSize <= 0 {
		return start, end, PageError
	}
	start = (page - 1) * pageSize
	if start > total {
		return start, end, PageOutOfRange
	}
	end = page*pageSize - 1
	if page*pageSize > total {
		end = total - 1
	}
	return start, end, nil
}

//// TO []interface{}
//func toInterface(data []string) []interface{} {
//	result := make([]interface{}, 0)
//	for _, d := range data {
//		result = append(result, d)
//	}
//	return result
//}

// 删除缓存所有的job
func DeleteAllJobs() error {
	keys, _, err := getJobKeys()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	if err := Pool.Del(keys...).Err(); err != nil {
		logrus.Errorf("删除redis key失败, error: %v", err)
		return err
	}
	return nil
}

// 获取job列表 from redis
func GetJobList(page, pageSize int) (sortJob, int, error) {
	result := make(sortJob, 0)
	if pageSize > MaxPageSize {
		return result, 0, PageSizeTooLong
	}
	keys, total, err := getJobKeys()
	if err != nil {
		return result, total, err
	}
	if total == 0 {
		return result, total, nil
	}
	start, end, err := getIndex(page, pageSize, len(keys))
	if err != nil {
		return result, total, err
	}
	jobs, err := Pool.MGet(keys...).Result()
	for _, j := range jobs {
		var jb *Job
		if jb, err = parseJob(j.(string)); err == nil {
			result = append(result, jb)
			continue
		}
		logrus.Errorf("解析job失败, error: %v", err)
	}
	sort.Sort(result)
	return result[start : end+1], total, nil
}
