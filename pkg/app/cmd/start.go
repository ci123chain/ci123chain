package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app"
	hnode "github.com/ci123chain/ci123chain/pkg/node"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abcis "github.com/tendermint/tendermint/abci/server"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/node"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	flagWithTendermint = "with-tendermint"
	flagAddress        = "address"
	flagTraceStore     = "trace-store"
	flagPruning        = "pruning"
	//flagLogLevel       = "log-level"
	flagStateDB 	   = "statedb"
	flagCiStateDBType  = "statedb_type"
	flagCiStateDBHost  = "statedb_host"
	flagCiStateDBTls   = "statedb_tls"
	flagCiStateDBPort  = "statedb_port"
	flagCiNodeDomain   = "node_domain"
	flagMasterDomain   = "master_domain"
	flagShardIndex     = "shardIndex"
	flagGenesis        = "genesis" //genesis.json
	flagNodeKey        = "nodeKey" //node_key.json
	flagPvs            = "pvs" //priv_validator_state.json
	flagPvk            = "pvk" //priv_validator_key.json
	version 		   = "CiChain testTM6"
)

func startCmd(ctx *app.Context, appCreator app.AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use: "start",
		Short: "Run the full node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !viper.GetBool(flagWithTendermint) {
				ctx.Logger.Info("Starting ABCI Without Tendermint")
				return startStandAlone(ctx, appCreator)
			}
			ctx.Logger.Info(version)
			ctx.Logger.Info("Starting ABCI with Tendermint")
			if len(viper.GetString(flagMasterDomain)) == 0 && len(viper.GetString(flagGenesis)) != 0 {
				preSetConfig(ctx)
			}
			_, err := StartInProcess(ctx, appCreator)
			if err != nil {
				return err
			}
			select {}
		},
	}

	cmd.Flags().Bool(flagWithTendermint, true, "Run abci app embedded in-process with tendermint")
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(flagPruning, "syncable", "Pruning strategy: syncable, nothing, everything")
	cmd.Flags().String(flagCiStateDBType, "redis", "database type")
	cmd.Flags().String(flagCiStateDBHost, "", "db host")
	cmd.Flags().Uint64(flagCiStateDBPort, 7443, "db port")
	cmd.Flags().Bool(flagCiStateDBTls, true, "use tls")
	cmd.Flags().String(flagCiNodeDomain, "", "node domain")
	cmd.Flags().String(flagShardIndex, "", "index of shard")
	cmd.Flags().String(flagMasterDomain, "", "master node")

	//cmd.Flags().String(flagLogLevel, "debug", "Run abci app with different log level")
	tcmd.AddNodeFlags(cmd)
	return cmd
}

func startStandAlone(ctx *app.Context, appCreator app.AppCreator) error {
	addr := viper.GetString(flagAddress)
	home := viper.GetString("home")
	traceStore := viper.GetString(flagTraceStore)
	stateDB := viper.GetString(flagStateDB)

	app, err := appCreator(home, ctx.Logger, stateDB, traceStore)
	if err != nil {
		return err
	}
	svr, err := abcis.NewServer(addr, "socket", app)
	if err != nil {
		return errors.Errorf("error creating listener: %v\n", err)
	}
	svr.SetLogger(ctx.Logger.With("module", "abci-server"))

	err = svr.Start()
	if err != nil {
		tos.Exit(err.Error())
	}

	tos.TrapSignal(ctx.Logger, func() {
		err = svr.Stop()
		if err != nil {
			tos.Exit(err.Error())
		}
	})
	return nil
}

func StartInProcess(ctx *app.Context, appCreator app.AppCreator) (*node.Node, error) {
	cfg := ctx.Config
	home := cfg.RootDir
	traceStore := viper.GetString(flagTraceStore)
	stateDB := ""

	dbType := viper.GetString(flagCiStateDBType)
	if dbType == "" {
		dbType = "redis"
	}
	dbHost := viper.GetString(flagCiStateDBHost)
	if dbHost == "" {
		return nil, errors.New(fmt.Sprintf("%s can not be empty", flagCiStateDBHost))
	}
	dbTls := viper.GetBool(flagCiStateDBTls)
	dbPort := viper.GetUint64(flagCiStateDBPort)
	p := strconv.FormatUint(dbPort, 10)

	switch dbType {
	case "redis":
		stateDB = "redisdb://" + dbHost + ":" + p
		if dbTls {
			stateDB += "#tls"
		}
	default:
		return nil, errors.New(fmt.Sprintf("type of db: %s, which is not reids not implement yet", dbType))
	}

	nodeDomain := viper.GetString(flagCiNodeDomain)

	if nodeDomain == "" {
		return nil, errors.New("node domain can not be empty")
	}

	gendoc, err := types.GenesisDocFromFile(cfg.GenesisFile())
	if err != nil {
		panic(err)
	}
	viper.Set("ShardID", gendoc.ChainID)

	app, err := appCreator(home, ctx.Logger, stateDB, traceStore)
	if err != nil {
		return nil, err
	}

	nodeKey, err := hnode.LoadNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}
	pv := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())

	///tcp://0.0.0.0:26656
	info := strings.Split(cfg.P2P.ListenAddress, "://")
	if len(info) == 2 {
		cfg.P2P.ListenAddress = info[0] + "://" + nodeDomain + "#" + info[1]
	}else {
		return nil, errors.New(fmt.Sprintf("unexpected p2p listen address: %v", cfg.P2P.ListenAddress))
	}

	tmNode, err := node.NewNode(
		cfg,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		ctx.Logger.With("module", "node"),
		)
	if err != nil{
		return nil, err
	}

	err = tmNode.Start()
	if err != nil {
		return nil, err
	}
	ctx.Logger.Info("Starting Node Server Success")

	// Sleep forever and then...
	tos.TrapSignal(ctx.Logger, func() {
		tmNode.Stop()
	})

	return tmNode, nil
}

func preSetConfig(ctx *app.Context) {
	cfg := ctx.Config
	genesis := viper.GetString(flagGenesis)
	nodeKey := viper.GetString(flagNodeKey)
	pvs := viper.GetString(flagPvs)
	pvk := viper.GetString(flagPvk)
	genesisBytes, _ := base64.StdEncoding.DecodeString(genesis)
	ioutil.WriteFile(cfg.GenesisFile(), genesisBytes, os.ModePerm)
	nodeKeyBytes, _ := base64.StdEncoding.DecodeString(nodeKey)
	ioutil.WriteFile(cfg.NodeKeyFile(), nodeKeyBytes, os.ModePerm)
	pvsBytes, _ := base64.StdEncoding.DecodeString(pvs)
	ioutil.WriteFile(cfg.PrivValidatorStateFile(), pvsBytes, os.ModePerm)
	pvkBytes, _ := base64.StdEncoding.DecodeString(pvk)
	ioutil.WriteFile(cfg.PrivValidatorKeyFile(), pvkBytes, os.ModePerm)
}