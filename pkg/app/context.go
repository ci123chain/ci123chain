package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/cmd"
	"github.com/ci123chain/ci123chain/pkg/config"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ci123chain/ci123chain/pkg/node"
	"github.com/ci123chain/ci123chain/pkg/util"
	val "github.com/ci123chain/ci123chain/pkg/validator"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/rand"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

//const (
//	flagMasterDomain   = "master_domain"
//	flagMasterPort	   = "master_port"
//	defaultMasterPort  = "80"
//	flagConfig         = "config" //config.toml
//	defaultConfigFilePath = "config.toml"
//	defaultConfigPath  = "config"
//	defaultDataPath    = "data"
//)

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
	root := viper.GetString(util.HomeFlag)
	c, err := config.GetConfig(root)
	if err == config.ErrConfigNotFound {
		master := viper.GetString(util.FlagMasterDomain)
		if len(master) != 0 {
			if os.Getenv("IDG_APPID") == "" {
				return errors.New("Can't use master domain in normal environment")
			}
			c, err = configFollowMaster(master, root)
			if err != nil {
				return err
			}
		} else {
			configEnv := viper.GetString(util.FlagConfig)
			if len(configEnv) != 0 {
				os.MkdirAll(filepath.Join(root, util.DefaultConfigPath), os.ModePerm)
				os.MkdirAll(filepath.Join(root, util.DefaultDataPath), os.ModePerm)
				configBytes, _ := base64.StdEncoding.DecodeString(configEnv)
				ioutil.WriteFile(filepath.Join(root, util.DefaultConfigPath, util.DefaultConfigFilePath), configBytes, os.ModePerm)
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
	port := viper.GetString(util.FlagMasterPort)
	if len(port) == 0 {
		port = util.DefaultMasterPort
	}
	resp, err := http.Get("http://"+ master + ":" + port + "/exportConfig")
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
	var configFiles cmd.ConfigFiles
	err = json.Unmarshal(res, &configFiles)
	if err != nil {
		return nil, err
	}

	c.P2P.PersistentPeers = configFiles.NodeID + "@" + master + ":7443@tls"

	config.SaveConfig(c)

	ioutil.WriteFile(c.GenesisFile(), configFiles.GenesisFile, os.ModePerm)

	var valKey ed25519.PrivKey
	validator := ed25519.GenPrivKey()
	cdc := amino.NewCodec()
	keyByte, err := cdc.MarshalJSON(validator)
	if err != nil {
		return nil, err
	}
	validatorKey := string(keyByte[1:len(keyByte)-1])
	privStr := fmt.Sprintf(`{"type":"%s","value":"%s"}`, ed25519.PrivKeyName, validatorKey)
	cdc = types.GetCodec()
	err = cdc.UnmarshalJSON([]byte(privStr), &valKey)
	if err != nil {
		return nil, err
	}

	pv := val.GenFilePV(
		c.PrivValidatorKeyFile(),
		c.PrivValidatorStateFile(),
		valKey,
	)

	//create node_key.json
	_, err = node.GenNodeKeyByPrivKey(c.NodeKeyFile(), pv.Key.PrivKey)
	if err != nil {
		return nil, err
	}

	return c, nil
}