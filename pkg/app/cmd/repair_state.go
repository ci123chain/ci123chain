package cmd

import (
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/spf13/cobra"
	"log"
)

const (
	FlagStartHeight string = "start-height"
)

func repairStateCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair-state",
		Short: "Repair the SMB(state machine broken) data of node",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- repair data start ---------")

			repairState(ctx)
			log.Println("--------- repair data success ---------")
		},
	}
	cmd.Flags().Int64(FlagStartHeight, 0, "Set the start block height for repair")
	return cmd
}



func repairState(ctx *app.Context) {
	// set ignore smb check
	//sm.SetIgnoreSmbCheck(true)
	//iavl.SetIgnoreVersionCheck(true)
	//
	//// load latest block height
	//rootDir := ctx.Config.RootDir
	//dataDir := filepath.Join(rootDir, "data")
	//latestBlockHeight := latestBlockHeight(dataDir)
	//startBlockHeight := types.GetStartBlockHeight()
	//if latestBlockHeight <= startBlockHeight+2 {
	//	panic(fmt.Sprintf("There is no need to repair data. The latest block height is %d, start block height is %d", latestBlockHeight, startBlockHeight))
	//}
	//
	//// create proxy app
	//proxyApp, repairApp, err := createRepairApp(ctx)
	//panicError(err)
	//
	//// load state
	//stateStoreDB, err := openDB(stateDB, dataDir)
	//panicError(err)
	//genesisDocProvider := node.DefaultGenesisDocProviderFunc(ctx.Config)
	//state, _, err := node.LoadStateFromDBOrGenesisDocProvider(stateStoreDB, genesisDocProvider)
	//panicError(err)
	//
	//// load start version
	//startVersion := viper.GetInt64(FlagStartHeight)
	//if startVersion == 0 {
	//	latestVersion := repairApp.getLatestVersion()
	//	startVersion = latestVersion - 2
	//}
	//if startVersion == 0 {
	//	panic("height too low, please restart from height 0 with genesis file")
	//}
	//err = repairApp.LoadStartVersion(startVersion)
	//panicError(err)
	//
	//// repair data by apply the latest two blocks
	//doRepair(ctx, state, stateStoreDB, proxyApp, startVersion, latestBlockHeight, dataDir)
	//repairApp.StopStore()
}


