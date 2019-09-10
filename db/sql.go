package db

import (
	"ashe/common"
	"ashe/library/database"
	"ashe/models"
	"log"
	"sync"
)

func Init() {
	var (
		err  error
		once sync.Once
	)
	once.Do(func() {
		models.Conn, err = database.NewConnect(common.Conf.Database)
		if err != nil {
			log.Fatalf("连接数据库失败, error: %v", err)
		}
		for _, data := range models.Tables {
			models.Conn.AutoMigrate(data)
		}
		if common.Env == "dev" {
			models.Conn.LogMode(true)
		}
	})
}
