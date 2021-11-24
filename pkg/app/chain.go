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
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ci123chain/ci123chain/pkg/mint"
	mint_module "github.com/ci123chain/ci123chain/pkg/mint/module"
	order_module "github.com/ci123chain/ci123chain/pkg/order/module"
	ordertypes "github.com/ci123chain/ci123chain/pkg/order/types"
	"github.com/ci123chain/ci123chain/pkg/redis"
	"github.com/ci123chain/ci123chain/pkg/registry"
	keeper2 "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	staking_module "github.com/ci123chain/ci123chain/pkg/staking/module"
	supply_module "github.com/ci123chain/ci123chain/pkg/supply/module"
	"github.com/ci123chain/ci123chain/pkg/upgrade"
	upgrade_module "github.com/ci123chain/ci123chain/pkg/upgrade/module"

	"github.com/ci123chain/ci123chain/pkg/util"
	vm_module "github.com/ci123chain/ci123chain/pkg/vm/module"
	"github.com/tendermint/tendermint/config"
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
	authtx "github.com/ci123chain/ci123chain/pkg/auth/tx"
	auth_types "github.com/ci123chain/ci123chain/pkg/auth/types"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	client "github.com/ci123chain/ci123chain/pkg/client/context"
	distr "github.com/ci123chain/ci123chain/pkg/distribution"
	ibctransfer "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer"
	ibctransferkeeper "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/keeper"
	ibctransfertypes "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	ibchost "github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	ibckeeper "github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
	"github.com/ci123chain/ci123chain/pkg/infrastructure"
	"github.com/ci123chain/ci123chain/pkg/order"
	orhandler "github.com/ci123chain/ci123chain/pkg/order/handler"
	"github.com/ci123chain/ci123chain/pkg/params"
	ptypes "github.com/ci123chain/ci123chain/pkg/params/types"
	prestaking "github.com/ci123chain/ci123chain/pkg/pre_staking"
	prestakingModule "github.com/ci123chain/ci123chain/pkg/pre_staking/module"

	"github.com/ci123chain/ci123chain/pkg/slashing"
	"github.com/ci123chain/ci123chain/pkg/staking"
	stakingTypes "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
	supply_types "github.com/ci123chain/ci123chain/pkg/supply/types"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/transfer/handler"
	"github.com/ci123chain/ci123chain/pkg/vm"
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
	flagNodeDomain = "IDG_HOST_80"
	flagShardIndex = "shardIndex"
	cacheName      = "cache"
	heightKey      = "s/k:order/OrderBook"
	StoreKey 	   = "main"
)

