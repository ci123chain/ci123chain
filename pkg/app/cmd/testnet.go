package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"
	"path/filepath"
)

var (
	nodeDirPrefix  = "node-dir-prefix"
	nValidators    = "validators-num"
	nNonValidators = "non-validators-num"
	outputDir      = "output-dir"

	startingIPAddress = "starting-ip-address"
)

// get cmd to initialize all files for tendermint testnet and application
func testnetFilesCmd(ctx *app.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a hmd testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:

	cid testnet --validators-num 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			err := testnetWithConfig(config, cdc)
			return err
		},
	}
	cmd.Flags().IntP(nValidators, "v", 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().StringP(outputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().String(nodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().IntP(nNonValidators, "n", 0,
		"Number of non-validators to initialize the testnet with")
	cmd.Flags().String(startingIPAddress, "192.168.0.1",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	//cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
	return cmd
}

func testnetWithConfig(c *cfg.Config, cdc *amino.Codec) error {
	outDir := viper.GetString(outputDir)
	numValidators := viper.GetInt(nValidators)
	//numNonValidators := viper.GetInt(nNonValidators)

	// Generate genesis.json and config.toml
	chainID := "chain-" + cmn.RandStr(6)
	genTime := tmtime.Now()
	//var genesisFilePath string
	for i := 0; i < numValidators; i++ {
		di := getDirsInfo(outDir, i)
		c.Moniker = di.DirName()
		c.SetRoot(di.NodeDir())

		initConfig := InitConfig{
			chainID,
			true,
			di.GenTxsDir(),
			true,
			genTime,
		}
		// Run `init` and generate genesis.json and config.toml
		_, _, _, err := InitWithConfig(cdc, app.AppInit{}, c, initConfig)
		if err != nil {
			return err
		}
		//if i == 0 {
		//	genesisFilePath = c.GenesisFile()
		//}
	}
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
	return filepath.Join(di.NodeRootDir(), "hmcli")
}

func (di dirsInfo) NodeDir() string {
	return filepath.Join(di.NodeRootDir(), "hmd")
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