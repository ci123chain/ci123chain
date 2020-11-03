package app

import (
	"encoding/base64"
	"github.com/spf13/viper"
	"github.com/ci123chain/ci123chain/pkg/config"
	"github.com/ci123chain/ci123chain/pkg/logger"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	flagConfig         = "config" //config.toml
	defaultConfigFilePath = "config.toml"
	defaultConfigPath  = "config"
	defaultDataPath    = "data"
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
		configEnv := viper.GetString(flagConfig)
		if len(configEnv) != 0 {
			os.MkdirAll(filepath.Join(root, defaultConfigPath), os.ModePerm)
			os.MkdirAll(filepath.Join(root, defaultDataPath), os.ModePerm)
			configBytes, _ := base64.StdEncoding.DecodeString(configEnv)
			ioutil.WriteFile(filepath.Join(root, defaultConfigPath, defaultConfigFilePath), configBytes, os.ModePerm)
			viper.ReadInConfig()
			c, err = config.GetConfig(root)
			if err != nil {
				return config.ErrGetConfig
			}
		} else {
			c, err = config.CreateConfig(common.RandStr(8), root)
			if err != nil {
				return config.ErrGetConfig
			}
			config.SaveConfig(c)
		}
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