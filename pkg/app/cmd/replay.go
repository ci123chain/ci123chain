package cmd

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	memmock "github.com/tendermint/tendermint/mempool/mock"
	"github.com/tendermint/tendermint/node"
	proto_type "github.com/tendermint/tendermint/proto/tendermint/state"
	"github.com/tendermint/tendermint/proxy"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"log"
	"path/filepath"
)

const (
	StartBlockHeight = 0
	blockStoreDB  = "blockstore"
	stateDB       = "state"
	flagHaltHeight = "halt_height"
	//pprofAddrFlag    = "pprof_addr"
	//runWithPprofFlag = "gen_pprof"
	//
	//saveBlock = "save_block"
	//
	//defaulPprofFileFlags = os.O_RDWR | os.O_CREATE | os.O_APPEND
	//defaultPprofFilePerm = 0644
)

func replayCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replay",
		Short: "Replay blocks from local db",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- replay start ---------")
			//pprofAddress := viper.GetString(pprofAddrFlag)
			//go func() {
			//	err := http.ListenAndServe(pprofAddress, nil)
			//	if err != nil {
			//		fmt.Println(err)
			//	}
			//}()
			id := viper.GetInt64(flagETHChainID)
			limit := viper.GetInt(flagIteratorLimit)
			util.Setup(id)
			util.SetLimit(limit)

			replayBlock(ctx)
			log.Println("--------- replay success ---------")
		},
	}
	cmd.Flags().IntP(flagHaltHeight, "", 0, "height for stop replaying")
	cmd.Flags().Int64(flagETHChainID, 1, "eth chain id")
	cmd.Flags().Int(flagIteratorLimit, 10, "iterator limit")
	//cmd.Flags().StringP(pprofAddrFlag, "p", "0.0.0.0:26661", "Address and port of pprof HTTP server listening")
	//cmd.Flags().BoolVarP(&state.IgnoreSmbCheck, "ignore-smb", "i", false, "ignore state machine broken")
	//cmd.Flags().String(server.FlagPruning, storetypes.PruningOptionNothing, "Pruning strategy (default|nothing|everything|custom)")
	//cmd.Flags().Uint64(server.FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	//cmd.Flags().Bool(config.FlagPprofAutoDump, false, "Enable auto dump pprof")
	//cmd.Flags().String(config.FlagPprofCollectInterval, "5s", "Interval for pprof dump loop")
	//cmd.Flags().Int(config.FlagPprofCpuTriggerPercentMin, 45, "TriggerPercentMin of cpu to dump pprof")
	//cmd.Flags().Int(config.FlagPprofCpuTriggerPercentDiff, 50, "TriggerPercentDiff of cpu to dump pprof")
	//cmd.Flags().Int(config.FlagPprofCpuTriggerPercentAbs, 50, "TriggerPercentAbs of cpu to dump pprof")
	//cmd.Flags().Int(config.FlagPprofMemTriggerPercentMin, 70, "TriggerPercentMin of mem to dump pprof")
	//cmd.Flags().Int(config.FlagPprofMemTriggerPercentDiff, 50, "TriggerPercentDiff of mem to dump pprof")
	//cmd.Flags().Int(config.FlagPprofMemTriggerPercentAbs, 75, "TriggerPercentAbs of cpu mem dump pprof")
	//cmd.Flags().String(config.FlagPprofCoolDown, "3m", "The cool down time after every type of pprof dump")
	//cmd.Flags().Int64(config.FlagPprofAbciElapsed, 5000, "Elapsed time of abci in millisecond for pprof dump")
	//cmd.Flags().Bool(config.FlagPprofUseCGroup, false, "Use cgroup when exchaind run in docker")
	//cmd.Flags().IntVar(&iavl.IavlCacheSize, iavl.FlagIavlCacheSize, 1000000, "Max size of iavl cache")
	//cmd.Flags().StringToIntVar(&tmiavl.OutputModules, tmiavl.FlagOutputModules, map[string]int{}, "decide which module in iavl to be printed")
	//cmd.Flags().Int64Var(&tmiavl.CommitIntervalHeight, tmiavl.FlagIavlCommitIntervalHeight, 100, "Max interval to commit node cache into leveldb")
	//cmd.Flags().Int64Var(&tmiavl.MinCommitItemCount, tmiavl.FlagIavlMinCommitItemCount, 500000, "Min nodes num to triggle node cache commit")
	//cmd.Flags().IntVar(&tmiavl.HeightOrphansCacheSize, tmiavl.FlagIavlHeightOrphansCacheSize, 8, "Max orphan version to cache in memory")
	//cmd.Flags().IntVar(&tmiavl.MaxCommittedHeightNum, tmiavl.FlagIavlMaxCommittedHeightNum, 8, "Max committed version to cache in memory")
	//cmd.Flags().BoolVar(&tmiavl.EnableAsyncCommit, tmiavl.FlagIavlEnableAsyncCommit, false, "Enable cache iavl node data to optimization leveldb pruning process")
	//cmd.Flags().Bool(runWithPprofFlag, false, "Dump the pprof of the entire replay process")
	//cmd.Flags().Bool(sm.FlagParalleledTx, false, "pall Tx")
	//cmd.Flags().Bool(saveBlock, false, "save block when replay")
	return cmd
}


