package cmd

import (
	"fmt"
	abci_store "github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/auth"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	distr "github.com/ci123chain/ci123chain/pkg/distribution"
	"github.com/ci123chain/ci123chain/pkg/gravity"
	ibctransfertypes "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	ibchost "github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/ci123chain/ci123chain/pkg/infrastructure"
	"github.com/ci123chain/ci123chain/pkg/mint"
	"github.com/ci123chain/ci123chain/pkg/order"
	"github.com/ci123chain/ci123chain/pkg/params"
	prestaking "github.com/ci123chain/ci123chain/pkg/pre_staking"
	"github.com/ci123chain/ci123chain/pkg/registry"
	"github.com/ci123chain/ci123chain/pkg/slashing"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/upgrade"
	"github.com/ci123chain/ci123chain/pkg/vm"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"reflect"

	"github.com/tendermint/tendermint/store"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
	dbm "github.com/tendermint/tm-db"
	"log"
	"time"
)

const flagPureHeightFrom = "pure_height_from"
const flagPureHeightTo = "pure_height_to"


const (
	flagHeight    = "height"
	flagDBBackend = "db_backend"

	blockDBName = "blockstore"
	stateDBName = "state"
	appDBName   = "ci123"
)

func pureStateCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prueState",
		Short: "prue state latest state",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- prue start ---------")
			pureHeightFrom := viper.GetInt(flagPureHeightFrom)
			if pureHeightFrom < 1 {
				fmt.Println("pure_height_from not provided")
			}
			pureHeightTo := viper.GetInt(flagPureHeightTo)
			if pureHeightTo <= pureHeightFrom {
				fmt.Println("pure_height_to invalid")
			}
			pureState(ctx, int64(pureHeightFrom), int64(pureHeightTo))
			log.Println("--------- replay success ---------")
		},
	}
	cmd.Flags().IntP(flagPureHeightFrom, "", 0, "from height")
	cmd.Flags().IntP(flagPureHeightTo, "", 0, "to height")
	cmd.Flags().String(flags.FlagHome, "", "The application home directory")

	return cmd
}

func pureState(ctx *app.Context, fromHeight, toHeight int64)  {
	config := ctx.Config
	config.SetRoot(viper.GetString(flags.FlagHome))

	blockStoreDB := initDB(config, blockDBName)
	stateDB := initDB(config, stateDBName)
	appDB := initDB(config, appDBName)


	pruneApp(appDB, fromHeight, toHeight)
	pruneStates(stateDB, fromHeight, toHeight)
	pruneBlocks(blockStoreDB, fromHeight, toHeight)
}

func initDB(config *cfg.Config, dbName string) dbm.DB {
	if dbName != blockDBName && dbName != stateDBName && dbName != appDBName {
		panic(fmt.Sprintf("unknow db name:%s", dbName))
	}

	db, err := node.DefaultDBProvider(&node.DBContext{dbName, config})
	panicError(err)

	return db
}


// pruneBlocks deletes blocks between the given heights (including from, excluding to).
func pruneBlocks(blockStoreDB dbm.DB, baseHeight, retainHeight int64) {

	log.Printf("Prune blocks [%d,%d)...", baseHeight, retainHeight)
	if retainHeight <= baseHeight {
		return
	}

	baseHeightBefore, sizeBefore := getBlockInfo(blockStoreDB)
	start := time.Now()
	blockstore := store.NewBlockStore(blockStoreDB)
	blockstore.Base()
	_, err := blockstore.PruneBlocks(retainHeight)
	if err != nil {
		panic(fmt.Errorf("failed to prune block store: %w", err))
	}

	baseHeightAfter, sizeAfter := getBlockInfo(blockStoreDB)
	log.Printf("Block db info [baseHeight,size]: [%d,%d] --> [%d,%d]\n", baseHeightBefore, sizeBefore, baseHeightAfter, sizeAfter)
	log.Printf("Prune blocks done in %v \n", time.Since(start))
}

// pruneStates deletes states between the given heights (including from, excluding to).
func pruneStates(stateDB dbm.DB, from, to int64) {

	log.Printf("Prune states [%d,%d)...", from, to)
	if to <= from {
		return
	}

	start := time.Now()
	store := sm.NewStore(stateDB)

	if err := store.PruneStates(from, to); err != nil {
		panic(fmt.Errorf("failed to prune state database: %w", err))
	}

	log.Printf("Prune states done in %v \n", time.Since(start))
}

