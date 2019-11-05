package app

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/baseapp"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	"github.com/tanhuiya/ci123chain/pkg/account/keeper"
	acc_types "github.com/tanhuiya/ci123chain/pkg/account/types"
	"github.com/tanhuiya/ci123chain/pkg/config"
	"github.com/tanhuiya/ci123chain/pkg/db"
	"github.com/tanhuiya/ci123chain/pkg/handler"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"encoding/json"
	"errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
	"io"
	"os"
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

	accModule 		*account.AppModule
}

func NewChain(logger log.Logger, tmdb tmdb.DB, traceStore io.Writer) *Chain {
	app := baseapp.NewBaseApp("ci123", logger, tmdb, transaction.DecodeTx)

	cdc := MakeCodec()

	c := &Chain{
		BaseApp: 			app,
		cdc: 				cdc,
		capKeyMainStore: 	MainStoreKey,
		contractStore: 		ContractStoreKey,
		txIndexStore: 		TxIndexStoreKey,
	}

	txm := transaction.NewTxIndexMapper(c.txIndexStore)
	sm := db.NewStateManager(c.contractStore)

	// 设置module
	// todo mainkey?
	accKeeper := keeper.NewAccountKeeper(cdc, c.capKeyMainStore, acc_types.ProtoBaseAccount)
	c.accModule = &account.AppModule{ accKeeper}


	c.SetHandler(handler.NewHandler(txm, accKeeper, sm))
	c.SetAnteHandler(handler.NewAnteHandler(accKeeper))
	c.SetInitChainer(GetInitChainer(*c.accModule))

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


	GetValidator func(pk crypto.PubKey, name string) types.GenesisValidator
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
		GetValidator:     AppGetValidator,
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

func AppGetValidator(pk crypto.PubKey, name string) types.GenesisValidator {
	validator := types.GenesisValidator{
		PubKey: pk,
		Power:  1,
		Name:   name,
	}
	return validator
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
