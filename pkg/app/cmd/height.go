package cmd

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/spf13/cobra"
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
			if err != nil {
				if _, err = leveldb.RecoverFile(dataDir, nil); err != nil {
					panic(fmt.Sprintf(`err while recoverfile%s : %s`, dataDir, err.Error()))
				}
			}
			panicError(err)
			originBlockStore := store.NewBlockStore(originBlockStoreDB)
			originLatestBlockHeight := originBlockStore.Height()
			fmt.Println(originLatestBlockHeight)
			return nil
		},
	}
	return cmd
}
