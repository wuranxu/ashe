package cronjob

import (
	exp "ashe/exception"
	"ashe/library/cache/redis"
	"encoding/json"
	"fmt"
	rds "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"log"
	"sort"
	"sync"
	"time"
)

const (
	JobPrefix   = "ASHE_CRON_JOB"
	MaxPageSize = 50
)

var (
	JobShouldNotNull = exp.ErrString("job不能为空")
	PageSizeTooLong  = exp.ErrString(fmt.Sprintf("任务列表pageSize不能超过%d条", MaxPageSize))
	PageOutOfRange   = exp.ErrString("job页数超出范围")
	PageError        = exp.ErrString("page/pageSize必须大于0")
)

var (
	Pool *rds.Pool
)

func InitRedisConnection(cfg redis.RedisCliInfo) {
	var once sync.Once
	once.Do(func() {
		if Pool = redis.NewCliPool(cfg); Pool == nil {
			log.Fatalf("redis连接失败")
		}
	})
}

type Job struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"type:varchar(64);not null;unique" json:"name"` // job name
	Command    string    `gorm:"type:text; not null" json:"command"`           // shell command
	IP         string    `gorm:"type:varchar(24)" json:"ip"`                   // 机器执行ip
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
	conn := Pool.Get()
	defer conn.Close()
	result, err := json.Marshal(j)
	if err != nil {
		return err
	}
	_, err = rds.String(conn.Do("SET", fmt.Sprintf(`%s:%d`, JobPrefix, j.ID), result))
	if err != nil {
		return err
	}
	return nil
}

// delete job
func DelJobFromRedis(id uint) error {
	conn := Pool.Get()
	defer conn.Close()
	_, err := rds.Bool(conn.Do("DEL", fmt.Sprintf(`%s:%d`, JobPrefix, id)))
	if err != nil {
		return err
	}
	return nil
}

// batch set
func SetBatchJob(job []*Job) (err error) {
	conn := Pool.Get()
	defer conn.Close()
	for _, j := range job {
		if err = SetJobToRedis(j); err != nil {
			return
		}
	}
	return
}

func GetJobFromRedis(id int) (*Job, error) {
	conn := Pool.Get()
	defer conn.Close()
	result, err := rds.String(conn.Do("GET", fmt.Sprintf(`%s:%d`, JobPrefix, id)))
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
func getJobKeys(conn rds.Conn) ([]string, int, error) {
	res, err := rds.Strings(conn.Do("KEYS", fmt.Sprintf(`%s:*`, JobPrefix)))
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

// TO []interface{}
func toInterface(data []string) []interface{} {
	result := make([]interface{}, 0)
	for _, d := range data {
		result = append(result, d)
	}
	return result
}

// 删除缓存所有的job
func DeleteAllJobs() error {
	conn := Pool.Get()
	defer conn.Close()
	var err error
	keys, _, err := getJobKeys(conn)
	if err != nil {
		return err
	}
	for _, k := range keys {
		if _, err = conn.Do("DEL", k); err != nil {
			logrus.Errorf("删除redis key失败, key: %v, error: %v", k, err)
			continue
		}
	}
	return err
}

// 获取job列表 from redis
func GetJobList(page, pageSize int) (sortJob, int, error) {
	conn := Pool.Get()
	defer conn.Close()
	result := make(sortJob, 0)
	if pageSize > MaxPageSize {
		return result, 0, PageSizeTooLong
	}
	keys, total, err := getJobKeys(conn)
	if err != nil {
		return result, total, err
	}
	start, end, err := getIndex(page, pageSize, len(keys))
	if err != nil {
		return result, total, err
	}
	jobs, err := rds.Strings(conn.Do("MGET", toInterface(keys)...))
	for _, j := range jobs {
		var jb *Job
		if jb, err = parseJob(j); err == nil {
			result = append(result, jb)
			continue
		}
		logrus.Errorf("解析job失败, error: %v", err)
	}
	sort.Sort(result)
	return result[start : end+1], total, nil
}
