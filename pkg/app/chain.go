package app

import (
	"encoding/json"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/abci/baseapp"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	acc_module "github.com/ci123chain/ci123chain/pkg/account/module"
	app_module "github.com/ci123chain/ci123chain/pkg/app/module"
	capabilitykeeper "github.com/ci123chain/ci123chain/pkg/capability/keeper"
	dist_module "github.com/ci123chain/ci123chain/pkg/distribution/module"
	"github.com/ci123chain/ci123chain/pkg/gateway/redissource"
	"github.com/ci123chain/ci123chain/pkg/gravity"
	"github.com/ci123chain/ci123chain/pkg/ibc"
	porttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/port/types"
	infrastructure_module "github.com/ci123chain/ci123chain/pkg/infrastructure/module"
	mint_module "github.com/ci123chain/ci123chain/pkg/mint/module"
	order_module "github.com/ci123chain/ci123chain/pkg/order/module"
	ordertypes "github.com/ci123chain/ci123chain/pkg/order/types"
	"github.com/ci123chain/ci123chain/pkg/redis"
	staking_module "github.com/ci123chain/ci123chain/pkg/staking/module"
	supply_module "github.com/ci123chain/ci123chain/pkg/supply/module"
	vm_module "github.com/ci123chain/ci123chain/pkg/vm/module"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
	"io/ioutil"
	"path/filepath"

	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	acc_types "github.com/ci123chain/ci123chain/pkg/account/types"
	app_types "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/auth"
	"github.com/ci123chain/ci123chain/pkg/auth/ante"
	_defer "github.com/ci123chain/ci123chain/pkg/auth/defer"
	auth_types "github.com/ci123chain/ci123chain/pkg/auth/types"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	distr "github.com/ci123chain/ci123chain/pkg/distribution"
	k "github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	ibctransfer "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer"
	ibctransferkeeper "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/keeper"
	ibctransfertypes "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	ibchost "github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	ibckeeper "github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
	"github.com/ci123chain/ci123chain/pkg/infrastructure"
	"github.com/ci123chain/ci123chain/pkg/mint"
	"github.com/ci123chain/ci123chain/pkg/order"
	orhandler "github.com/ci123chain/ci123chain/pkg/order/handler"
	"github.com/ci123chain/ci123chain/pkg/params"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"github.com/ci123chain/ci123chain/pkg/slashing"
	stakingTypes "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
	supply_types "github.com/ci123chain/ci123chain/pkg/supply/types"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/transfer/handler"
	"github.com/ci123chain/ci123chain/pkg/vm"
	wasm_types "github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	cmn "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
	"io"
	"os"
)

const (
	flagAddress    = "address"
	flagName       = "name"
	flagClientHome = "home-client"
	flagNodeDomain = "node_domain"
	flagShardIndex = "shardIndex"
	cacheName      = "cache"
	heightKey      = "s/k:order/OrderBook"
)

