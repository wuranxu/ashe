package cronjob

import (
	"ashe/library/cache/redis"
	"fmt"
	"testing"
	"time"
)

func TestSetJobToRedis(t *testing.T) {
	job := []*Job{{
		ID:         435443553,
		Name:       "李逍遥",
		Command:    "echo hel2lo",
		IP:         "128.2.2.1",
		UpdateTime: time.Now(),
	},{
		ID:         332234,
		Name:       "李逍遥3",
		Command:    "echo he2l2lo",
		IP:         "128.2.2.1",
		UpdateTime: time.Now(),
	}}
	cl := redis.RedisCliInfo{
		Server:      "106.13.173.14:6379",
		Password:    "ashetest",
		Port:        6379,
		Db:          1,
		MaxIdle:     50,
		MaxActive:   50,
		IdleTimeout: 5000,
	}
	if pl := redis.NewCliPool(cl); pl != nil {
		conn := pl.Get()
		defer conn.Close()
		fmt.Println(SetBatchJob(job))
		//fmt.Println(getJobKeys(conn))
		fmt.Println(GetJobList(1, 5))
		//fmt.Println(GetJobFromRedis(job.ID,conn))
	}

}