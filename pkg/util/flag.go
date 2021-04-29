package util


//// flag which is used.

const (
	///root cmd.

	AppName = "ci123"
	DefaultConfDir = "$HOME/.ci123"
	FlagLogLevel = "log_level"
	HomeFlag     = "home"
	LogINFO      = "state:info,x/ibc/client:info,x/ibc/connection:info,x/ibc/channel:info,*:error"
	LogDEBUG      = "*:debug"
	LogERROR     = "*:error"
	LogNONE      = "*:none"

	//common.
	//FlagHOMEDIR = "home"

	FlagETHChainID = "eth_chain_id"

	FlagCiStateDBHost = "statedb_host"

	FlagCiStateDBPort = "statedb_port"

	FlagCiStateDBTls = "statedb_tls"

	FlagCiStateDBType = "statedb_type"

	FlagCiNodeDomain = "node_domain"

	///init cmd.

	FlagOverwrite = "overwrite"

	FlagValidatorKey = "validator_key"

	FlagChainID = "chain_id"

	FlagCoinName = "denom"

	FlagName = "name"

	////start cmd.

	FlagPruning = "pruning"

	FlagWithTendermint = "with-tendermint"


	FlagAddress = "address"

	FlagTraceStore = "trace-store"

	FlagMasterDomain = "master_domain"

	FlagShardIndex = "shardIndex"

	FlagGenesis = "genesis"

	FlagNodeKey = "nodeKey"

	FlagPvs = "pvs"

	FlagPvk = "pvk"

	Version 		   = "CiChain v1.4.15"
	//// rest-server cmd.
	FlagListenAddr = "laddr"
	FlagMaxOpenConnections = "max-open"
	FlagRPCReadTimeout     = "read-timeout"
	FlagRPCWriteTimeout    = "write-timeout"
	FlagWebsocket		   = "wsport"
	GenesisFile			   = "genesis.json"
	PrivValidatorKey	   = "priv_validator_key.json"

	///lite cmd
	FlagNode = "node"
	FlagLiteHomeDir = "home-dir"

	//client.
	MinPassLength = 4

	FlagBlocked = "blocked"
	FlagHeight = "height"
	FlagHomeDir = "clihome"
	FlagVerbose = "verbose"
	FlagPassword = "password"

	FlagFile = "file"
	FlagGas = "gas"
	FlagPrivateKey = "privateKey"
	//FlagMsg = "msg"
	FlagArgs = "args"
	FlagCodeHash = "codeHash"
	FlagVersion = "version"
	FlagAuthor = "author"
	FlagEmail = "email"
	FlagDescribe = "describe"
	FlagHash = "codeHash"
	FlagFunds = "funds"
	FlagContractAddress = "contractAddress"
	FlagPrivate = "private"
	FlagSilent   = "silent"
	FlagMnemonic = "mnemonic"
	FlagHDWPath  = "hdw_path"

	//transfer.
	FlagFrom    = "from"
	FlagTo 		= "to"
	FlagAmount  = "amount"
	FlagKey		= "privKey"
	FlagIsFabric= "isFabric"

	FlagMasterPort	   = "master_port"
	DefaultMasterPort  = "80"
	FlagConfig         = "config" //config.toml
	DefaultConfigFilePath = "config.toml"
	DefaultConfigPath  = "config"
	DefaultDataPath    = "data"

	CacheName      = "cache"
	CacheSize      = "cache-size"
	HeightKey      = "s/k:order/OrderBook"
)
