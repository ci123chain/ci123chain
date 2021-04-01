package cmd

import (
	"encoding/base64"
	///"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/config"
	"github.com/ci123chain/ci123chain/pkg/node"
	"github.com/ci123chain/ci123chain/pkg/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	cmn "github.com/tendermint/tendermint/libs/os"
	tmpro "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"net"
	"path/filepath"
	"regexp"
	"time"
	//"regexp"
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

type ValidatorAccount struct {
	Address     string    `json:"address"`
	PublicKey   string     `json:"public_key"`
	PrivateKey  string    `json:"private_key"`
}

//
//func GenTxCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "gen-tx",
//		Short: "Create genesis transfer file (under [--home]/config/gentx/gentx-[nodeID].json)",
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
//			genTxConfig := config.GenTx{
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
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Println("validator_key:", viper.GetString(FlagWithValidator))
			fmt.Println("home:", viper.GetString(tmcli.HomeFlag))
			fmt.Println("chainid:", viper.GetString(FlagChainID))

			ctxConfig := ctx.Config
			ctxConfig.SetRoot(viper.GetString(tmcli.HomeFlag))

			initConfig := InitConfig{
				ChainID: viper.GetString(FlagChainID),
				//viper.GetBool(FlagWithTxs),
				//filepath.Join(config.RootDir, "config", "gentx"),
				Overwrite: viper.GetBool(FlagOverwrite),
				//tmtime.Now(),
			}
			if initConfig.ChainID == "" {
				panic(errors.New("chain id can not be empty"))
			}
			chainID, nodeID, appMessage, pubKey, err := InitWithConfig(cdc, appInit, ctxConfig, initConfig)
			if err != nil {
				return types.ErrInitWithCfg(types.DefaultCodespace, err)
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
				return abci.ErrInternal("Marshal failed")
			}
			fmt.Println(string(out))
			return nil
		},
	}

	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	//cmd.Flags().Bool(FlagWithTxs, false, "apply existing genesis transactions from [--home]/config/gentx/")
	//cmd.Flags().AddFlagSet(appInit.FlagsAppGenState)
	//cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx) // need to add this flagset for when no GenTx's provided
	//cmd.AddCommand(GenTxCmd(ctx, cdc, appInit))
	cmd.Flags().String(FlagStateDB, "couchdb://couchdb-service:5984/ci123", "fetch new shard from db")
	cmd.Flags().String(FlagWithValidator, "", "the validator key")
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
	pkey, err := pv.GetPubKey()
	if err != nil {
		return
	}

	appGenTx, cliPrint, val, err := appInit.AppGenTx(cdc, pkey, genTxConfig)
	if err != nil {
		return
	}
	tx := app.GenesisTx{
		NodeID:    nodeID,
		IP:        genTxConfig.IP,
		Validator: val,
		AppGenTx:  appGenTx,
	}

	bz, err := types.MarshalJSONIndent(cdc, tx)
	if err != nil {
		return
	}

	/*genTxFile = json.RawMessage(bz)
	if err != nil {
		return
	}*/

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

/*
func GetChainID() (string, error){

	var id string

	statedb := viper.GetString(FlagStateDB)
	db, err := app.GetStateDB("", statedb)
	key := ortypes.ModuleName + "//" + order.OrderBookKey
	var ob order.OrderBook

	res := db.Get([]byte(key))

	err = order.ModuleCdc.UnmarshalBinaryLengthPrefixed(res, &ob)
	if err != nil {
		return "", errors.New("failed to unmarshal")
	}
	if len(ob.Actions) == 1{
		if ob.Actions[0].Type == order.OpADD {
			id = ob.Actions[0].Name
		}
	}else {
		for i := 0; i < len(ob.Actions) - 1; i++ {
			if ob.Actions[i].Type == order.OpADD {
				id = ob.Actions[i].Name
				break
			}
		}
	}
	return id, nil
}
*/

/*
func InitWithConfig(cdc *amino.Codec, appInit app.AppInit, c *cfg.Config, initConfig InitConfig)(
	chainID string, nodeID string, appMessage json.RawMessage, pubKey string, err error) {
	var validatorKey secp256k1.PrivKeySecp256k1
	var privStr string
	nodeKey, err := node.LoadNodeKey(c.NodeKeyFile())
	privBz := viper.GetString(FlagWithValidator)
	if len(privBz) > 0 {
		//1.match length
		priByt := []byte(privBz)
		length := len(priByt)
		if length != 44 {
			panic(errors.New(fmt.Sprintf("length of validator key does not match, expected %d, got %d",44 ,length)))
		}

		//2.regex match
		rule := `=$`
		reg := regexp.MustCompile(rule)
		if !reg.MatchString(privBz) {
			panic(errors.New("the end of the validator key string should be an equal sign"))
		}

		//3.match base64 encoding
		_,err := base64.StdEncoding.DecodeString(privBz)
		if err != nil {
			panic(err)
		}

		privStr = fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PrivKeyAminoName, privBz)
		err = cdc.UnmarshalJSON([]byte(privStr), &validatorKey)
		if err != nil {
			panic(err)
		}
	}else {

		panic(errors.New("validator key can not be empty"))
	}

	pv := validator.GenFilePV(
		c.PrivValidatorKeyFile(),
		c.PrivValidatorStateFile(),
		validatorKey,
	)

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

	validator := appInit.GetValidator(nodeKey.PubKey(), viper.GetString(FlagName))
	validators := []tmtypes.GenesisValidator{validator}

	appState, err := appInit.AppGenState(validators)

	if err != nil {
		return
	}
	err = writeGenesisFile(cdc, genFile, initConfig.ChainID, validators, appState, initConfig.GenesisTime)
	if err != nil {
		return
	}
	return
}

*/
func InitWithConfig(cdc *amino.Codec, appInit app.AppInit, c *cfg.Config, initConfig InitConfig)(
	chainID string, nodeID string, appMessage json.RawMessage, pubKey crypto.PubKey, err error) {
	var validatorKey ed25519.PrivKey
	var privStr string
	nodeKey, err := node.LoadNodeKey(c.NodeKeyFile())
	privBz := viper.GetString(FlagWithValidator)
	if len(privBz) > 0 {
		//1.match length
		//priByt := []byte(privBz)
		//length := len(priByt)
		//if length != 44 {
		//	panic(errors.New(fmt.Sprintf("length of validator key does not match, expected %d, got %d",44 ,length)))
		//}

		//2.regex match
		rule := `=$`
		reg := regexp.MustCompile(rule)
		if !reg.MatchString(privBz) {
			panic(errors.New("the end of the validator key string should be an equal sign"))
		}

		//3.match base64 encoding
		_,err := base64.StdEncoding.DecodeString(privBz)
		if err != nil {
			panic(err)
		}

		privStr = fmt.Sprintf(`{"type":"%s","value":"%s"}`, ed25519.PrivKeyName, privBz)
		err = cdc.UnmarshalJSON([]byte(privStr), &validatorKey)
		if err != nil {
			panic(err)
		}
	}else {
		panic(errors.New("validator key can not be empty"))
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
			Validator: tmpro.ValidatorParams{PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeSecp256k1}},
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