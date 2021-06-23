package logger

import (
	"github.com/tendermint/tendermint/libs/log"
	logger "gitlab.oneitfarm.com/bifrost/cilog/v2"
)

func Init() {
	log.InitOneitfarmLogger()
}

func Error(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func Info(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}
