package app

import (
	"encoding/base64"
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/config"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ci123chain/ci123chain/pkg/node"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/rand"
	pvm "github.com/tendermint/tendermint/privval"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	flagMasterDomain   = "master_domain"
	flagMasterPort	   = "master_port"
	defaultMasterPort  = "443"
	flagConfig         = "config" //config.toml
	defaultConfigFilePath = "config.toml"
	defaultConfigPath  = "config"
	defaultDataPath    = "data"
	flagValidatorKey   = "VALIDATOR_KEY"
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
		master := viper.GetString(flagMasterDomain)
		if len(master) != 0 {
			//if os.Getenv("IDG_APPID") == "" {
			//	return errors.New("Can't use master domain in normal environment")
			//}
			c, err = configFollowMaster(master, root)
			if err != nil {
				return err
			}
		} else {
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
				c, err = config.CreateConfig(rand.Str(8), root)
				if err != nil {
					return config.ErrGetConfig
				}
				config.SaveConfig(c)
			}
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

func configFollowMaster(master, root string) (*cfg.Config, error){
	port := viper.GetString(flagMasterPort)
	if len(port) == 0 {
		port = defaultMasterPort
	}
	prefix := util.DefaultHTTP
	if os.Getenv(util.IDG_APPID) != "" {
		prefix = util.DefaultHTTPS
	}
	resp, err := http.Get(prefix + master + ":" + port + "/exportConfig")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err :=ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c, err := config.CreateConfig(rand.Str(8), root)
	if err != nil {
		return nil, err
	}
	var configFiles types.ConfigFiles
	err = json.Unmarshal(res, &configFiles)
	if err != nil {
		return nil, err
	}

	if os.Getenv(util.IDG_APPID) != "" {
		c.P2P.PersistentPeers = configFiles.NodeID + "@" + master + ":7443@tls"
	} else {
		c.P2P.PersistentPeers = configFiles.NodeID + "@" + master + ":26656"
	}

	config.SaveConfig(c)

	if err := ioutil.WriteFile(c.GenesisFile(), configFiles.GenesisFile, os.ModePerm); err != nil {
		panic(err)
	}
	
	pv := pvm.LoadOrGenFilePV(c.PrivValidatorKeyFile(), c.PrivValidatorStateFile())
	//create node_key.json
	_, err = node.GenNodeKeyByPrivKey(c.NodeKeyFile(), pv.Key.PrivKey)
	if err != nil {
		return nil, err
	}

	return c, nil
}