var (

	GasPriceConfig 	*config.GasPriceConfig

// default home directories for expected binaries
//	MainStoreKey     = sdk.NewKVStoreKey("main")
//	ContractStoreKey = sdk.NewKVStoreKey("contract")
//	TxIndexStoreKey  = sdk.NewTransientStoreKey("tx_index")

//	AccountStoreKey  = sdk.NewKVStoreKey(account.StoreKey)
//	ParamStoreKey  	 = sdk.NewKVStoreKey(params.StoreKey)
//	ParamTransStoreKey  = sdk.NewTransientStoreKey(params.TStoreKey)
//	AuthStoreKey 	 = sdk.NewKVStoreKey(auth.StoreKey)
//	SupplyStoreKey   = sdk.NewKVStoreKey(supply.StoreKey)
//	OrderStoreKey	 = sdk.NewKVStoreKey(order.StoreKey)
//	IBCStoreKey 	 = sdk.NewKVStoreKey(ibchost.StoreKey)
//
//	DisrtStoreKey    = sdk.NewKVStoreKey(k.DisrtKey)
//	StakingStoreKey  = sdk.NewKVStoreKey(staking.StoreKey)
//	preStakingStorekey = sdk.NewKVStoreKey(prestaking.StoreKey)
//	SlashingStoreKey  = sdk.NewKVStoreKey(slashing.StoreKey)
//	GravityStoreKey  = sdk.NewKVStoreKey(gravity.StoreKey)
//	WasmStoreKey     = sdk.NewKVStoreKey(vm.StoreKey)
//	MintStoreKey     = sdk.NewKVStoreKey(mint.StoreKey)
//	InfrastructureStoreKey = sdk.NewKVStoreKey(infrastructure.StoreKey)
//	IbcTransferStoreKey = sdk.NewKVStoreKey(ibctransfertypes.StoreKey)
//	CapabilityStoreKey  = sdk.NewKVStoreKey(capabilitytypes.StoreKey)

	keys = sdk.NewKVStoreKeys(StoreKey, account.StoreKey, params.StoreKey, auth.StoreKey,
		supply.StoreKey, order.StoreKey, ibchost.StoreKey, distr.StoreKey, staking.StoreKey,
		prestaking.StoreKey, slashing.StoreKey, gravity.StoreKey, vm.StoreKey, mint.StoreKey,
		infrastructure.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey, upgrade.StoreKey, registry.StoreKey,
	)
	tkeys = sdk.NewTransientStoreKeys(params.TStoreKey)
	memKeys = sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	maccPerms = map[string][]string{
		auth.FeeCollectorName: nil,
		distr.ModuleName:      nil,
		mint.ModuleName:       {supply.Minter},
		ibc.ModuleName: nil,
		gravity.ModuleName: {supply.Minter},
		ibctransfer.ModuleName: nil,
		stakingTypes.BondedPoolName: {supply.Burner, supply.Staking, supply.Minter},
		stakingTypes.NotBondedPoolName: {supply.Burner, supply.Staking},
		prestaking.ModuleName: nil,
		registry.ModuleName: nil,
	}
)


type Chain struct {
	*baseapp.BaseApp

	cdc    *amino.Codec
	appCodec codec.Marshaler
	interfaceRegistry codectypes.InterfaceRegistry

	// keys to access the substores
	//capKeyMainStore *sdk.KVStoreKey
	//txIndexStore    *sdk.TransientStoreKey
	AccountKeeper   account.AccountKeeper
	AuthKeeper 		auth.AuthKeeper
	ParamsKeepr 	params.Keeper
	StakingKeeper   keeper2.StakingKeeper
	SupplyKeeper 	supply.Keeper
	SlashingKeeper 	slashing.Keeper
	MintKeeper 		mint.Keeper
	DistrKeeper		distr.Keeper
	InfrastructureKeeper infrastructure.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	ScopedIBCKeeper  capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	PrestakingKeeper prestaking.Keeper
	RegistryKeeper  registry.Keeper
	IBCKeeper 		ibc.Keeper
	VMKeeper 		vm.Keeper
	GravityKeeper   gravity.Keeper
	UpgradeKeeper 	upgrade.Keeper
	// the module manager
	mm *module.AppManager
}