var (
	// default home directories for expected binaries
	MainStoreKey     = sdk.NewKVStoreKey("main")
	ContractStoreKey = sdk.NewKVStoreKey("contract")
	TxIndexStoreKey  = sdk.NewTransientStoreKey("tx_index")
	AccountStoreKey  = sdk.NewKVStoreKey(account.StoreKey)
	ParamStoreKey  	 = sdk.NewKVStoreKey(params.StoreKey)
	ParamTransStoreKey  = sdk.NewTransientStoreKey(params.TStoreKey)
	AuthStoreKey 	 = sdk.NewKVStoreKey(auth.StoreKey)
	SupplyStoreKey   = sdk.NewKVStoreKey(supply.StoreKey)
	OrderStoreKey	 = sdk.NewKVStoreKey(order.StoreKey)
	IBCStoreKey 	 = sdk.NewKVStoreKey(ibchost.StoreKey)

	DisrtStoreKey    = sdk.NewKVStoreKey(k.DisrtKey)
	StakingStoreKey  = sdk.NewKVStoreKey(staking.StoreKey)
	SlashingStoreKey  = sdk.NewKVStoreKey(slashing.StoreKey)
	GravityStoreKey  = sdk.NewKVStoreKey(gravity.StoreKey)
	WasmStoreKey     = sdk.NewKVStoreKey(vm.StoreKey)
	MintStoreKey     = sdk.NewKVStoreKey(mint.StoreKey)
	InfrastructureStoreKey = sdk.NewKVStoreKey(infrastructure.StoreKey)
	IbcTransferStoreKey = sdk.NewKVStoreKey(ibctransfertypes.StoreKey)
	CapabilityStoreKey  = sdk.NewKVStoreKey(capabilitytypes.StoreKey)

	memKeys = sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	maccPerms = map[string][]string{
		auth.FeeCollectorName: nil,
		distr.ModuleName:      nil,
		mint.ModuleName:       {supply.Minter},
		ibc.ModuleName: nil,
		ibctransfer.ModuleName: nil,
		stakingTypes.BondedPoolName: {supply.Burner, supply.Staking},
		stakingTypes.NotBondedPoolName: {supply.Burner, supply.Staking},
	}
)


type Chain struct {
	*baseapp.BaseApp

	logger log.Logger
	cdc    *amino.Codec
	appCodec codec.Marshaler
	interfaceRegistery codectypes.InterfaceRegistry

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey
	contractStore   *sdk.KVStoreKey
	txIndexStore    *sdk.TransientStoreKey

	authKeeper 		auth.AuthKeeper
	paramsKeepr 	params.Keeper
	// the module manager
	mm *module.AppManager
}