// replayBlock replays blocks from db, if something goes wrong, it will panic with error message.
func replayBlock(ctx *app.Context) {
	proxyApp, err := createProxyApp(ctx)
	panicError(err)

	res, err := proxyApp.Query().InfoSync(proxy.RequestInfo)
	panicError(err)
	currentBlockHeight := res.LastBlockHeight
	currentAppHash := res.LastBlockAppHash
	log.Println("current block height", "height", currentBlockHeight)
	log.Println("current app hash", "appHash", fmt.Sprintf("%X", currentAppHash))

	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	stateStoreDB, err := openDB(stateDB, dataDir)
	panicError(err)

	stateStore := sm.NewStore(stateStoreDB)

	genesisDocProvider := node.DefaultGenesisDocProviderFunc(ctx.Config)
	state, genDoc, err := node.LoadStateFromDBOrGenesisDocProvider(stateStoreDB, genesisDocProvider)
	panicError(err)

	// If startBlockHeight == 0 it means that we are at genesis and hence should initChain.
	if currentBlockHeight == 0 {
		err := initChain(state, stateStore, genDoc, proxyApp)
		panicError(err)
		state, _ = stateStore.Load()
		// blockchain start at no-zero height
		if state.InitialHeight > 0 {
			currentBlockHeight = state.InitialHeight - 1
		}
	}
	log.Println("current block height", "height", currentBlockHeight)
	log.Println("current app hash", "appHash", fmt.Sprintf("%X", currentAppHash))

	// replay
	doReplay(ctx, state, stateStore, proxyApp, currentAppHash, currentBlockHeight)
}

// panic if error is not nil
func panicError(err error) {
	if err != nil {
		panic(err)
	}
}


func openDB(dbName string, dataDir string) (db dbm.DB, err error) {
	return sdk.NewLevelDB(dbName, dataDir)
}

func createProxyApp(ctx *app.Context) (proxy.AppConns, error) {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")

	db, err := openDB(app.AppName, dataDir)
	panicError(err)
	app := app.NewApp(ctx.Logger, db, nil, nil)
	clientCreator := proxy.NewLocalClientCreator(app)
	return createAndStartProxyAppConns(clientCreator)
}


func createAndStartProxyAppConns(clientCreator proxy.ClientCreator) (proxy.AppConns, error) {
	proxyApp := proxy.NewAppConns(clientCreator)
	if err := proxyApp.Start(); err != nil {
		return nil, fmt.Errorf("error starting proxy app connections: %v", err)
	}
	return proxyApp, nil
}



func doReplay(ctx *app.Context, state sm.State, stateStoreDB sm.Store,
	proxyApp proxy.AppConns, lastAppHash []byte, lastBlockHeight int64) {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	originBlockStoreDB, err := openDB(blockStoreDB, dataDir)
	panicError(err)
	originBlockStore := store.NewBlockStore(originBlockStoreDB)
	originLatestBlockHeight := originBlockStore.Height()
	log.Println("origin latest block height", "height", originLatestBlockHeight)

	haltheight := viper.GetInt64(flagHaltHeight)
	if haltheight == 0 {
		haltheight = originLatestBlockHeight
	}
	if haltheight <= lastBlockHeight+1 {
		panic("haltheight <= startBlockHeight please check data or height")
	}

	if haltheight > originLatestBlockHeight {
		panic("haltheight > originLatestBlockHeight,  please check data or height")
	}

	log.Println("replay stop block height", "height", haltheight)

	// Replay blocks up to the latest in the blockstore.
	if lastBlockHeight == state.LastBlockHeight+1 {
		abciResponses, err := stateStoreDB.LoadABCIResponses(lastBlockHeight)
		//abciResponses, err := sm.LoadABCIResponses(stateStoreDB, lastBlockHeight)
		panicError(err)
		mockApp := newMockProxyApp(lastAppHash, abciResponses)
		block := originBlockStore.LoadBlock(lastBlockHeight)
		meta := originBlockStore.LoadBlockMeta(lastBlockHeight)
		blockExec := sm.NewBlockExecutor(stateStoreDB, ctx.Logger, mockApp, memmock.Mempool{}, sm.EmptyEvidencePool{})
		//blockExec.SetIsAsyncDeliverTx(false) // mockApp not support parallel tx
		state, _, err = blockExec.ApplyBlock(state, meta.BlockID, block)
		panicError(err)
	}

	blockExec := sm.NewBlockExecutor(stateStoreDB, ctx.Logger, proxyApp.Consensus(), memmock.Mempool{}, sm.EmptyEvidencePool{})
	//if viper.GetBool(runWithPprofFlag) {
	//	startDumpPprof()
	//	defer stopDumpPprof()
	//}
	//needSaveBlock := viper.GetBool(saveBlock)
	for height := lastBlockHeight + 1; height <= haltheight; height++ {

		block := originBlockStore.LoadBlock(height)
		meta := originBlockStore.LoadBlockMeta(height)
		//blockExec.SetIsAsyncDeliverTx(viper.GetBool(sm.FlagParalleledTx))
		state, _, err = blockExec.ApplyBlock(state, meta.BlockID, block)
		panicError(err)
		log.Println("replaying: ", height, " apphash: ", fmt.Sprintf("%X", state.AppHash))
	}
}


