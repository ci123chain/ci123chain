package app

import (
	"encoding/json"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/abci/baseapp"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	app_module "github.com/ci123chain/ci123chain/pkg/app/module"
	dist_module "github.com/ci123chain/ci123chain/pkg/distribution/module"
	mint_module "github.com/ci123chain/ci123chain/pkg/mint/module"
	order_module "github.com/ci123chain/ci123chain/pkg/order/module"
	staking_module "github.com/ci123chain/ci123chain/pkg/staking/module"
	supply_module "github.com/ci123chain/ci123chain/pkg/supply/module"
	wasm_module "github.com/ci123chain/ci123chain/pkg/wasm/module"

	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	acc_types "github.com/ci123chain/ci123chain/pkg/account/types"
	app_types "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/auth"
	"github.com/ci123chain/ci123chain/pkg/auth/ante"
	_defer "github.com/ci123chain/ci123chain/pkg/auth/defer"
	auth_types "github.com/ci123chain/ci123chain/pkg/auth/types"
	"github.com/ci123chain/ci123chain/pkg/config"
	"github.com/ci123chain/ci123chain/pkg/couchdb"
	distr "github.com/ci123chain/ci123chain/pkg/distribution"
	k "github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc"
	"github.com/ci123chain/ci123chain/pkg/mint"
	"github.com/ci123chain/ci123chain/pkg/mortgage"
	"github.com/ci123chain/ci123chain/pkg/order"
	orhandler "github.com/ci123chain/ci123chain/pkg/order/handler"
	"github.com/ci123chain/ci123chain/pkg/params"
	"github.com/ci123chain/ci123chain/pkg/staking"
	stakingTypes "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
	supply_types "github.com/ci123chain/ci123chain/pkg/supply/types"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/transfer/handler"
	"github.com/ci123chain/ci123chain/pkg/wasm"
	wasm_types "github.com/ci123chain/ci123chain/pkg/wasm/types"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"
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
	AccountStoreKey  = sdk.NewKVStoreKey(account.StoreKey)
	ParamStoreKey  	 = sdk.NewKVStoreKey(params.StoreKey)
	ParamTransStoreKey  = sdk.NewTransientStoreKey(params.TStoreKey)
	AuthStoreKey 	 = sdk.NewKVStoreKey(auth.StoreKey)
	SupplyStoreKey   = sdk.NewKVStoreKey(supply.StoreKey)
	MortgageStoreKey = sdk.NewKVStoreKey(mortgage.StoreKey)
	IBCStoreKey 	 = sdk.NewKVStoreKey(ibc.StoreKey)
	OrderStoreKey	 = sdk.NewKVStoreKey(order.StoreKey)


	disrtStoreKey    = sdk.NewKVStoreKey(k.DisrtKey)
	stakingStoreKey  = sdk.NewKVStoreKey(staking.StoreKey)
	wasmStoreKey     = sdk.NewKVStoreKey(wasm.StoreKey)
	mintStoreKey     = sdk.NewKVStoreKey(mint.StoreKey)



	maccPerms = map[string][]string{
		//mortgage.ModuleName: nil,
		auth.FeeCollectorName: nil,
		distr.ModuleName:      nil,
		mint.ModuleName:       {supply.Minter},
		ibc.ModuleName: nil,
		stakingTypes.BondedPoolName: {supply.Burner, supply.Staking},
		stakingTypes.NotBondedPoolName: {supply.Burner, supply.Staking},
	}
)


type Chain struct {
	*baseapp.BaseApp

	logger log.Logger
	cdc    *amino.Codec

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey
	contractStore   *sdk.KVStoreKey
	txIndexStore    *sdk.TransientStoreKey

	authKeeper 		auth.AuthKeeper

	// the module manager
	mm *module.AppManager
}

