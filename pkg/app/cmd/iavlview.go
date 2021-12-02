package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	acc_types "github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/cosmos/iavl"
	"github.com/spf13/cobra"
	dbm "github.com/tendermint/tm-db"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

const (
	KeyAcc          = "accounts"
)
var printKeysDict = map[string]printKey{
	//KeyEvm:          evmPrintKey,
	KeyAcc:          accPrintKey,
	//KeyParams:       paramsPrintKey,
	//KeyStaking:      stakingPrintKey,
	//KeyGov:          govPrintKey,
	//KeyDistribution: distributionPrintKey,
	//KeySlashing:     slashingPrintKey,
	//KeyMain:         mainPrintKey,
	//KeyToken:        tokenPrintKey,
	//KeyMint:         mintPrintKey,
	//KeySupply:       supplyPrintKey,
}

type (
	printKey func(cdc *codec.Codec, key []byte, value []byte)
)

const 	DefaultCacheSize int = 100000

func iaviewerCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iaviewer",
		Short: "Iaviewer key-value from leveldb",
	}

	cmd.AddCommand(
		//readAll(cdc),
		readDiff(cdc),
	)

	return cmd
}


func readDiff(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff [data_dir] [compare_data_dir] [height] [module]",
		Short: "Read different key-value from leveldb according two paths",
		Run: func(cmd *cobra.Command, args []string) {
			var moduleList []string
			if len(args) == 4 {
				moduleList = []string{args[3]}
			}
			//else {
			//	moduleList = make([]string, 0, len(app.ModuleBasics))
			//	for m := range app.ModuleBasics {
			//		moduleList = append(moduleList, fmt.Sprintf("s/k:%s/", m))
			//	}
			//}
			height, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				panic("The input height is wrong")
			}
			IaviewerPrintDiff(cdc, args[0], args[1], moduleList, int(height))
		},
	}
	return cmd
}


// getKVs, get all key-values by mutableTree
func getKVs(tree *iavl.MutableTree, dataMap map[string][32]byte, wg *sync.WaitGroup) {
	tree.Iterate(func(key []byte, value []byte) bool {
		dataMap[hex.EncodeToString(key)] = sha256.Sum256(value)
		return false
	})
	wg.Done()
}

// IaviewerPrintDiff reads different key-value from leveldb according two paths
func IaviewerPrintDiff(cdc *codec.Codec, dataDir string, compareDir string, modules []string, height int) {
	for _, module := range modules {
		os.Remove(path.Join(dataDir, "/LOCK"))
		os.Remove(path.Join(compareDir, "/LOCK"))
		wmodule := wholeMoudleName(module)
		//get all key-values
		tree, err := ReadTree(dataDir, height, []byte(wmodule), DefaultCacheSize)
		if err != nil {
			log.Println("Error reading data: ", err)
			os.Exit(1)
		}
		compareTree, err := ReadTree(compareDir, height, []byte(wmodule), DefaultCacheSize)
		if err != nil {
			log.Println("Error reading compareTree data: ", err)
			os.Exit(1)
		}
		if bytes.Equal(tree.Hash(), compareTree.Hash()) {
			continue
		}

		var wg sync.WaitGroup
		wg.Add(2)
		dataMap := make(map[string][32]byte, tree.Size())
		compareDataMap := make(map[string][32]byte, compareTree.Size())
		go getKVs(tree, dataMap, &wg)
		go getKVs(compareTree, compareDataMap, &wg)
		wg.Wait()

		//get all keys
		keySize := tree.Size()
		if compareTree.Size() > keySize {
			keySize = compareTree.Size()
		}
		allKeys := make(map[string]bool, keySize)
		for k, _ := range dataMap {
			allKeys[k] = false
		}
		for k, _ := range compareDataMap {
			allKeys[k] = false
		}

		log.Println(fmt.Sprintf("==================================== %s begin ====================================", module))
		//find diff value by each key
		for key, _ := range allKeys {
			value, ok := dataMap[key]
			compareValue, compareOK := compareDataMap[key]
			keyByte, _ := hex.DecodeString(key)
			if ok && compareOK {
				if value == compareValue {
					continue
				}
				log.Println("\nvalue is different--------------------------------------------------------------------")
				log.Println("dir key-value :")
				printByKey(cdc, tree, module, keyByte)
				log.Println("compareDir key-value :")
				printByKey(cdc, compareTree, module, keyByte)
				log.Println("value is different--------------------------------------------------------------------")
				continue
			}
			if ok {
				log.Println("\nOnly be in dir--------------------------------------------------------------------")
				printByKey(cdc, tree, module, keyByte)
				continue
			}
			if compareOK {
				log.Println("\nOnly be in compare dir--------------------------------------------------------------------")
				printByKey(cdc, compareTree, module, keyByte)
				continue
			}

		}
		log.Println(fmt.Sprintf("==================================== %s end ====================================", module))
	}
}

func wholeMoudleName(moudle string) string {
	return "s/k:" + moudle + "/"
}

func printByKey(cdc *codec.Codec, tree *iavl.MutableTree, module string, key []byte) {
	_, value := tree.Get(key)
	if impl, exit := printKeysDict[module]; exit {
		impl(cdc, key, value)
	} else {
		log.Println("Not Imp for moudle %s", module)
		//printKey := parseWeaveKey(key)
		//digest := hex.EncodeToString(value)
		//log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}


func accPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	if key[0] == acc_types.AddressStoreKeyPrefix[0] {
		var acc exported.Account
		bz := value
		cdc.MustUnmarshalBinaryLengthPrefixed(bz, &acc)
		log.Println(fmt.Sprintf("address:%s; account:%s\n", hex.EncodeToString(key), acc.String()))
		return
	} else if bytes.Equal(key, acc_types.GlobalAccountNumberKey) {
		log.Println(fmt.Sprintf("%s:%s\n", string(key), hex.EncodeToString(value)))
		return
	} else {
		log.Println("Not Imp for key %s", key)
		//printKey := parseWeaveKey(key)
		//digest := hex.EncodeToString(value)
		//log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}


// ReadTree loads an iavl tree from the directory
// If version is 0, load latest, otherwise, load named version
// The prefix represents which iavl tree you want to read. The iaviwer will always set a prefix.
func ReadTree(dir string, version int, prefix []byte, cacheSize int) (*iavl.MutableTree, error) {
	db, err := OpenDB(dir)
	if err != nil {
		return nil, err
	}
	if len(prefix) != 0 {
		db = dbm.NewPrefixDB(db, prefix)
	}

	tree, err := iavl.NewMutableTree(db, cacheSize)
	if err != nil {
		return nil, err
	}
	ver, err := tree.LoadVersion(int64(version))
	log.Println(fmt.Sprintf("%s Got version: %d\n", string(prefix), ver))
	return tree, err
}


func OpenDB(dir string) (dbm.DB, error) {
	switch {
	case strings.HasSuffix(dir, ".db"):
		dir = dir[:len(dir)-3]
	case strings.HasSuffix(dir, ".db/"):
		dir = dir[:len(dir)-4]
	default:
		return nil, fmt.Errorf("database directory must end with .db")
	}
	//doesn't work on windows!
	cut := strings.LastIndex(dir, "/")
	if cut == -1 {
		return nil, fmt.Errorf("cannot cut paths on %s", dir)
	}
	name := dir[cut+1:]
	db, err := dbm.NewGoLevelDB(name, dir[:cut])
	if err != nil {
		return nil, err
	}
	return db, nil
}