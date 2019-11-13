package cmd

import (
	"encoding/json"
	"errors"
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
	tmcli "github.com/tendermint/tendermint/libs/cli"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"net"
	"path/filepath"
	"time"
)

var (
	FlagName = "name"
	FlagClientHome = "home-client"
	FlagOWK = "owk"
)

var (
	FlagOverwrite = "overwrite"
	FlagWithTxs = "with-txs"
	FlagIP = "ip"
	FlagChainID = "chain-id"
)


type GenesisTx struct{
	NodeID 	string	`json:"node_id"`
	IP 		string	`json:"ip"`
	Validator tmtypes.GenesisValidator `json:"validator"`
	AppGenTx json.RawMessage  `json:"app_gen_tx"`
}

type InitConfig struct{
	ChainID 	string
	GenTxs 		bool
	GenTxsDir 	string
	Overwrite 	bool
	GenesisTime time.Time
}

func GenTxCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-tx",
		Short: "Create genesis transfer file (under [--home]/config/gentx/gentx-[nodeID].json)",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			c := ctx.Config
			c.SetRoot(viper.GetString(tmcli.HomeFlag))

			ip := viper.GetString(FlagIP)
			if len(ip) == 0 {
				eip, err := externalIP()
				if err != nil {
					return err
				}
				ip = eip
			}
			genTxConfig := config.GenTx{
				viper.GetString(FlagName),
				viper.GetString(FlagClientHome),
				viper.GetBool(FlagOWK),
				ip,
			}
			cliPrint, genTxFile, err := gentxWithConfig(cdc, appInit, c, genTxConfig)
			if err != nil {
				return err
			}
			toPrint := struct {
				AppMessage 	json.RawMessage `json:"app_message"`
				GenTxFile 	json.RawMessage `json:"gen_tx_file"`
			}{
				cliPrint,
				genTxFile,
			}
			out, err := app.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		},
	}
	cmd.Flags().String(FlagIP, "", "external facing IP to use if left blank IP will be retrieved from this machine")
	cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
	return cmd
}

func initCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use: "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(tmcli.HomeFlag))

			initConfig := InitConfig{
				viper.GetString(FlagChainID),
				viper.GetBool(FlagWithTxs),
				filepath.Join(config.RootDir, "config", "gentx"),
				viper.GetBool(FlagOverwrite),
				tmtime.Now(),
			}

			chainID, nodeID, appMessage, err := InitWithConfig(cdc, appInit, config, initConfig)
			if err != nil {
				return err
			}
			// print out some types information
			toPrint := struct {
				ChainID    string          `json:"chain_id"`
				NodeID     string          `json:"node_id"`
				AppMessage json.RawMessage `json:"app_message"`
			}{
				chainID,
				nodeID,
				appMessage,
			}
			out, err := app.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		},
	}
	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Bool(FlagWithTxs, false, "apply existing genesis transactions from [--home]/config/gentx/")
	cmd.Flags().AddFlagSet(appInit.FlagsAppGenState)
	cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx) // need to add this flagset for when no GenTx's provided
	cmd.AddCommand(GenTxCmd(ctx, cdc, appInit))
	return cmd
}

func gentxWithConfig(cdc *amino.Codec, appInit app.AppInit, config *cfg.Config, genTxConfig config.GenTx) (
	cliPrint json.RawMessage, genTxFile json.RawMessage, err error ) {

	pv := validator.GenFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
		secp256k1.GenPrivKey(),
		)
	nodeKey, err := node.GenNodeKeyByPrivKey(config.NodeKeyFile(), pv.Key.PrivKey)

	if err != nil {
		return
	}
	nodeID := string(nodeKey.ID())

	appGenTx, cliPrint, validator, err := appInit.AppGenTx(cdc, pv.GetPubKey(), genTxConfig)
	if err != nil {
		return
	}
	tx := app.GenesisTx{
		NodeID: nodeID,
		IP: 	genTxConfig.IP,
		Validator: validator,
		AppGenTx: 	appGenTx,
	}

	bz, err := app.MarshalJSONIndent(cdc, tx)
	if err != nil {
		return
	}

	genTxFile = json.RawMessage(bz)
	if err != nil {
		return
	}

	genTxFile = json.RawMessage(bz)
	name := fmt.Sprintf("gentx-%v.json", nodeID)
	writePath := filepath.Join(config.RootDir, "config", "gentx")
	file := filepath.Join(writePath, name)
	err = cmn.EnsureDir(writePath, 0700)
	if err != nil {
		return
	}
	err = cmn.WriteFile(file, bz, 0644)
	if err != nil {
		return
	}
	// Write updated config with moniker
	//config.Moniker = genTxConfig.Name
	configFilePath := filepath.Join(config.RootDir, "config", "config.toml")
	cfg.WriteConfigFile(configFilePath, config)
	return
}

func InitWithConfig(cdc *amino.Codec, appInit app.AppInit, c *cfg.Config, initConfig InitConfig)(
	chainID string, nodeID string, appMessage json.RawMessage, err error) {

	nodeKey, err := node.LoadNodeKey(c.NodeKeyFile())
	if err != nil {
		pv := validator.GenFilePV(
			c.PrivValidatorKeyFile(),
			c.PrivValidatorStateFile(),
			secp256k1.GenPrivKey(),
		)
		nodeKey, err = node.GenNodeKeyByPrivKey(c.NodeKeyFile(), pv.Key.PrivKey)
	}
	nodeID = string(nodeKey.ID())

	if initConfig.ChainID == "" {
		initConfig.ChainID = fmt.Sprintf("test-chain-%v", cmn.RandStr(6))
	}
	chainID = initConfig.ChainID

	genFile := c.GenesisFile()
	if !initConfig.Overwrite && cmn.FileExists(genFile) {
		err = fmt.Errorf("genesis.json file already exists: %v", genFile)
		return
	}

	validator := appInit.GetValidator(nodeKey.PubKey(), viper.GetString(FlagName))
	validators := []tmtypes.GenesisValidator{validator}

	appState, err := appInit.AppGenState()

	if err != nil {
		return
	}
	err = writeGenesisFile(cdc, genFile, initConfig.ChainID, validators, appState, initConfig.GenesisTime)
	if err != nil {
		return
	}
	return
}


func writeGenesisFile(cdc *amino.Codec, genesisFile string, chainID string,
	validators []tmtypes.GenesisValidator, appState json.RawMessage, genesisTime time.Time) error {

	genDoc := tmtypes.GenesisDoc{
		GenesisTime: 	genesisTime,
		ChainID: 		chainID,
		Validators: 	validators,
		AppState:		appState,
		ConsensusParams: &tmtypes.ConsensusParams{
			Block: tmtypes.DefaultBlockParams(),
			Evidence: tmtypes.DefaultEvidenceParams(),
			Validator: tmtypes.ValidatorParams{PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeSecp256k1}},
		},
	}
	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}
	return genDoc.SaveAs(genesisFile)
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if skipInterface(iface) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip := addrToIP(addr)
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func skipInterface(iface net.Interface) bool {
	if iface.Flags&net.FlagUp == 0 {
		return true // interface down
	}
	if iface.Flags&net.FlagLoopback != 0 {
		return true // loopback interface
	}
	return false
}

func addrToIP(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip
}