func initChain(state sm.State, stateStore sm.Store, genDoc *types.GenesisDoc, proxyApp proxy.AppConns) error {
	validators := make([]*types.Validator, len(genDoc.Validators))
	for i, val := range genDoc.Validators {
		validators[i] = types.NewValidator(val.PubKey, val.Power)
	}
	validatorSet := types.NewValidatorSet(validators)
	nextVals := types.TM2PB.ValidatorUpdates(validatorSet)
	csParams := types.TM2PB.ConsensusParams(genDoc.ConsensusParams)
	req := abci.RequestInitChain{
		Time:            genDoc.GenesisTime,
		ChainId:         genDoc.ChainID,
		InitialHeight:   genDoc.InitialHeight,
		ConsensusParams: csParams,
		Validators:      nextVals,
		AppStateBytes:   genDoc.AppState,
	}
	res, err := proxyApp.Consensus().InitChainSync(req)
	if err != nil {
		return err
	}
	if state.LastBlockHeight == StartBlockHeight { //we only update state when we are in initial state
		// If the app returned validators or consensus params, update the state.
		if len(res.Validators) > 0 {
			vals, err := types.PB2TM.ValidatorUpdates(res.Validators)
			if err != nil {
				return err
			}
			state.Validators = types.NewValidatorSet(vals)
			state.NextValidators = types.NewValidatorSet(vals).CopyIncrementProposerPriority(1)
		} else if len(genDoc.Validators) == 0 {
			// If validator set is not set in genesis and still empty after InitChain, exit.
			return fmt.Errorf("validator set is nil in genesis and still empty after InitChain")
		}

		if res.ConsensusParams != nil {
			state.ConsensusParams = types.UpdateConsensusParams(state.ConsensusParams, res.ConsensusParams)
			state.Version.Consensus.App = state.ConsensusParams.Version.AppVersion
		}
		// We update the last results hash with the empty hash, to conform with RFC-6962.
		state.LastResultsHash = merkle.HashFromByteSlices(nil)

		err := stateStore.Save(state)
		panicError(err)
	}
	return nil
}

func newMockProxyApp(appHash []byte, abciResponses *proto_type.ABCIResponses) proxy.AppConnConsensus {
	clientCreator := proxy.NewLocalClientCreator(&mockProxyApp{
		appHash:       appHash,
		abciResponses: abciResponses,
	})
	cli, _ := clientCreator.NewABCIClient()
	err := cli.Start()
	if err != nil {
		panic(err)
	}
	return proxy.NewAppConnConsensus(cli)
}

type mockProxyApp struct {
	abci.BaseApplication

	appHash       []byte
	txCount       int
	abciResponses *proto_type.ABCIResponses
}

func (mock *mockProxyApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	r := mock.abciResponses.DeliverTxs[mock.txCount]
	mock.txCount++
	if r == nil { //it could be nil because of amino unMarshall, it will cause an empty ResponseDeliverTx to become nil
		return abci.ResponseDeliverTx{}
	}
	return *r
}

func (mock *mockProxyApp) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	mock.txCount = 0
	return *mock.abciResponses.EndBlock
}

func (mock *mockProxyApp) Commit() abci.ResponseCommit {
	return abci.ResponseCommit{Data: mock.appHash}
}

