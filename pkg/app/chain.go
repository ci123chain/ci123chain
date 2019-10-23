package app

import (
	"CI123Chain/pkg/abci/baseapp"
	"CI123Chain/pkg/account"
	"CI123Chain/pkg/config"
	"CI123Chain/pkg/db"
	"CI123Chain/pkg/handler"
	"CI123Chain/pkg/transaction"
	"encoding/json"
	"errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
	tmtypes "github.com/tendermint/tendermint/types"
	"io"
	"os"
	sdk "CI123Chain/pkg/abci/types"
)

const (
	flagAddress    = "address"
	flagName       = "name"
	flagClientHome = "home-client"
)

var (
	// default home directories for expected binaries
	DefaultCLIHome  = os.ExpandEnv("$HOME/.cicli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.ci123d")

	MainStoreKey     = sdk.NewKVStoreKey("main")
	ContractStoreKey = sdk.NewKVStoreKey("contract")
	TxIndexStoreKey  = sdk.NewTransientStoreKey("tx_index")
)


type Chain struct {
	*baseapp.BaseApp

	logger log.Logger
	cdc    *amino.Codec

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey
	contractStore   *sdk.KVStoreKey
	txIndexStore    *sdk.TransientStoreKey
}

func NewChain(logger log.Logger, tmdb tmdb.DB, traceStore io.Writer) *Chain {
	app := baseapp.NewBaseApp("ci123", logger, tmdb, transaction.DecodeTx)
	c := &Chain{
		BaseApp: 			app,
		cdc: 				cdc,
		capKeyMainStore: 	MainStoreKey,
		contractStore: 		ContractStoreKey,
		txIndexStore: 		TxIndexStoreKey,
	}
	am := account.NewAccountMapper(c.capKeyMainStore)

	txm := transaction.NewTxIndexMapper(c.txIndexStore)
	sm := db.NewStateManager(c.contractStore)

	c.SetHandler(handler.NewHandler(txm, am, sm))
	c.SetAnteHandler(handler.NewAnteHandler(am))
	c.SetInitChainer(GetInitChainer(am))

	err := c.mountStores()
	if err != nil {
		common.Exit(err.Error())
	}

	return c
}

func (c *Chain) mountStores() error {
	keys := []*sdk.KVStoreKey{
		c.capKeyMainStore, c.contractStore,
	}

	c.MountStoresIAVL(keys...)
	c.MountStoresTransient(c.txIndexStore)

	for _, key := range keys {
		if err := c.LoadLatestVersion(key); err != nil {
			return err
		}
	}

	return nil
}

func (c *Chain) ExportAppStateJSON() (json.RawMessage, []types.GenesisValidator, error) {
	// TODO: Implement
	// Currently non-functional, just enough to compile
	return nil, nil, errors.New("not implemented error")
}

//_____________________________________________________________________

// Core functionality passed from the application to the server init command
type AppInit struct {

	// flags required for application init functions
	FlagsAppGenState *pflag.FlagSet
	FlagsAppGenTx    *pflag.FlagSet

	// create the application genesis tx
	AppGenTx func(cdc *amino.Codec, pk crypto.PubKey, genTxConfig config.GenTx) (
		appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error)

	// AppGenState creates the core parameters initialization. It takes in a
	// pubkey meant to represent the pubkey of the validator of this machine.
	AppGenState func(cdc *amino.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error)
}


func NewAppInit() AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagAddress, "", "address, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")

	return AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         CreateAppGenTx,
		AppGenState:      AppGenStateJSON,
	}
}


// simple genesis tx
type GenesisTx struct {
	NodeID    string                 `json:"node_id"`
	IP        string                 `json:"ip"`
	Validator types.GenesisValidator `json:"validator"`
	AppGenTx  json.RawMessage        `json:"app_gen_tx"`
}

type AppGenTx struct {
	// currently takes address as string because unmarshaling Ether address fails
	Address string `json:"address"`
}

// Generate a genesis transaction with flags
// pk: publickey of validator
func CreateAppGenTx(cdc *amino.Codec, pk crypto.PubKey, gentTxConfig config.GenTx) (
	appGenTx, cliPrint json.RawMessage, validator types.GenesisValidator, err error) {
	addrString := viper.GetString(flagAddress)

	bz, err := cdc.MarshalJSON("success")
	if err != nil {
		panic(err)
	}
	cliPrint = json.RawMessage(bz)
	appGenTx, _, validator, err = CreateAppGenTxNF(cdc, pk, addrString, gentTxConfig)
	return
}

// Generate a genesis transaction without flags
func CreateAppGenTxNF(cdc *amino.Codec, pk crypto.PubKey, addr string, gentTxConfig config.GenTx) (
	appGenTx, cliPrint json.RawMessage, validator types.GenesisValidator, err error) {

	var bz []byte
	tx := AppGenTx{
		Address: addr,
	}
	bz, err = MarshalJSONIndent(cdc, tx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)
	validator = types.GenesisValidator{
		PubKey: pk,
		Power:  1,
		Name:   gentTxConfig.Name,
	}
	return
}
