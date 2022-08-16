package cmd

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/store"
	"path/filepath"
)

func storeHeight(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store-height",
		Short: "Print the store block height",
		RunE: func(cmd *cobra.Command, args []string) error {
			rootDir := ctx.Config.RootDir
			dataDir := filepath.Join(rootDir, "data")
			originBlockStoreDB, err := openDB(blockStoreDB, dataDir)
			panicError(err)
			originBlockStore := store.NewBlockStore(originBlockStoreDB)
			originLatestBlockHeight := originBlockStore.Height()
			fmt.Println(originLatestBlockHeight)
			return nil
		},
	}
	return cmd
}

func recoverLevelDB(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recover-leveldb",
		Short: "recover leveldb",
		RunE: func(cmd *cobra.Command, args []string) error {
			rootDir := ctx.Config.RootDir
			dataDir := filepath.Join(rootDir, "data")
			filename := viper.GetString("filename")
			db, err := leveldb.RecoverFile(dataDir+"/"+filename+".db", nil)
			if db == nil {
				panic("DB nil")
			}
			if err != nil {
				panic(err)
			}
			fmt.Println("Success")
			db.Close()
			return nil
		},
	}
	cmd.Flags().String("filename", "", "recover filename in data directory")
	return cmd
}
