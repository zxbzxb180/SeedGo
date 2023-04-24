package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
)

var Logger *logrus.Logger

const (
	// LogName 日志文件名
	LogName = "all"
	// LogSuffix 日志文件后缀
	LogSuffix = ".log"
	// LogSize 单个日志文件大小，单位MB
	LogSize = 50
	// LogBackup 日志文件个数
	LogBackup = 10
	// LogDate 日志文件最大天数
	LogDate = 7
)

func init() {
	Logger = logrus.New()
	//Logger.Formatter = &logrus.TextFormatter{
	//	FullTimestamp:          true,
	//	ForceColors:            true,
	//	DisableLevelTruncation: true,
	//	TimestampFormat:        "2006-01-02 15:04:05",
	//}
	fmt.Println(path.Join("./logs/", LogName+LogSuffix))
	logconf := &lumberjack.Logger{
		Filename:   path.Join("./logs/", LogName+LogSuffix),
		MaxSize:    LogSize,   // 日志文件大小，单位是 MB
		MaxBackups: LogBackup, // 最大过期日志保留个数
		MaxAge:     LogDate,   // 保留过期文件最大时间，单位 天
		Compress:   false,     // 是否压缩日志，默认是不压缩。这里设置为true，压缩日志
	}
	Logger.SetOutput(io.MultiWriter(logconf, os.Stdout))
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.DebugLevel)
	//Logger.SetOutput(os.Stdout)
	//logrus.SetOutput(io.MultiWriter(os.Stdout, hook))

}