func NewChain(logger log.Logger, ldb tmdb.DB, cdb tmdb.DB, traceStore io.Writer) *Chain {
	cdc := app_types.GetCodec()
	encodingConfig := app_types.GetEncodingConfig()
	appCodec := encodingConfig.Marshaler
	interfaceRegister := encodingConfig.InterfaceRegistry
	cacheDir := os.ExpandEnv(filepath.Join(viper.GetString(cli.HomeFlag) , cacheName))
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, os.ModePerm)
		os.Chmod(cacheDir, os.ModePerm)
	}
	app := baseapp.NewBaseApp("ci123", logger, ldb, cdb, cacheDir, app_types.DefaultTxDecoder(cdc))
	cache := filepath.Join(cacheDir, cacheName)
	if _, err := os.Stat(cache); !os.IsNotExist(err) {
		//cache exist, check latest version
		err := handleCache(cdb, cache, cdc, app)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	c := &Chain{
		BaseApp: 			app,
		cdc: 				cdc,
		appCodec:           appCodec,
		interfaceRegistery: interfaceRegister,
		capKeyMainStore: 	MainStoreKey,
		contractStore: 		ContractStoreKey,
		txIndexStore: 		TxIndexStoreKey,
	}

	accountKeeper := keeper.NewAccountKeeper(cdc, AccountStoreKey, acc_types.ProtoBaseAccount)

	c.paramsKeepr = initAppParamsKeeper(cdc, ParamStoreKey, ParamTransStoreKey)

	supplyKeeper := supply.NewKeeper(cdc, SupplyStoreKey, accountKeeper, maccPerms)

	c.authKeeper = auth.NewAuthKeeper(cdc, AuthStoreKey, c.GetSubspace(auth.ModuleName))

	stakingKeeper := staking.NewKeeper(cdc, StakingStoreKey, accountKeeper, supplyKeeper, c.GetSubspace(staking.ModuleName), cdb)

	slashingKeeper := slashing.NewKeeper(cdc, SlashingStoreKey, stakingKeeper, c.GetSubspace(slashing.ModuleName))

	gravityKeeper := gravity.NewKeeper(cdc, GravityStoreKey, c.GetSubspace(gravity.ModuleName), stakingKeeper, supplyKeeper, slashingKeeper)

	distrKeeper := k.NewKeeper(cdc, DisrtStoreKey, supplyKeeper, accountKeeper, auth.FeeCollectorName, c.GetSubspace(distr.ModuleName), stakingKeeper, cdb)

	mintKeeper := mint.NewKeeper(cdc, MintStoreKey, c.GetSubspace(mint.ModuleName), stakingKeeper, supplyKeeper, auth.FeeCollectorName)

	infrastructureKeeper := infrastructure.NewKeeper(cdc, InfrastructureStoreKey)
	stakingKeeper.SetHooks(staking.NewMultiStakingHooks(distrKeeper.Hooks(), slashingKeeper.Hooks(), gravityKeeper.Hooks()))
	

	capabilityKeeper := capabilitykeeper.NewKeeper(cdc, CapabilityStoreKey, memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	// Create IBC Keeper
	IBCKeeper := ibckeeper.NewKeeper(
		appCodec, IBCStoreKey, c.GetSubspace(ibchost.ModuleName), stakingKeeper, scopedIBCKeeper,
	)

	// Create Transfer Keepers
	ibcTransferKeeper := ibctransferkeeper.NewKeeper(
		appCodec, IbcTransferStoreKey, c.GetSubspace(ibctransfertypes.ModuleName),
		IBCKeeper.ChannelKeeper, &IBCKeeper.PortKeeper,
		supplyKeeper, scopedTransferKeeper,
	)

	odb := toRedisdb(cdb)
	orderKeeper := order.NewKeeper(odb, OrderStoreKey, accountKeeper)

	homeDir := viper.GetString(cli.HomeFlag)
	var wasmconfig wasm_types.WasmConfig
	vmKeeper := vm.NewKeeper(cdc, WasmStoreKey, homeDir, wasmconfig, c.GetSubspace(vm.ModuleName), accountKeeper, stakingKeeper, cdb)

	ibcTransferModule := ibc.NewTransferModule(ibcTransferKeeper)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, ibcTransferModule)
	IBCKeeper.SetRouter(ibcRouter)

	{ // setup module
		moduleNameOrder := []string{
			auth_types.ModuleName,
			acc_types.ModuleName,
			supply_types.ModuleName,
			distr.ModuleName,
			order.ModuleName,
			stakingTypes.ModuleName,
			slashing.ModuleName,
			gravity.ModuleName,
			vm.ModuleName,
			mint.ModuleName,
			infrastructure.ModuleName,
			ibc.ModuleName,
			ibctransfer.ModuleName,
		}
		// 设置modules
		c.mm = module.NewManager(
			moduleNameOrder,
			auth.AppModule{AuthKeeper: c.authKeeper},
			acc_module.AppModule{AccountKeeper: accountKeeper},
			supply_module.AppModule{Keeper: supplyKeeper},
			dist_module.AppModule{DistributionKeeper: distrKeeper, AccountKeeper: accountKeeper, SupplyKeeper: supplyKeeper},
			order_module.AppModule{OrderKeeper: &orderKeeper},
			staking_module.AppModule{StakingKeeper: stakingKeeper, AccountKeeper: accountKeeper, SupplyKeeper: supplyKeeper},
			slashing.AppModule{Keeper: slashingKeeper, AccountKeeper: accountKeeper, StakingKeeper: stakingKeeper},
			gravity.AppModule{Keeper: gravityKeeper, AccKeeper: accountKeeper},
			vm_module.AppModule{Keeper: &vmKeeper},
			mint_module.AppModule{Keeper: mintKeeper},
			infrastructure_module.AppModule{Keeper: infrastructureKeeper},
			ibc.NewCoreModule(&IBCKeeper),
			ibcTransferModule,
		)
	}
	c.mm.RegisterServices(module.NewConfigurator(c.MsgServiceRouter(),c.GRPCQueryRouter()))
	{
		// invoke router
		c.Router().AddRoute(transfer.RouteKey, handler.NewHandler(accountKeeper))
		c.Router().AddRoute(order.RouteKey, orhandler.NewHandler(&orderKeeper))
		c.Router().AddRoute(staking.RouteKey, staking.NewHandler(stakingKeeper))
		c.Router().AddRoute(slashing.RouteKey, slashing.NewHandler(slashingKeeper))
		c.Router().AddRoute(gravity.RouteKey, gravity.NewHandler(gravityKeeper))
		c.Router().AddRoute(distr.RouteKey, distr.NewHandler(distrKeeper))
		c.Router().AddRoute(vm.RouteKey, vm.NewHandler(&vmKeeper))
		c.Router().AddRoute(infrastructure.RouteKey, infrastructure.NewHandler(infrastructureKeeper))
		c.Router().AddRoute(ibc.RouterKey, ibc.NewHandler(IBCKeeper))
		c.Router().AddRoute(ibctransfertypes.RouterKey, ibctransfer.NewHandler(ibcTransferKeeper))
	}
	{
		// query router
		c.QueryRouter().AddRoute(distr.RouteKey, distr.NewQuerier(distrKeeper))
		c.QueryRouter().AddRoute(order.RouteKey, order.NewQuerier(&orderKeeper))
		c.QueryRouter().AddRoute(staking.RouteKey, staking.NewQuerier(stakingKeeper))
		c.QueryRouter().AddRoute(slashing.RouteKey, slashing.NewQuerier(slashingKeeper, cdc))
		c.QueryRouter().AddRoute(gravity.RouteKey, gravity.NewQuerier(gravityKeeper))
		c.QueryRouter().AddRoute(account.RouteKey, account.NewQuerier(accountKeeper))
		c.QueryRouter().AddRoute(vm.RouteKey, vm.NewQuerier(vmKeeper))
		c.QueryRouter().AddRoute(mint.RouteKey, mint.NewQuerier(mintKeeper))
		c.QueryRouter().AddRoute(infrastructure.RouteKey, infrastructure.NewQuerier(infrastructureKeeper))
		c.QueryRouter().AddRoute(ibc.RouterKey, ibc.NewQuerier(IBCKeeper))
		c.QueryRouter().AddRoute(ibctransfertypes.RouterKey, ibc.NewTransferQuerier(ibcTransferKeeper))
	}

	c.SetAnteHandler(ante.NewAnteHandler(c.authKeeper, accountKeeper, supplyKeeper))
	c.SetDeferHandler(_defer.NewDeferHandler(accountKeeper))
	c.SetBeginBlocker(c.BeginBlocker)
	c.SetInitChainer(c.InitChainer)
	c.SetEndBlocker(c.EndBlocker)

	// 设置分片前缀标示
	setSharedIdentifier()

	err := c.mountStores()
	if err != nil {
		cmn.Exit(err.Error())
	}
	ctx := c.BaseApp.NewUncachedContext(true, tmtypes.Header{})
	capabilityKeeper.InitializeAndSeal(ctx)

	return c
}

