package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	//Logger.Formatter = &logrus.TextFormatter{
	//	FullTimestamp:          true,
	//	ForceColors:            true,
	//	DisableLevelTruncation: true,
	//	TimestampFormat:        "2006-01-02 15:04:05",
	//}
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetOutput(os.Stdout)

}
