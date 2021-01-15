package db

import (
	"github.com/wuranxu/library/conf"
	"github.com/wuranxu/library/dao"
	"log"
)

func Init(tables []interface{}) *dao.Cursor {
	conn, err := dao.NewConnect(conf.Conf.Database)
	if err != nil {
		log.Fatalf("连接数据库失败, error: %v", err)
	}
	if tables != nil {
		for _, data := range tables {
			conn.AutoMigrate(data)
		}
	}
	conn.LogMode(conf.Conf.Database.LogMode)
	return conn
}