func setSharedIdentifier () {
	shardID := viper.GetString("ShardID")
	sdk.CommitInfoKeyFmt = shardID + "s/%d"
	sdk.LatestVersionKey = shardID + "s/latest"
}

func (c *Chain) mountStores() error {
	keys := []*sdk.KVStoreKey{
		c.capKeyMainStore,
		AccountStoreKey,
		c.contractStore,
		ParamStoreKey,
		AuthStoreKey,
		SupplyStoreKey,
		IBCStoreKey,
		DisrtStoreKey,
		OrderStoreKey,
		StakingStoreKey,
		SlashingStoreKey,
		GravityStoreKey,
		WasmStoreKey,
		MintStoreKey,
		InfrastructureStoreKey,
		IbcTransferStoreKey,
		CapabilityStoreKey,
	}
	c.MountStoresIAVL(keys...)
	c.MountStoreMemory(memKeys)
	c.MountStoresTransient(c.txIndexStore, ParamTransStoreKey)

	for _, key := range keys {
		if err := c.LoadLatestVersion(key); err != nil {
			return err
		}
	}

	return nil
}


// initParamsKeeper init params keeper and its subspaces
func initAppParamsKeeper(cdc *codec.Codec, key *sdk.KVStoreKey, tkey *sdk.TransientStoreKey) params.Keeper {
	paramsKeeper := params.NewKeeper(cdc, key, tkey, params.DefaultCodespace)

	paramsKeeper.Subspace(account.ModuleName)
	paramsKeeper.Subspace(auth.ModuleName)
	paramsKeeper.Subspace(stakingTypes.ModuleName)
	paramsKeeper.Subspace(slashing.ModuleName)
	paramsKeeper.Subspace(gravity.ModuleName)
	paramsKeeper.Subspace(mint.ModuleName)
	paramsKeeper.Subspace(distr.ModuleName)
	paramsKeeper.Subspace(vm.ModuleName)
	paramsKeeper.Subspace(ibc.ModuleName)
	paramsKeeper.Subspace(ibctransfer.ModuleName)
	//paramsKeeper.Subspace(slashingtypes.ModuleName)
	//paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	//paramsKeeper.Subspace(crisistypes.ModuleName)

	return paramsKeeper
}


