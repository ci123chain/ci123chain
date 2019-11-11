package app

import (
	"github.com/tanhuiya/ci123chain/pkg/config"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tanhuiya/ci123chain/pkg/logger"
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

func SetupContext(ctx *Context) error {
	root := viper.GetString(tmcli.HomeFlag)
	c, err := config.GetConfig(root)
	if err == config.ErrConfigNotFound {
		c, err = config.CreateConfig(common.RandStr(8), root)
		if err != nil {
			return err
		}
		config.SaveConfig(c)
	}
	if err != nil {
		return err
	}
	c.SetRoot(root)
	lg := logger.GetDefaultLogger(c.LogLevel)
	ctx.Config = c
	ctx.Logger = lg
	return nil
}