package db

import (
	"ashe/common"
	"ashe/library/database"
	"log"
)

func Init(tables []interface{}) *database.Cursor{
	conn, err := database.NewConnect(common.Conf.Database)
	if err != nil {
		log.Fatalf("连接数据库失败, error: %v", err)
	}
	if tables != nil {
		for _, data := range tables {
			conn.AutoMigrate(data)
		}
	}
	if common.Env == "dev" {
		conn.LogMode(true)
	}
	return conn
}
