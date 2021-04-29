package cmd

import (
	"encoding/json"
	"time"

	//"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/node"
	otypes "github.com/ci123chain/ci123chain/pkg/order/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ci123chain/ci123chain/pkg/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	cmn "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	//"io"
	//"io/ioutil"
	"os"
	"path/filepath"
)
//const nodeDirPerm = 0700
var (
	//nodeDirPrefix  = "node-dir-prefix"
	nNodes   = "node-num"
	//nNonValidators = "non-validators-num"
	//outputDir      = "output-dir"
	chainPrefix		   = "chain-prefix"
)

type genState struct {
	Accounts json.RawMessage `json:"accounts"`
	Auth json.RawMessage	`json:"auth"`
	Order otypes.GenesisState 	`json:"order"`
}

// get cmd to initialize all files for tendermint testnet and application
func bootstrapGenCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "boot-gen",
		Short: "Initialize files for bootstrap",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, configs, etc.).

Note, strict routability for addresses is turned off in the configs file.

Example:
	cid boot-gen --chain-pre=xxxx --node-num=4 --output-dir=./output --validatorKey="privKey"
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			err := bootstrapGenWithConfig(config, cdc, appInit)
			return err
		},
	}
	cmd.Flags().StringVarP(&chainPrefix, "chain-pre", "c","",
		"These shards chain-prefix")
	cmd.Flags().IntP(nNodes, "N", 3,
		"Number of nodes")
	cmd.Flags().StringP(outputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().String(util.FlagValidatorKey, "", "the validator key")
	//cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
	return cmd
}

func bootstrapAddCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "boot-add",
		Short: "add a node for bootstrap",
		Long: `testnet will add a directory into {chainID} testnet and populate each with
necessary files (private validator, genesis, configs, etc.).

Note, strict routability for addresses is turned off in the configs file.

Example:
	cid boot-add --chain-pre=xxxx --output=./output
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			err := bootstrapAddNode(config, cdc, appInit)
			return err
		},
	}
	cmd.Flags().StringVarP(&chainPrefix, "chain-pre", "c","",
		"The chain-prefix")
	//cmd.Flags().IntP(nValidators, "v", 4,
	//	"Number of validators to initialize the testnet with")
	cmd.Flags().StringP(outputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	//cmd.Flags().String(nodeDirPrefix, "node",
	//	"Prefix the directory name for each node with (node results in node0, node1, ...)")
	//cmd.Flags().IntP(nNonValidators, "n", 0,
	//	"Number of non-validators to initialize the testnet with")
	//cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
	return cmd
}

// 结构为：/{output}/.{chainID}/node0/configs/genesis.json
//										   /configs.toml
//										   /node_key.json
//										   /priv_validator_key.json
//								  	/data
//							  /node1
func bootstrapGenWithConfig(c *cfg.Config, cdc *amino.Codec, appInit app.AppInit) error {
	outDir := viper.GetString(outputDir)
	nodes := viper.GetInt(nNodes)
	var validators []types.GenesisValidator
	var genFilePath string
	var validatorKey secp256k1.PrivKey
	var privStr string
	privBz := viper.GetString(util.FlagValidatorKey)
	if len(privBz) > 0 {
		privStr = fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PrivKeyName, privBz)
		err := cdc.UnmarshalJSON([]byte(privStr), &validatorKey)
		if err != nil {
			panic(err)
		}
	}else {
		validatorKey = secp256k1.GenPrivKey()
	}
	//生成chainID和rootDir
	if chainPrefix == "" {
		chainPrefix = cmn.Str(6) + "-"
	}
	rootDir := outDir

	//生成nodedir和Key和config
	for i := 0; i < nodes; i++ {
		chainName := fmt.Sprintf("%s%d", chainPrefix, i)
		c.SetRoot(filepath.Join(rootDir, chainName))
		c.Moniker = chainName
		cfg.EnsureRoot(filepath.Join(rootDir, chainName))

		pv := validator.GenFilePV(
			c.PrivValidatorKeyFile(),
			c.PrivValidatorStateFile(),
			validatorKey,
		)
		_, err := node.GenNodeKeyByPrivKey(c.NodeKeyFile(), pv.Key.PrivKey)
		if err != nil {
			return err
		}
		genTime := tmtime.Now()

		validator, appState, err := getValidator(cdc, c, appInit)
		if err != nil{
			return err
		}
		if i == 0 {
			validators = append(validators, *validator)
		}

		var appByte []byte
		appByte, _ = appState.MarshalJSON()
		var gs genState
		json.Unmarshal(appByte, &gs)
		var ls []otypes.Lists
		for j := 0; j < nodes; j++ {
			var l otypes.Lists
			l.Name = fmt.Sprintf("%s%d", chainPrefix, j)
			l.Height = 0
			ls = append(ls, l)
		}
		gs.Order.Params.OrderBook.Lists = ls
		genState, _ := cdc.MarshalJSON(gs)
		gst := json.RawMessage(string(genState))
		genFilePath = c.GenesisFile()
		err = writeGenesisFile(cdc, genFilePath, chainName, validators, gst, genTime)
		if err != nil {
				return err
			}
	}
	fmt.Printf("Successfully initialized node directories node=%v\n", viper.GetInt(nNodes))
	return nil
}

func bootstrapAddNode(c *cfg.Config, cdc *amino.Codec, appInit app.AppInit) error{
	outDir := viper.GetString(outputDir)
	if chainPrefix == ""{
		return errors.New("chainPrefix cannot be nil")
	}
	rootDir := outDir
	_, err := os.Stat(rootDir)
	if err != nil {
		return err
	}
	nodeNum, err := getNodeNum(rootDir)
	if err != nil {
		return err
	}

	id := nodeNum
	chainName := fmt.Sprintf("%s%d", chainPrefix, id)
	c.Moniker = chainName
	c.SetRoot(filepath.Join(rootDir, chainName))
	cfg.EnsureRoot(filepath.Join(rootDir, chainName))

	if err != nil{
		return err
	}
	cpDir := rootDir + "/" + chainPrefix + "0"
	nodeKeyPath := filepath.Join(cpDir, "configs/node_key.json")
	privKeyPath := filepath.Join(cpDir, "configs/priv_validator_key.json")
	privStatePath := filepath.Join(cpDir, "data/priv_validator_state.json")
	if err := CopyFile(nodeKeyPath, filepath.Join(c.RootDir, "configs/node_key.json")); err != nil {
		return err
	}
	if err := CopyFile(privKeyPath, filepath.Join(c.RootDir, "configs/priv_validator_key.json")); err != nil {
		return err
	}
	if err := CopyFile(privStatePath, filepath.Join(c.RootDir, "data/priv_validator_state.json")); err != nil {
		return err
	}

	var validators []types.GenesisValidator
	validator, appState, err := getValidator(cdc, c, appInit)
	if err != nil{
		return err
	}
	validators = append(validators, *validator)

	var appByte []byte
	appByte, _ = appState.MarshalJSON()
	var gs genState
	json.Unmarshal(appByte, &gs)
	var ls []otypes.Lists
	for j := 0; j < nodeNum - 1; j++ {
		var l otypes.Lists
		l.Name = fmt.Sprintf("%s%d", chainPrefix, j)
		l.Height = 0
		ls = append(ls, l)
	}
	genTime := time.Now()
	gs.Order.Params.OrderBook.Lists = ls
	genState, _ := cdc.MarshalJSON(gs)
	gst := json.RawMessage(string(genState))
	genFilePath := c.GenesisFile()
	err = writeGenesisFile(cdc, genFilePath, chainName, validators, gst, genTime)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully add node%d directories \n", nodeNum)
	return nil
}
