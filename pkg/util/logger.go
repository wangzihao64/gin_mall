package util

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

var LogrousObj *logrus.Logger

func init() {
	src, _ := setOutPutFile()
	if LogrousObj != nil {
		LogrousObj.Out = src
		return
	}
	//实例化
	logger := logrus.New()
	logger.Out = src                   //设置输出
	logger.SetLevel(logrus.DebugLevel) //设置日志级别
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	LogrousObj = logger
}
func setOutPutFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil { //获取工作目录
		logFilePath = dir + "/logs"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	logFileName := now.Format("2006-01-02") + ".log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	_, err = os.Stat(fileName)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(fileName, 0777); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	return src, nil
}