func NewChain(logger log.Logger, ldb tmdb.DB, cdb tmdb.DB, traceStore io.Writer) *Chain {
	cdc := app_types.MakeCodec()
	app := baseapp.NewBaseApp("ci123", logger, ldb, cdb, app_types.DefaultTxDecoder(cdc))

	c := &Chain{
		BaseApp: 			app,
		cdc: 				cdc,
		capKeyMainStore: 	MainStoreKey,
		contractStore: 		ContractStoreKey,
		txIndexStore: 		TxIndexStoreKey,
	}

	// todo mainkey?
	accKeeper := keeper.NewAccountKeeper(cdc, AccountStoreKey, acc_types.ProtoBaseAccount)

	paramsKeeper := params.NewKeeper(cdc, ParamStoreKey, ParamTransStoreKey, params.DefaultCodespace)

	supplyKeeper := supply.NewKeeper(cdc, SupplyStoreKey, accKeeper, maccPerms)

	//mortgageKeeper := mortgage.NewKeeper(MortgageStoreKey, supplyKeeper)

	authSubspace := paramsKeeper.Subspace(auth.DefaultCodespace)
	c.authKeeper = auth.NewAuthKeeper(cdc, AuthStoreKey, authSubspace)

	ibcKeeper := ibc.NewKeeper(IBCStoreKey, accKeeper, supplyKeeper)

	//fcKeeper := fc.NewFcKeeper(cdc, fcStoreKey, accKeeper)

	stakingKeeper := staking.NewKeeper(cdc, stakingStoreKey, accKeeper,supplyKeeper, paramsKeeper.Subspace(params.ModuleName), cdb)

	distrKeeper := k.NewKeeper(cdc, disrtStoreKey, supplyKeeper, accKeeper, auth.FeeCollectorName, paramsKeeper.Subspace(distr.DefaultCodespace), stakingKeeper, cdb)

	mintSubspace := paramsKeeper.Subspace(mint.DefaultCodeSpce)
	mintKeeper := mint.NewKeeper(cdc, mintStoreKey, mintSubspace, stakingKeeper, supplyKeeper, auth.FeeCollectorName)


	odb := cdb.(*couchdb.GoCouchDB)
	orderKeeper := order.NewKeeper(odb, OrderStoreKey, accKeeper)


	homeDir := viper.GetString(cli.HomeFlag)
	var wasmconfig wasm_types.WasmConfig
	wasmKeeper := wasm.NewKeeper(cdc, wasmStoreKey,homeDir, wasmconfig, accKeeper, stakingKeeper, cdb)

	stakingKeeper.SetHooks(staking.NewMultiStakingHooks(distrKeeper.Hooks()))
	module_order := []string{auth_types.ModuleName, acc_types.ModuleName, supply_types.ModuleName, distr.ModuleName, order.ModuleName,stakingTypes.ModuleName, wasm_types.ModuleName, mint.ModuleName}
	// 设置modules
	c.mm = module.NewManager(
		module_order,
		auth.AppModule{AuthKeeper: c.authKeeper},
		account.AppModule{AccountKeeper: accKeeper},
		supply_module.AppModule{Keeper: supplyKeeper},
		dist_module.AppModule{DistributionKeeper: distrKeeper, AccountKeeper:accKeeper, SupplyKeeper:supplyKeeper},
		order_module.AppModule{OrderKeeper: &orderKeeper},
		staking_module.AppModule{StakingKeeper: stakingKeeper, AccountKeeper:accKeeper, SupplyKeeper:supplyKeeper},
		wasm_module.AppModule{WasmKeeper: &wasmKeeper},
		mint_module.AppModule{Keeper: mintKeeper},
		)
	// invoke router
	c.Router().AddRoute(transfer.RouteKey, handler.NewHandler(accKeeper))
	c.Router().AddRoute(ibc.RouterKey, ibc.NewHandler(ibcKeeper))
	c.Router().AddRoute(order.RouteKey, orhandler.NewHandler(&orderKeeper))
	c.Router().AddRoute(staking.RouteKey, staking.NewHandler(stakingKeeper))
	c.Router().AddRoute(distr.RouteKey, distr.NewHandler(distrKeeper))
	c.Router().AddRoute(wasm.RouteKey, wasm.NewHandler(wasmKeeper))
	// query router
	c.QueryRouter().AddRoute(ibc.RouterKey, ibc.NewQuerier(ibcKeeper))

	c.QueryRouter().AddRoute(distr.RouteKey, distr.NewQuerier(distrKeeper))

	c.QueryRouter().AddRoute(order.RouteKey, order.NewQuerier(&orderKeeper))

	c.QueryRouter().AddRoute(staking.RouteKey, staking.NewQuerier(stakingKeeper))
	c.QueryRouter().AddRoute(account.RouteKey, account.NewQuerier(accKeeper))
	c.QueryRouter().AddRoute(wasm.RouteKey, wasm.NewQuerier(wasmKeeper))
	c.QueryRouter().AddRoute(mint.RouteKey, mint.NewQuerier(mintKeeper))


	c.SetAnteHandler(ante.NewAnteHandler(c.authKeeper, accKeeper, supplyKeeper))
	c.SetDeferHandler(_defer.NewDeferHandler(accKeeper))
	c.SetBeginBlocker(c.BeginBlocker)
	c.SetCommitter(c.Committer)
	c.SetInitChainer(c.InitChainer)
	c.SetEndBlocker(c.EndBlocker)
	shardID := viper.GetString("ShardID")
	sdk.CommitInfoKeyFmt = shardID + "s/%d"
	sdk.LatestVersionKey = shardID + "s/latest"

	err := c.mountStores()
	if err != nil {
		common.Exit(err.Error())
	}

	return c
}

