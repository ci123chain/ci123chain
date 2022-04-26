package cmd

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/version"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	hnode "github.com/ci123chain/ci123chain/pkg/node"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abcis "github.com/tendermint/tendermint/abci/server"
	v0 "github.com/tendermint/tendermint/blockchain/v0"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/node"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/rpc/client/local"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	flagCiNodeDomain   = "IDG_HOST_80"
	flagShardIndex     = "shardIndex"

	flagETHChainID     = "eth_chain_id"
	flagIteratorLimit  = "iterator_limit"
	flagRunMode		   = "mode"
	flagStartFromExport = "export"
	flagStartFromExportFile = "export_file"

	flagSkipWAL  = "skipwal"
)

func startCmd(ctx *app.Context, appCreator app.AppCreator, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "start",
		Short: "Run the full node",
		RunE: func(cmd *cobra.Command, args []string) error {
			limit := viper.GetInt(flagIteratorLimit)
			util.Setup(int64(ctx.Config.EthChainID))
			util.SetLimit(limit)
			v0.SetSkipWAL(viper.GetBool(flagSkipWAL))
			loadExport := viper.GetBool(flagStartFromExport)
			exportFile := viper.GetString(flagStartFromExportFile)

			filePath := filepath.Join(ctx.Config.RootDir, "config/genesis.json")
			if loadExport && !tos.FileExists(filePath){
				err := downloadGenesis(exportFile, filePath)
				//by, err := ioutil.ReadFile(exportFile)
				if err != nil {
					ctx.Logger.Error("start node with export genesis file, but download fail ", "path", exportFile, "err", err.Error())
					panic(err)
				} else {
					ctx.Logger.Info("Load Export Genesis OK", "path", exportFile)
				}
			}
			if !viper.GetBool(flagWithTendermint) {
				ctx.Logger.Info("Starting ABCI Without Tendermint")
				return startStandAlone(ctx, appCreator)
			}
			ctx.Logger.Info("Version: ", version.NewInfo().Version)
			ctx.Logger.Info("Starting ABCI with Tendermint")

			_, err := StartInProcess(ctx, appCreator, cdc)
			if err != nil {
				return err
			}
			select {}
		},
	}

	cmd.Flags().Bool(flagWithTendermint, true, "Run abci app embedded in-process with tendermint")
	cmd.Flags().String(FlagListenAddr, "tcp://0.0.0.0:1317", "The address for the server to listen on")
	cmd.Flags().String(helper.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(flagPruning, "syncable", "Pruning strategy: syncable, nothing, everything")
	cmd.Flags().String(flagCiStateDBType, "redis", "database types")
	cmd.Flags().String(flagCiStateDBHost, "", "db host")
	cmd.Flags().Uint64(flagCiStateDBPort, 7443, "db port")
	cmd.Flags().Bool(flagCiStateDBTls, true, "use tls")
	cmd.Flags().String(flagCiNodeDomain, "", "node domain")
	cmd.Flags().String(flagShardIndex, "", "index of shard")
	cmd.Flags().Int(flagIteratorLimit, 10, "iterator limit")
	cmd.Flags().String(app.FlagValidatorKey, "", "validator_key")
	cmd.Flags().String(flagRunMode, "light", "run chain mode")
	cmd.Flags().Bool(flagStartFromExport, false, "start with export file")
	cmd.Flags().String(flagStartFromExportFile, "/opt/exportFile.json", "start with export file")
	cmd.Flags().String(FlagWebsocket, "8546", "websocket port to listen to")
	cmd.Flags().Uint(FlagMaxOpenConnections, 1000, "The number of maximum open connections")
	cmd.Flags().Uint(FlagRPCReadTimeout, 60, "The RPC read timeout")
	cmd.Flags().Uint(FlagRPCWriteTimeout, 60, "The RPC write timeout")
	cmd.Flags().String(flagTokenName, "stake", "Chain token name")
	cmd.Flags().String(app.FlagMasterDomain, "", "master node")

	cmd.Flags().Bool(flagSkipWAL, false, "start with skipWAL")

	//cmd.Flags().String(flagLogLevel, "debug", "Run abci app with different log level")
	tcmd.AddNodeFlags(cmd)
	return cmd
}

func downloadGenesis(genesisUrl string, filepath string) error {
	client := http.DefaultClient;
	client.Timeout = time.Second * 60 //设置超时时间
	resp, err := client.Get(genesisUrl)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("unKnownError: " + genesisUrl)
	}
	if resp.StatusCode != http.StatusOK {
		ret, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(ret))
	}
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(file, resp.Body)
	return err
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

func StartInProcess(ctx *app.Context, appCreator app.AppCreator, cdc *codec.Codec) (*node.Node, error) {
	cfg := ctx.Config
	home := cfg.RootDir
	traceStore := viper.GetString(flagTraceStore)
	stateDB := ""

	lightMode := viper.GetString(flagRunMode)
	if lightMode != store.ModeLight {
		dbType := viper.GetString(flagCiStateDBType)
		if dbType == "" {
			dbType = "redis"
		}
		dbHost := viper.GetString(flagCiStateDBHost)
		if dbHost == "" {
			var err error
			dbHost, err = util.GetDomain()
			if err != nil {
				ctx.Logger.Error("get remote db host failed", "err", err.Error())
				return nil, err
			}
		}
		ctx.Logger.Info("discovery remote db host", "host", dbHost)
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
			return nil, errors.New(fmt.Sprintf("types of db: %s, which is not reids not implement yet", dbType))
		}
	}



	//nodeDomain := viper.GetString(flagCiNodeDomain)
	nodeDomain := os.Getenv(flagCiNodeDomain)
	if nodeDomain == "" {
		nodeDomain = util.GetLocalAddress()
		if nodeDomain == "" {
			return nil, errors.New("you have no valid ip address which used to be dial by other peer")
		}
	}

	appState, gendoc, err := app.GenesisStateFromGenFile(cdc, cfg.GenesisFile())
	if err != nil {
		return nil, err
	}
	var stakingGenesisState staking.GenesisState
	if _, ok := appState[staking.ModuleName]; !ok{
		return nil, errors.New("unexpected genesisState of staking")
	} else {
		cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState)
	}
	types.SetCoinDenom(stakingGenesisState.Params.BondDenom)
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

	///tcp://0.0.0.0:26656 change to custom domain, send to other peer
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

	cliCtx, err := client.NewClientContext()
	if err != nil {
		return nil, err
	}
	cliCtx.WithClient(local.New(tmNode))
	app.RegisterTxService(cliCtx)

	go func() {
		for {
			err := StartRestServer(cdc, tmNode, viper.GetString(FlagListenAddr))
			ctx.Logger.Error("Rest-Server ", "err", err.Error())
		}
	}()

	// Sleep forever and then...
	tos.TrapSignal(ctx.Logger, func() {
		tmNode.Stop()
	})

	return tmNode, nil
}