// pruneApp deletes app states between the given heights (including from, excluding to).
func pruneApp(appDB dbm.DB, from, to int64) {

	log.Printf("Prune applcation [%d,%d)...", from, to)
	if to <= from {
		return
	}

	rs := initAppStore(appDB)
	latestV := rs.GetLatestVersion()
	if to > latestV {
		return
	}
	//versions := rs.GetVersions()
	//if len(versions) == 0 {
	//	return
	//}
	//pruneHeights := rs.GetPruningHeights()

	//newVersions := make([]int64, 0)
	//newPruneHeights := make([]int64, 0)
	//deleteVersions := make([]int64, 0)
	//
	//for _, v := range pruneHeights {
	//	if v >= to || v < from {
	//		newPruneHeights = append(newPruneHeights, v)
	//		continue
	//	}
	//	deleteVersions = append(deleteVersions, v)
	//}
	//
	//for _, v := range versions {
	//	if v >= to || v < from {
	//		newVersions = append(newVersions, v)
	//		continue
	//	}
	//	deleteVersions = append(deleteVersions, v)
	//}
	//log.Printf("Prune application: Versions=%v, PruneVersions=%v", len(versions)+len(pruneHeights), len(deleteVersions))

	keysNumBefore, kvSizeBefore := calcKeysNum(appDB)
	start := time.Now()
	for key, store := range rs.GetStores() {
		if store.GetStoreType() == sdk.StoreTypeIAVL {
			// If the store is wrapped with an inter-block cache, we must first unwrap
			// it to get the underlying IAVL store.
			store = rs.GetCommitKVStore(key)

			if reflect.TypeOf(store).Elem() == reflect.TypeOf(abci_store.IavlStore{}){

			}

			if err := store.(*abci_store.IavlStore).DeleteVersions(from, to); err != nil {
				log.Printf("failed to delete version: %s", err)
			}
		}
	}

	commitID := rs.Commit()
	log.Printf("Prune application done commitID %v \n", commitID)

	//rs.FlushPruneHeights(newPruneHeights, newVersions)
	//
	keysNumAfter, kvSizeAfter := calcKeysNum(appDB)
	log.Printf("Application db key info [keysNum,kvSize]: [%d,%d] --> [%d,%d]\n", keysNumBefore, kvSizeBefore, keysNumAfter, kvSizeAfter)
	log.Printf("Prune application done in %v \n", time.Since(start))
}


func getBlockInfo(blockStoreDB dbm.DB) (baseHeight, size int64) {
	blockStore := store.NewBlockStore(blockStoreDB)
	baseHeight = blockStore.Base()
	size = blockStore.Size()
	return
}


func initAppStore(appDB dbm.DB) sdk.CommitMultiStore {
	cms := abci_store.NewCommitMultiStore(appDB, nil, "")

	keys := sdk.NewKVStoreKeys(app.StoreKey, account.StoreKey, params.StoreKey, auth.StoreKey,
		supply.StoreKey, order.StoreKey, ibchost.StoreKey, distr.StoreKey, staking.StoreKey,
		prestaking.StoreKey, slashing.StoreKey, gravity.StoreKey, vm.StoreKey, mint.StoreKey,
		infrastructure.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey, upgrade.StoreKey, registry.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	//c.MountKVStores(keys)
	//c.MountStoreMemory(memKeys)
	//c.MountKVStoresTransient(tkeys)
	//
	//for _, key := range keys {
	//	if err := c.LoadLatestVersion(key); err != nil {
	//		return err
	//	}
	//}

	for _, key := range memKeys {
		cms.MountStoreWithDB(key, sdk.StoreTypeMemory, nil)
	}

	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
	}
	for _, key := range tkeys {
		cms.MountStoreWithDB(key, sdk.StoreTypeTransient, nil)
	}

	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	//rs, ok := cms.(*rootmulti.Store)
	//if !ok {
	//	panic("cms of from app is not rootmulti store")
	//}

	return cms
}

func calcKeysNum(db dbm.DB) (keys, kvSize uint64) {
	iter, err := db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	for ; iter.Valid(); iter.Next() {
		keys++
		kvSize += uint64(len(iter.Key())) + uint64(len(iter.Value()))
	}
	iter.Close()
	return
}