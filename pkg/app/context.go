package app

import (
	"encoding/base64"
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/cmd"
	"github.com/ci123chain/ci123chain/pkg/node"
	"github.com/spf13/viper"
	"github.com/ci123chain/ci123chain/pkg/config"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	val "github.com/ci123chain/ci123chain/pkg/validator"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"fmt"
)

const (
	flagMasterDomain   = "master_domain"
	flagMasterPort	   = "master_port"
	defaultMasterPort  = "80"
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
		master := viper.GetString(flagMasterDomain)
		if len(master) != 0 {
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
				c, err = config.CreateConfig(common.RandStr(8), root)
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
	resp, err := http.Get("http://"+ master + ":" + port + "/exportConfig")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err :=ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c, err := config.CreateConfig(common.RandStr(8), root)
	if err != nil {
		return nil, err
	}
	var configFiles cmd.ConfigFiles
	err = json.Unmarshal(res, &configFiles)
	if err != nil {
		return nil, err
	}

	c.P2P.PersistentPeers = configFiles.NodeID + "@" + master + ":26656"
	//c.P2P.PersistentPeers = configFiles.NodeID + "@" + master + ":26656@tls"

	config.SaveConfig(c)

	ioutil.WriteFile(c.GenesisFile(), configFiles.GenesisFile, os.ModePerm)

	var valKey secp256k1.PrivKeySecp256k1
	validator := secp256k1.GenPrivKey()
	cdc := amino.NewCodec()
	keyByte, err := cdc.MarshalJSON(validator)
	if err != nil {
		return nil, err
	}
	validatorKey := string(keyByte[1:len(keyByte)-1])
	privStr := fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PrivKeyAminoName, validatorKey)
	cdc = types.MakeCodec()
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