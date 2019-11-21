package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/config"
	"github.com/tanhuiya/ci123chain/pkg/node"
	"github.com/tanhuiya/ci123chain/pkg/validator"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"io"
	"os"
	"path/filepath"
)
const nodeDirPerm = 0700
var (
	nodeDirPrefix  = "node-dir-prefix"
	nValidators    = "validators-num"
	nNonValidators = "non-validators-num"
	outputDir      = "output-dir"
	chainID		   = ""
)

// get cmd to initialize all files for tendermint testnet and application
func testnetGenCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-net",
		Short: "Initialize files for a hmd testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	cid gen-net --chain-id=xxxx --validator-num=4 --non-validator=3 --output=./output
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			err := testnetGenWithConfig(config, cdc, appInit)
			return err
		},
	}
	cmd.Flags().StringVarP(&chainID, "chain-id", "c","",
		"The testnet chain-id")
	cmd.Flags().IntP(nValidators, "v", 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().StringP(outputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().String(nodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().IntP(nNonValidators, "n", 0,
		"Number of non-validators to initialize the testnet with")
	//cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
	return cmd
}

func testnetAddCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-net",
		Short: "Initialize files for a hmd testnet",
		Long: `testnet will add a directory into {chainID} testnet and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	cid add-net --chain-id=xxxx --output=./output
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			err := testnetAddNode(config, cdc, appInit)
			return err
		},
	}
	cmd.Flags().StringVarP(&chainID, "chain-id", "c","",
		"The testnet chain-id")
	//cmd.Flags().IntP(nValidators, "v", 4,
	//	"Number of validators to initialize the testnet with")
	cmd.Flags().StringP(outputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().String(nodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	//cmd.Flags().IntP(nNonValidators, "n", 0,
	//	"Number of non-validators to initialize the testnet with")
	//cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
	return cmd
}

// 结构为：/{output}/.{chainID}/node0/config/genesis.json
//										   /config.toml
//										   /node_key.json
//										   /priv_validator_key.json
//								  	/data
//							  /node1
func testnetGenWithConfig(c *cfg.Config, cdc *amino.Codec, appInit app.AppInit) error {
	outDir := viper.GetString(outputDir)
	numValidators := viper.GetInt(nValidators)
	numNonValidators := viper.GetInt(nNonValidators)

	var genesisFilePath string
	var genFile string
	var validators []types.GenesisValidator
	//生成chainID和rootDir
	if chainID == "" {
		chainID = "chain-" + cmn.RandStr(6)
	}
	rootDir := filepath.Join(outDir, "."+chainID)

	//生成nodedir和Key和config
	for i := 0; i < numValidators+numNonValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", viper.GetString(nodeDirPrefix), i)
		di := getDirsInfo(rootDir, i)
		c.SetRoot(di.NodeDir())
		c.Moniker = nodeDirName
		cfg.EnsureRoot(di.NodeDir())

		pv := validator.GenFilePV(
			c.PrivValidatorKeyFile(),
			c.PrivValidatorStateFile(),
			secp256k1.GenPrivKey(),
		)
		_, err := node.GenNodeKeyByPrivKey(c.NodeKeyFile(), pv.Key.PrivKey)
		if err != nil {
			return err
		}
	}

	genTime := tmtime.Now()
	//遍历所有的validator，获取validators，直到最后一个获取完毕，为node0生成genesis.json
	for i := 0; i < numValidators; i++ {
		di := getDirsInfo(rootDir, i)
		initConfig := InitConfig{
			chainID,
			true,
			di.GenTxsDir(),
			true,
			genTime,
		}
		c.Moniker = di.DirName()
		c.SetRoot(di.NodeDir())

		genfile, validator, appState, err := getValidator(cdc, c, appInit)
		if err != nil{
			return err
		}
		validators = append(validators, *validator)
		if i == 0{
			genFile = genfile
			genesisFilePath = c.GenesisFile()
		}
		if i == numValidators-1 {
			err := writeGenesisFile(cdc, genFile, initConfig.ChainID, validators, *appState, initConfig.GenesisTime)
			if err != nil {
				return err
			}
		}
	}
	//把genesis.json copy给其他的node
	for i := 1; i < numValidators+numNonValidators; i++ {
		id := i
		di := getDirsInfo(rootDir, id)
		c.Moniker = di.DirName()
		c.SetRoot(di.NodeDir())
		cfg.EnsureRoot(di.NodeDir())

		if err := CopyFile(genesisFilePath, filepath.Join(c.RootDir, "config/genesis.json")); err != nil {
			return err
		}
	}
	fmt.Printf("Successfully initialized node directories val=%v nval=%v\n", viper.GetInt(nValidators), viper.GetInt(nNonValidators))
	return nil
}

func testnetAddNode(c *cfg.Config, cdc *amino.Codec, appInit app.AppInit) error{
	return nil
}


type dirsInfo struct {
	rootDir string
	dirName string
}

func (di dirsInfo) DirName() string {
	return di.dirName
}

func (di dirsInfo) NodeRootDir() string {
	return filepath.Join(di.rootDir, di.dirName)
}

func (di dirsInfo) ClientDir() string {
	return filepath.Join(di.NodeRootDir(), "cicli")
}

func (di dirsInfo) NodeDir() string {
	return filepath.Join(di.NodeRootDir(), "cid")
}

func (di dirsInfo) ConfigDir() string {
	return filepath.Join(di.NodeDir(), "config")
}

func (di dirsInfo) GenTxsDir() string {
	return filepath.Join(di.rootDir, "gentxs")
}

func getDirsInfo(rootDir string, id int) dirsInfo {
	dirName := fmt.Sprintf("%s%d", viper.GetString(nodeDirPrefix), id)
	return dirsInfo{
		rootDir: rootDir,
		dirName: dirName,
	}
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func getValidator(cdc *amino.Codec, c *cfg.Config, appInit app.AppInit) (string, *types.GenesisValidator, *json.RawMessage, error){
	genTxConfig := config.GenTx{
		viper.GetString(FlagName),
		viper.GetString(FlagClientHome),
		viper.GetBool(FlagOWK),
		"127.0.0.1",
	}
	genFile := c.GenesisFile()
	// Write updated config with moniker
	// c.Moniker = genTxConfig.Name
	// config.SaveConfig(c)
	nodeKey, err := node.LoadNodeKey(c.NodeKeyFile())
	if err != nil {
		return "", nil, nil, err
	}

	_, _, validator, err := appInit.AppGenTx(cdc, nodeKey.PubKey(), genTxConfig)
	if err != nil {
		return "", nil, nil, err
	}
	//validators := []types.GenesisValidator{validator}
	//appGenTxs := []json.RawMessage{appGenTx}
	appState, err := appInit.AppGenState()
	if err != nil {
		return "", nil, nil, err
	}
	return genFile, &validator, &appState, nil
}