func (c *Chain) mountStores() error {
	keys := []*sdk.KVStoreKey{
		c.capKeyMainStore,
		AccountStoreKey,
		c.contractStore,
		ParamStoreKey,
		AuthStoreKey,
		SupplyStoreKey,
		MortgageStoreKey,
		IBCStoreKey,
		disrtStoreKey,
		OrderStoreKey,
		stakingStoreKey,
		wasmStoreKey,
		mintStoreKey,
	}
	c.MountStoresIAVL(keys...)

	c.MountStoresTransient(c.txIndexStore, ParamTransStoreKey)

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
	//FlagsAppGenState *pflag.FlagSet
	//FlagsAppGenTx    *pflag.FlagSet

	// create the application genesis tx
	AppGenTx func(cdc *amino.Codec, pk crypto.PubKey, genTxConfig config.GenTx) (
		appGenTx, cliPrint json.RawMessage, validator types.GenesisValidator, err error)

	// AppGenState creates the core parameters initialization. It takes in a
	// pubkey meant to represent the pubkey of the validator of this machine.
	AppGenState func(validators []types.GenesisValidator) (appState json.RawMessage, err error)


	GetValidator func(pk crypto.PubKey, name string) types.GenesisValidator
}


func NewAppInit() AppInit {
	//fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)
	//fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	//fsAppGenTx.String(flagAddress, "", "address, required")
	//fsAppGenTx.String(flagClientHome, DefaultCLIHome,
	//	"home directory for the client, used for types generation")

	return AppInit{
		//FlagsAppGenState: fsAppGenState,
		//FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         CreateAppGenTx,
		AppGenState:      AppGenStateJSON,
		GetValidator:     app_module.AppGetValidator,
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



// Generate a genesis transfer with flags
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

// Generate a genesis transfer without flags
func CreateAppGenTxNF(cdc *amino.Codec, pk crypto.PubKey, addr string, gentTxConfig config.GenTx) (
	appGenTx, cliPrint json.RawMessage, validator types.GenesisValidator, err error) {

	var bz []byte
	tx := AppGenTx{
		Address: addr,
	}
	bz, err = app_types.MarshalJSONIndent(cdc, tx)
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
