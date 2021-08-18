package util

import (
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"os"
)

func CallBack(err error) {
	//logger.Init()
	logger.Error("get info from remote discovery failed", "error", err.Error())
	os.Exit(1)
}
