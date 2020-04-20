package app

import (
	"github.com/spf13/viper"
	"github.com/ci123chain/ci123chain/pkg/config"
	"github.com/ci123chain/ci123chain/pkg/logger"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

type Context struct {
	Config *cfg.Config
	Logger log.Logger
}


func NewDefaultContext() *Context {
	return NewContext(
		cfg.DefaultConfig(),
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	)
}

func NewContext(config *cfg.Config, logger log.Logger) *Context {
	return &Context{config, logger}
}

func SetupContext(ctx *Context, level string) error {
	root := viper.GetString(cli.HomeFlag)
	c, err := config.GetConfig(root)
	if err == config.ErrConfigNotFound {
		c, err = config.CreateConfig(common.RandStr(8), root)
		if err != nil {
			return config.ErrGetConfig
		}
		config.SaveConfig(c)
	}
	if err != nil {
		return config.ErrGetConfig
	}
	c.SetRoot(root)
	lg := logger.GetDefaultLogger(level)
	ctx.Config = c
	ctx.Logger = lg
	return nil
}