package cmd

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/baseapp"
	abcitypes "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/cosmos/iavl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	memmock "github.com/tendermint/tendermint/mempool/mock"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"io"

	"path/filepath"
)

const (
	FlagStartHeight string = "start_height"
)

func repairStateCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair-state",
		Short: "Repair the SMB(state machine broken) data of node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("--------- repair data start ---------")

			limit := viper.GetInt(flagIteratorLimit)
			util.Setup(int64(ctx.Config.EthChainID))
			util.SetLimit(limit)

			{
				cfg := ctx.Config
				cdc := types2.GetCodec()
				appState, _, _ := app.GenesisStateFromGenFile(cdc, cfg.GenesisFile())
				var stakingGenesisState staking.GenesisState
				cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState)
				abcitypes.SetCoinDenom(stakingGenesisState.Params.BondDenom)
			}

			repairState(ctx)
			fmt.Println("--------- repair data success ---------")
		},
	}
	cmd.Flags().Int64(FlagStartHeight, 0, "Set the start block height for repair")
	cmd.Flags().Int(flagIteratorLimit, 10, "iterator limit")
	return cmd
}

type repairApp struct {
	db dbm.DB
	*app.Chain
}

func createRepairApp(ctx *app.Context) (proxy.AppConns, *repairApp, error) {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := openDB(app.AppName, dataDir)
	panicError(err)
	repairApp := newRepairApp(ctx.Logger, db, nil)

	clientCreator := proxy.NewLocalClientCreator(repairApp)
	// Create the proxyApp and establish connections to the ABCI app (consensus, mempool, query).
	proxyApp, err := createAndStartProxyAppConns(clientCreator)
	return proxyApp, repairApp, err
}

func newRepairApp(logger log.Logger, db dbm.DB, traceStore io.Writer) *repairApp {
	return &repairApp{db, app.NewChain(
		logger,
		db,
		nil,
		nil,
		baseapp.SetGasPriceConfig(app.GasPriceConfig),
	)}
}

func (app *repairApp) getLatestVersion() int64 {
	rs := initAppStore(app.db)
	return rs.GetLatestVersion()
}

func repairState(ctx *app.Context) {
	sm.SetIgnoreSmbCheck(true)
	iavl.SetIgnoreVersionCheck(true)

	// load latest block height
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	latestBlockHeight := latestBlockHeight(dataDir)

	// create proxy app
	proxyApp, repairApp, err := createRepairApp(ctx)
	panicError(err)

	// load state
	stateStoreDB, err := openDB(stateDB, dataDir)
	panicError(err)
	genesisDocProvider := node.DefaultGenesisDocProviderFunc(ctx.Config)
	state, _, err := node.LoadStateFromDBOrGenesisDocProvider(stateStoreDB, genesisDocProvider)
	panicError(err)

	stateStoreDBBack, err := openDB("state_back", dataDir)
	panicError(err)

	// load start version
	startVersion := viper.GetInt64(FlagStartHeight)
	if startVersion == 0 {
		latestVersion := repairApp.getLatestVersion()
		startVersion = latestVersion - 2
	}
	if startVersion == 0 {
		panic("height too low, please restart from height 0 with genesis file")
	}
	err = repairApp.LoadStartVersion(startVersion)
	panicError(err)

	// repair data by apply the latest two blocks
	doRepair(ctx, state, stateStoreDB, stateStoreDBBack, proxyApp, startVersion, latestBlockHeight, dataDir)
	//repairApp.StopStore()
}

func doRepair(ctx *app.Context, state sm.State, stateStoreDB, stateStoreDBBack dbm.DB,
	proxyApp proxy.AppConns, startHeight, latestHeight int64, dataDir string) {
	stateCopy := state.Copy()
	state = constructStartState(state, stateStoreDB, startHeight+1)
	smstore := sm.NewStore(stateStoreDB)
	var err error
	//smstoreBack := sm.NewStore(stateStoreDBBack)
	//backVal, err := smstoreBack.LoadValidators(startHeight)
	//fmt.Println("BackUpa: ", backVal.Proposer.String())
	//fmt.Println("Current: ", state.Validators.Proposer.String())

	blockExec := sm.NewBlockExecutor(smstore, ctx.Logger, proxyApp.Consensus(), memmock.Mempool{}, sm.EmptyEvidencePool{})
	for height := startHeight + 1; height <= latestHeight; height++ {
		repairBlock, repairBlockMeta := loadBlock(height, dataDir)

		state, _, err = blockExec.ApplyBlock(state, repairBlockMeta.BlockID, repairBlock)
		panicError(err)

		//backVal, err = smstoreBack.LoadValidators(height)
		//fmt.Println("BackUpa: ", backVal.Proposer.String())
		//fmt.Println("Current: ", state.Validators.Proposer.String())

		//use stateCopy to correct the repaired state
		if state.LastBlockHeight == stateCopy.LastBlockHeight {
			//fmt.Println("newState", fmt.Sprintf("%+v", state))
			//fmt.Println("oldState", fmt.Sprintf("%+v", stateCopy))

			state.LastHeightConsensusParamsChanged = stateCopy.LastHeightConsensusParamsChanged
			state.LastHeightValidatorsChanged = stateCopy.LastHeightValidatorsChanged
			//state.LastValidators = stateCopy.LastValidators.Copy()
			//state.Validators = stateCopy.Validators.Copy()
			//state.NextValidators = state.NextValidators.Copy()
			//smstore.Save(state)
			//sm.SaveState(stateStoreDB, state)
		}
		res, err := proxyApp.Query().InfoSync(proxy.RequestInfo)
		panicError(err)
		repairedBlockHeight := res.LastBlockHeight
		repairedAppHash := res.LastBlockAppHash
		fmt.Println("Repaired block height", repairedBlockHeight)
		fmt.Println("Repaired app hash", fmt.Sprintf("%X", repairedAppHash))
	}
}

func constructStartState(state sm.State, stateStoreDB dbm.DB, startHeight int64) sm.State {
	stateCopy := state.Copy()
	smstore := sm.NewStore(stateStoreDB)

	validators, err := smstore.LoadValidators(startHeight)
	lastValidators, err := smstore.LoadValidators(startHeight - 1)
	if err != nil {
		return stateCopy
	}
	nextValidators, err := smstore.LoadValidators(startHeight + 1)
	if err != nil {
		return stateCopy
	}
	consensusParams, err := smstore.LoadConsensusParams(startHeight + 1)
	if err != nil {
		return stateCopy
	}
	stateCopy.Validators = validators
	stateCopy.LastValidators = lastValidators
	stateCopy.NextValidators = nextValidators
	stateCopy.ConsensusParams = consensusParams
	stateCopy.LastBlockHeight = startHeight
	return stateCopy
}

func loadBlock(height int64, dataDir string) (*types.Block, *types.BlockMeta) {
	//rootDir := ctx.Config.RootDir
	//dataDir := filepath.Join(rootDir, "data")
	storeDB, err := openDB(blockStoreDB, dataDir)
	defer storeDB.Close()
	blockStore := store.NewBlockStore(storeDB)
	panicError(err)
	block := blockStore.LoadBlock(height)
	meta := blockStore.LoadBlockMeta(height)
	return block, meta
}

func latestBlockHeight(dataDir string) int64 {
	storeDB, err := openDB(blockStoreDB, dataDir)
	panicError(err)
	defer storeDB.Close()
	blockStore := store.NewBlockStore(storeDB)
	return blockStore.Height()
}
