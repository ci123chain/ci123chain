package cmd

import (
	"encoding/base64"
	///"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	abcitypes "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/node"
	"github.com/ci123chain/ci123chain/pkg/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	cmn "github.com/tendermint/tendermint/libs/os"
	tmpro "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"net"
	"regexp"
	"time"
	//"regexp"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
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
	FlagChainID = "chain_id"
	FlagStateDB = "statedb"
	//FlagDBName = "dbname"
	FlagWithValidator = "validator_key"
	FlagCoinName = "denom"
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
	Export 		bool
	ValidatorKey	string
}

type ValidatorAccount struct {
	Address     string    `json:"address"`
	PublicKey   string     `json:"public_key"`
	PrivateKey  string    `json:"private_key"`
}

//
//func GenTxCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "gen-tx",
//		Short: "Create genesis transfer file (under [--home]/configs/gentx/gentx-[nodeID].json)",
//		Args:  cobra.NoArgs,
//		RunE: func(_ *cobra.Command, args []string) error {
//			c := ctx.Config
//			c.SetRoot(viper.GetString(tmcli.HomeFlag))
//
//			ip := viper.GetString(FlagIP)
//			if len(ip) == 0 {
//				eip, err := externalIP()
//				if err != nil {
//					return err
//				}
//				ip = eip
//			}
//			genTxConfig := configs.GenTx{
//				viper.GetString(FlagName),
//				viper.GetString(FlagClientHome),
//				viper.GetBool(FlagOWK),
//				ip,
//			}
//			cliPrint, genTxFile, err := gentxWithConfig(cdc, appInit, c, genTxConfig)
//			if err != nil {
//				return err
//			}
//			toPrint := struct {
//				AppMessage 	json.RawMessage `json:"app_message"`
//				GenTxFile 	json.RawMessage `json:"gen_tx_file"`
//			}{
//				cliPrint,
//				genTxFile,
//			}
//			out, err := app.MarshalJSONIndent(cdc, toPrint)
//			if err != nil {
//				return err
//			}
//			fmt.Println(string(out))
//			return nil
//		},
//	}
//	cmd.Flags().String(FlagIP, "", "external facing IP to use if left blank IP will be retrieved from this machine")
//	cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
//	return cmd
//}

func initCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use: "init",
		Short: "Initialize genesis configs, priv-validator file, and p2p-node file",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			if denom := viper.GetString(FlagCoinName); denom != "" {
				abcitypes.SetCoinDenom(viper.GetString(FlagCoinName))
			}

			chainID1 := viper.GetString(FlagChainID);
			if chainID1 == "" {
				panic(errors.New("chain id can not be empty"))
			}

			exportMode := viper.GetBool(flagStartFromExport)
			validatorKey := viper.GetString(FlagWithValidator)
			if exportMode && len(validatorKey) == 0 {
				panic("validator key should provide for export mode")
			}

			ctxConfig := ctx.Config
			ctxConfig.SetRoot(viper.GetString(tmcli.HomeFlag))

			initConfig := InitConfig{
				ChainID: chainID1,
				Overwrite: viper.GetBool(FlagOverwrite),
				Export: exportMode,
				ValidatorKey: validatorKey,
			}

			chainID, nodeID, appMessage, pubKey, err := InitWithConfig(cdc, appInit, ctxConfig, initConfig)
			if err != nil {
				return sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("invalid params: %v", err.Error()))
			}

			// print out some types information
			toPrint := struct {
				ChainID    string          `json:"chain_id"`
				NodeID     string          `json:"node_id"`
				AppMessage json.RawMessage `json:"app_message"`
				PubKey     crypto.PubKey          `json:"pub_key"`
			}{
				chainID,
				nodeID,
				appMessage,
				pubKey,
			}
			out, err := types.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("marshal failed: %v", err.Error()))
			}
			fmt.Println(string(out))
			return nil
		},
	}

	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	//cmd.Flags().Bool(FlagWithTxs, false, "apply existing genesis transactions from [--home]/configs/gentx/")
	//cmd.Flags().AddFlagSet(appInit.FlagsAppGenState)
	//cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx) // need to add this flagset for when no GenTx's provided
	//cmd.AddCommand(GenTxCmd(ctx, cdc, appInit))
	cmd.Flags().String(FlagStateDB, "couchdb://couchdb-service:5984/ci123", "fetch new shard from db")
	cmd.Flags().String(FlagWithValidator, "", "the validator key")
	cmd.Flags().String(FlagCoinName, "stake",  "coin name")
	return cmd
}


func InitWithConfig(cdc *amino.Codec, appInit app.AppInit, c *cfg.Config, initConfig InitConfig)(
	chainID string, nodeID string, appMessage json.RawMessage, pubKey crypto.PubKey, err error) {
	var validatorKey ed25519.PrivKey
	var privStr string
	nodeKey, err := node.LoadNodeKey(c.NodeKeyFile())

	if len(initConfig.ValidatorKey) > 0 {
		//2.regex match
		rule := `=$`
		reg := regexp.MustCompile(rule)
		if !reg.MatchString(initConfig.ValidatorKey) {
			panic(errors.New("the end of the validator key string should be an equal sign"))
		}

		//3.match base64 encoding
		_,err := base64.StdEncoding.DecodeString(initConfig.ValidatorKey)
		if err != nil {
			panic(err)
		}

		privStr = fmt.Sprintf(`{"type":"%s","value":"%s"}`, ed25519.PrivKeyName, initConfig.ValidatorKey)
		err = cdc.UnmarshalJSON([]byte(privStr), &validatorKey)
		if err != nil {
			panic(err)
		}
	}else {
		//panic(errors.New("validator key can not be empty"))
		validatorKey = ed25519.GenPrivKey()
	}

	//create priv_validator_key.json
	pv := validator.GenFilePV(
		c.PrivValidatorKeyFile(),
		c.PrivValidatorStateFile(),
		validatorKey,
	)

	//create node_key.json
	nodeKey, err = node.GenNodeKeyByPrivKey(c.NodeKeyFile(), pv.Key.PrivKey)
	if err != nil {
		panic(err)
	}
	nodeID = string(nodeKey.ID())

	chainID = initConfig.ChainID

	genFile := c.GenesisFile()
	if !initConfig.Overwrite && cmn.FileExists(genFile) {
		err = fmt.Errorf("genesis.json file already exists: %v", genFile)
		return
	}

	val := appInit.GetValidator(nodeKey.PubKey(), viper.GetString(FlagName))
	validators := []tmtypes.GenesisValidator{val}


	pubKey = nodeKey.PubKey()//hex.EncodeToString(cdc.MustMarshalJSON(nodeKey.PubKey()))

	appState, err := appInit.AppGenState(validators)

	if err != nil {
		return
	}

	//create genesis.json
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
		ConsensusParams: &tmpro.ConsensusParams{
			Block: tmtypes.DefaultBlockParams(),
			Evidence: tmtypes.DefaultEvidenceParams(),
			Validator: tmpro.ValidatorParams{PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519}},
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