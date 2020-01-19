package logger

import (
	"github.com/spf13/viper"
	"github.com/tanhuiya/ci123chain/pkg/logger/spliter"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	"io"
	"os"
	"path/filepath"
)

type Logger = log.Logger

var logger Logger

var fileLogger spliter.FileLogger

func GetDefaultLogger(lv string) Logger {
	fname := "dailylog"
	logDir := os.ExpandEnv(filepath.Join(viper.GetString(cli.HomeFlag) , "logs"))
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, os.ModePerm)
		os.Chmod(logDir, os.ModePerm)
	}
	fullFileName := filepath.Join(logDir, fname)
	file, err := os.OpenFile(fullFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.NewTMLogger(log.NewSyncWriter(io.MultiWriter(os.Stdout, file)))

	logger, err = tmflags.ParseLogLevel(lv, logger, cfg.DefaultLogLevel())
	if err != nil {
		panic(err)
	}
	if viper.GetBool(cli.TraceFlag) {
		logger = log.NewTracingLogger(logger)
	}

	fileLogger = spliter.NewFileLogger(logDir, fname)
	logger = logger.With("module", "main")
	return logger
}

func SetLogger(lg Logger) {
	if logger == nil {
		logger = lg
	}
}

func GetLogger() Logger {
	return logger
}