func NewChain(logger log.Logger, ldb tmdb.DB, cdb tmdb.DB, traceStore io.Writer, baseAppOptions ...func(*baseapp.BaseApp)) *Chain {
	cdc := app_types.GetCodec()
	encodingConfig := app_types.GetEncodingConfig()
	appCodec := encodingConfig.Marshaler
	interfaceRegister := encodingConfig.InterfaceRegistry
	homeDir := viper.GetString(cli.HomeFlag)
	cacheDir := os.ExpandEnv(filepath.Join(homeDir , cacheName))
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, os.ModePerm)
		os.Chmod(cacheDir, os.ModePerm)
	}
	app := baseapp.NewBaseApp("ci123", logger, ldb, cdb, cacheDir, app_types.DefaultTxDecoder(cdc), baseAppOptions...)
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
		interfaceRegistry:  interfaceRegister,
		//capKeyMainStore: 	MainStoreKey,
		//txIndexStore: 		TxIndexStoreKey,
	}
	c.UpgradeKeeper = upgrade.NewKeeper(nil, keys[upgrade.StoreKey], cdc)

	c.AccountKeeper = keeper.NewAccountKeeper(cdc, keys[account.StoreKey], acc_types.ProtoBaseAccount)

	c.ParamsKeepr = initAppParamsKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	c.SetParamStore(c.ParamsKeepr.Subspace(baseapp.Paramspace).WithKeyTable(ptypes.ConsensusParamsKeyTable()))

	c.SupplyKeeper = supply.NewKeeper(cdc, keys[supply.StoreKey], c.AccountKeeper, maccPerms)
	c.AuthKeeper = auth.NewAuthKeeper(cdc, keys[auth.StoreKey], c.GetSubspace(auth.ModuleName))
	c.StakingKeeper = staking.NewKeeper(cdc, keys[staking.StoreKey], c.AccountKeeper, c.SupplyKeeper, c.GetSubspace(staking.ModuleName), cdb)

	//prestakingKeeper := prestaking.NewKeeper(cdc, preStakingStorekey, accountKeeper, supplyKeeper, stakingKeeper, c.GetSubspace(prestaking.ModuleName),cdb)
	c.SlashingKeeper = slashing.NewKeeper(cdc, keys[slashing.StoreKey], c.StakingKeeper, c.GetSubspace(slashing.ModuleName))
	c.DistrKeeper = distr.NewKeeper(cdc, keys[distr.StoreKey], c.SupplyKeeper, c.AccountKeeper, auth.FeeCollectorName, c.GetSubspace(distr.ModuleName), c.StakingKeeper, cdb)
	c.MintKeeper = mint.NewKeeper(cdc, keys[mint.StoreKey], c.GetSubspace(mint.ModuleName), c.StakingKeeper, c.SupplyKeeper, auth.FeeCollectorName)
	c.InfrastructureKeeper = infrastructure.NewKeeper(cdc, keys[infrastructure.StoreKey])
	c.CapabilityKeeper = capabilitykeeper.NewKeeper(cdc, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	c.ScopedIBCKeeper = c.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	c.ScopedTransferKeeper = c.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	// Create IBC Keeper
	c.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], c.GetSubspace(ibchost.ModuleName), c.StakingKeeper, c.ScopedIBCKeeper,
	)



	odb := toRedisdb(cdb)
	orderKeeper := order.NewKeeper(odb, keys[order.StoreKey], c.AccountKeeper)

	c.VMKeeper = vm.NewKeeper(cdc, keys[vm.StoreKey], homeDir, c.GetSubspace(vm.ModuleName), c.AccountKeeper, c.StakingKeeper, c.UpgradeKeeper)
	c.SupplyKeeper.SetVMKeeper(c.VMKeeper)

	c.GravityKeeper = gravity.NewKeeper(cdc, keys[gravity.StoreKey], c.GetSubspace(gravity.ModuleName), c.AccountKeeper, c.StakingKeeper, c.SupplyKeeper, c.SlashingKeeper)
	c.StakingKeeper.SetHooks(staking.NewMultiStakingHooks(c.DistrKeeper.Hooks(), c.SlashingKeeper.Hooks(), c.GravityKeeper.Hooks()))

	c.PrestakingKeeper = prestaking.NewKeeper(cdc, keys[prestaking.StoreKey], c.AccountKeeper, c.SupplyKeeper, c.StakingKeeper, c.UpgradeKeeper, c.GetSubspace(prestaking.ModuleName),cdb)

	c.RegistryKeeper = registry.NewKeeper(cdc, keys[registry.StoreKey], c.SupplyKeeper, c.UpgradeKeeper)

	// Create Transfer Keepers
	ibcTransferKeeper := ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], c.GetSubspace(ibctransfertypes.ModuleName),
		c.IBCKeeper.ChannelKeeper, &c.IBCKeeper.PortKeeper,
		c.SupplyKeeper, c.ScopedTransferKeeper,
	)
	ibcTransferModule := ibc.NewTransferModule(ibcTransferKeeper)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, ibcTransferModule)
	c.IBCKeeper.SetRouter(ibcRouter)


	{ // setup module
		moduleNameOrder := []string{
			upgrade.ModuleName,
			auth_types.ModuleName,
			acc_types.ModuleName,
			supply_types.ModuleName,
			distr.ModuleName,
			order.ModuleName,
			vm.ModuleName,
			stakingTypes.ModuleName,
			prestaking.ModuleName,
			slashing.ModuleName,
			//vm.ModuleName,
			gravity.ModuleName,
			mint.ModuleName,
			infrastructure.ModuleName,
			ibc.ModuleName,
			ibctransfer.ModuleName,
		}
		// 设置modules
		c.mm = module.NewManager(
			moduleNameOrder,
			upgrade_module.AppModule{UpgradeKeeper: c.UpgradeKeeper, AccountKeeper: c.AccountKeeper, SupplyKeeper: c.SupplyKeeper},
			auth.AppModule{AuthKeeper: c.AuthKeeper},
			acc_module.AppModule{AccountKeeper: c.AccountKeeper, Cdc: cdc},
			supply_module.AppModule{Keeper: c.SupplyKeeper},
			dist_module.AppModule{DistributionKeeper: c.DistrKeeper, AccountKeeper: c.AccountKeeper, SupplyKeeper: c.SupplyKeeper},
			order_module.AppModule{OrderKeeper: &orderKeeper},
			vm_module.AppModule{Keeper: &c.VMKeeper, AccountKeeper: c.AccountKeeper},
			staking_module.AppModule{StakingKeeper: c.StakingKeeper, AccountKeeper: c.AccountKeeper, SupplyKeeper: c.SupplyKeeper},
			prestakingModule.AppModule{Keeper: c.PrestakingKeeper},
			slashing.AppModule{Keeper: c.SlashingKeeper, AccountKeeper: c.AccountKeeper, StakingKeeper: c.StakingKeeper},
			gravity.AppModule{Keeper: c.GravityKeeper, AccKeeper: c.AccountKeeper},
			//vm_module.AppModule{Keeper: &vmKeeper, AccountKeeper:accountKeeper},
			mint_module.AppModule{Keeper: c.MintKeeper},
			infrastructure_module.AppModule{Keeper: c.InfrastructureKeeper},
			ibc.NewCoreModule(&c.IBCKeeper),
			ibcTransferModule,
		)
	}
	c.mm.RegisterServices(module.NewConfigurator(c.MsgServiceRouter(),c.GRPCQueryRouter()))
	{
		// invoke router
		c.Router().AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(c.UpgradeKeeper))
		c.Router().AddRoute(transfer.RouteKey, handler.NewHandler(c.AccountKeeper))
		c.Router().AddRoute(order.RouteKey, orhandler.NewHandler(&orderKeeper))
		c.Router().AddRoute(staking.RouteKey, staking.NewHandler(c.StakingKeeper))
		c.Router().AddRoute(prestaking.RouteKey, prestaking.NewHandler(c.PrestakingKeeper))
		c.Router().AddRoute(slashing.RouteKey, slashing.NewHandler(c.SlashingKeeper))
		c.Router().AddRoute(gravity.RouteKey, gravity.NewHandler(c.GravityKeeper))
		c.Router().AddRoute(distr.RouteKey, distr.NewHandler(c.DistrKeeper))
		c.Router().AddRoute(vm.RouteKey, vm.NewHandler(&c.VMKeeper))
		c.Router().AddRoute(infrastructure.RouteKey, infrastructure.NewHandler(c.InfrastructureKeeper))
		c.Router().AddRoute(ibc.RouterKey, ibc.NewHandler(c.IBCKeeper))
		c.Router().AddRoute(ibctransfertypes.RouterKey, ibctransfer.NewHandler(ibcTransferKeeper))
		c.Router().AddRoute(slashing.RouteKey, slashing.NewHandler(c.SlashingKeeper))
	}
	{
		// query router
		c.QueryRouter().AddRoute(upgrade.RouterKey, upgrade.NewQuerier(c.UpgradeKeeper))
		c.QueryRouter().AddRoute(distr.RouteKey, distr.NewQuerier(c.DistrKeeper))
		c.QueryRouter().AddRoute(order.RouteKey, order.NewQuerier(&orderKeeper))
		c.QueryRouter().AddRoute(staking.RouteKey, staking.NewQuerier(c.StakingKeeper))
		c.QueryRouter().AddRoute(prestaking.RouteKey, prestaking.NewQuerier(c.PrestakingKeeper))
		c.QueryRouter().AddRoute(slashing.RouteKey, slashing.NewQuerier(c.SlashingKeeper, cdc))
		c.QueryRouter().AddRoute(gravity.RouteKey, gravity.NewQuerier(c.GravityKeeper))
		c.QueryRouter().AddRoute(account.RouteKey, account.NewQuerier(c.AccountKeeper))
		c.QueryRouter().AddRoute(vm.RouteKey, vm.NewQuerier(c.VMKeeper))
		c.QueryRouter().AddRoute(mint.RouteKey, mint.NewQuerier(c.MintKeeper))
		c.QueryRouter().AddRoute(infrastructure.RouteKey, infrastructure.NewQuerier(c.InfrastructureKeeper))
		c.QueryRouter().AddRoute(ibc.RouterKey, ibc.NewQuerier(c.IBCKeeper))
		c.QueryRouter().AddRoute(ibctransfertypes.RouterKey, ibc.NewTransferQuerier(ibcTransferKeeper))
	}

	c.SetAnteHandler(ante.NewAnteHandler(c.AuthKeeper, c.AccountKeeper, c.SupplyKeeper))
	c.SetDeferHandler(_defer.NewDeferHandler(c.AccountKeeper))
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
	c.CapabilityKeeper.InitializeAndSeal(ctx)

	return c
}

