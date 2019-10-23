package logger

import (
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	cfg "github.com/tendermint/tendermint/config"
	"os"
)

type Logger = log.Logger

var logger Logger

func GetDefaultLogger(lv string) Logger {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger, err := tmflags.ParseLogLevel(lv, logger, cfg.DefaultLogLevel())
	if err != nil {
		panic(err)
	}
	if viper.GetBool(cli.TraceFlag) {
		logger = log.NewTracingLogger(logger)
	}
	return logger.With("module", "main")
}

func SetLogger(lg Logger) {
	if logger == nil {
		logger = lg
	}
}

func GetLogger() Logger {
	return logger
}