func (c *Chain) ExportAppStateJSON() (json.RawMessage, []types.GenesisValidator, error) {
	// TODO: Implement
	// Currently non-functional, just enough to compile
	return nil, nil, errors.New("not implemented error")
}

//_____________________________________________________________________

// Core functionality passed from the application to the server init command
type AppInit struct {
	// AppGenState creates the collactor parameters initialization. It takes in a
	// pubkey meant to represent the pubkey of the validator of this machine.
	AppGenState func(validators []types.GenesisValidator) (appState json.RawMessage, err error)

	GetValidator func(pk crypto.PubKey, name string) types.GenesisValidator
}


func NewAppInit() AppInit {

	return AppInit{
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

func toRedisdb(cdb tmdb.DB) *redis.RedisDB {
	//odb := cdb.(*couchdb.GoCouchDB)
	var nodeList []string
	odb := cdb.(*redis.RedisDB)

	nodeListBytes, err := odb.Get([]byte(redissource.FlagNodeList))
	if err != nil {
		cmn.Exit(err.Error())
	}
	if nodeListBytes != nil {
		err := json.Unmarshal(nodeListBytes, &nodeList)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}
	nodeList = append(nodeList, viper.GetString(flagNodeDomain))
	nodeListBytes, _ = json.Marshal(nodeList)
	odb.Set([]byte(redissource.FlagNodeList), nodeListBytes)

	return odb
}


func handleCache(cdb tmdb.DB, cache string, cdc *codec.Codec, app *baseapp.BaseApp) error{
	var cacheMap []store.CacheMap
	cacheFile, _ := ioutil.ReadFile(cache)
	err := json.Unmarshal(cacheFile, &cacheMap)
	if err != nil {
		return errors.New("cache file unmarshal to cacheMap failed")
	}
	batch := cdb.NewBatch()
	defer batch.Close()

	var orderBook ordertypes.OrderBook
	for _, v := range cacheMap {
		if string(v.Key) == heightKey {
			cdc.UnmarshalJSON(v.Value, &orderBook)
			break
		}
	}

	if orderBook.Lists[orderBook.Current.Index].Height == app.GetLatestVersion() + 1 {
		os.Remove(cache)
	} else {
		for _, v := range cacheMap {
			batch.Set(v.Key, v.Value)
		}
		batch.Write()
		//remove cache
		os.Remove(cache)
	}
	return nil
}

// NOTE: This is solely to be used for testing purposes.
func (app *Chain) GetSubspace(moduleName string) params.Subspace {
	subspace, _ := app.paramsKeepr.GetSubspace(moduleName)
	return subspace
}