func setSharedIdentifier () {
	shardID := viper.GetString("ShardID")
	sdk.CommitInfoKeyFmt = shardID + "s/%d"
	sdk.LatestVersionKey = shardID + "s/latest"
}

func (c *Chain) mountStores() error {

	c.MountKVStores(keys)
	c.MountStoreMemory(memKeys)
	c.MountKVStoresTransient(tkeys)

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


func (c *Chain) ExportAppStateJSON() (ExportedApp, error) {
	ctx := c.BaseApp.NewContext(true, tmtypes.Header{Height: c.BaseApp.LastBlockHeight()})
	height := c.LastBlockHeight()
	genState :=  c.mm.ExportGenesis(ctx)
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, c.StakingKeeper)
	cp := c.BaseApp.GetConsensusParams(ctx)
	return ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: cp,
	}, err
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
	if cdb == nil {
		return nil
	}
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
	//nodeList = append(nodeList, viper.GetString(flagNodeDomain))
	n := os.Getenv(flagNodeDomain)
	if n == "" {
		n = util.GetLocalAddress()
	}
	nodeList = append(nodeList, n)
	nodeListBytes, _ = json.Marshal(nodeList)
	_ = odb.Set([]byte(redissource.FlagNodeList), nodeListBytes)

	return odb
}


func handleCache(cdb tmdb.DB, cache string, cdc *codec.Codec, app *baseapp.BaseApp) error{
	if cdb == nil {
		return nil
	}
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
	subspace, _ := app.ParamsKeepr.GetSubspace(moduleName)
	return subspace
}

func (app *Chain) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

var _ Application = (*Chain)(nil)



func NewApp(lg log.Logger, ldb tmdb.DB, cdb tmdb.DB,traceStore io.Writer) Application{
	logger.SetLogger(lg)
	return NewChain(lg, ldb, cdb, traceStore, baseapp.SetGasPriceConfig(GasPriceConfig))
}

func (app *Chain) LoadStartVersion(height int64) error {
	return app.LoadVersion(height, keys[StoreKey])
}