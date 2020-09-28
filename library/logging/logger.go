package logging

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type CommonFormat struct {
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

//格式详情
func (s *CommonFormat) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05.000")
	var file string
	var length int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		length = entry.Caller.Line
	}
	//fmt.Println(entry.Data)
	msg := fmt.Sprintf("[%s] [%s] [%s:%d][GOID:%d] %s\n", strings.ToUpper(entry.Level.String()), timestamp, file, length, getGID(), entry.Message)
	return []byte(msg), nil
}

func NewLog(filename string) *logrus.Logger {
	log := logrus.New()
	// log to console and file
	f, err := os.OpenFile(fmt.Sprintf("logs/%s.log", filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetReportCaller(true)
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&CommonFormat{})
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	return log
}
