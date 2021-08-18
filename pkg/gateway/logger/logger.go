package logger

import (
	"fmt"
	tmlogger "github.com/tendermint/tendermint/libs/log"
	"os"
)

//func Init() {
//	log.InitOneitfarmLogger()
//}

var logger = tmlogger.NewTMLogger(tmlogger.NewSyncWriter(os.Stdout))

func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Error(msg)
}

func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Warn(msg)
}

func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Info(msg)
}

func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Debug(msg)
}
