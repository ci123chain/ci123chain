package app

import (
	"encoding/json"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/abci"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/couchdb"
	"github.com/spf13/viper"
	sdk "github.com/tendermint/tendermint/abci/types"
	"strings"

	//"github.com/ci123chain/ci123chain/pkg/couchdb"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"path/filepath"
)

const DefaultDBName = "ci123"
const DefaultDataDBName = "ci123_data"
const DBAuthUser = "DBAuthUser"
const DBAuthPwd = "DBAuthPwd"
const DBAuth = "DBAuth"
const DBName = "DBName"
type (
	AppCreator func(home string, logger log.Logger, statedb, commitDB, traceStore string) (sdk.Application, error)

	AppExporter func(home string, logger log.Logger, traceStore string) (json.RawMessage, []tmtypes.GenesisValidator, error)

	AppCreatorInit func(logger log.Logger, db, commitDB dbm.DB, writer io.Writer) sdk.Application

	AppExporterInit func(logger log.Logger, db, commitDB dbm.DB, writer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error)
)

func ConstructAppCreator(appFn AppCreatorInit, name string) AppCreator {

	return func(rootDir string, logger log.Logger, statedb, commitDB, traceStore string) (sdk.Application, error) {
		dataDir := filepath.Join(rootDir, "data")
		db, err := GetStateDB(dataDir, statedb, true)
		if err != nil {
			return nil, types.ErrNewDB(types.DefaultCodespace, err)
		}
		cdb, cerr := GetStateDB(dataDir, commitDB, false)
		if cerr != nil {
			return nil, types.ErrNewDB(types.DefaultCodespace, cerr)
		}
		var traceStoreWriter io.Writer
		if traceStore != "" {
			traceStoreWriter, err = os.OpenFile(
				traceStore,
				os.O_WRONLY | os.O_APPEND | os.O_CREATE,
				0666,
				)
			if err != nil {
				return nil, abci.ErrInternal("Open file failed")
			}
		}
		app := appFn(logger, db, cdb, traceStoreWriter)
		return app, nil
	}
}

func ConstructAppExporter(appFn AppExporterInit, name string) AppExporter {
	return func(rootDir string, logger log.Logger, traceStore string) (json.RawMessage, []tmtypes.GenesisValidator, error) {
		dataDir := filepath.Join(rootDir, "data")

		db, err := dbm.NewGoLevelDB(name, dataDir)
		if err != nil {
			return nil, nil, types.ErrNewDB(types.DefaultCodespace, err)
		}
		cdb, cerr := dbm.NewGoLevelDB(name, dataDir)
		if cerr != nil {
			return nil, nil, types.ErrNewDB(types.DefaultCodespace, cerr)
		}
		var traceStoreWriter io.Writer
		if traceStore != "" {
			traceStoreWriter, err = os.OpenFile(
				traceStore,
				os.O_WRONLY | os.O_APPEND | os.O_CREATE,
				0666,
				)
			if err != nil {
				return nil, nil, abci.ErrInternal("Open file failed")
			}
		}
		return appFn(logger, db, cdb, traceStoreWriter)
	}
}

func GetStateDB(path, statedb string, isDataDB bool) (db dbm.DB, err error) {

	var DeFaultName string
	if isDataDB {
		DeFaultName = DefaultDataDBName
	}else {
		DeFaultName = DefaultDBName
	}

	var dbname string
	if statedb == "leveldb" {
		db, err = dbm.NewGoLevelDB(DeFaultName, path)
		return
	} else {
		// couchdb://admin:password@192.168.2.89:5984/dbname
		// couchdb://192.168.2.89:5984/dbname
		s := strings.Split(statedb, "://")
		if len(s) < 2 {
			return nil, errors.New("statedb format error")
		}
		if s[0] != "couchdb" {
			return nil, errors.New("statedb format error")
		}
		auths := strings.Split(s[1], "@")

		if len(auths) < 2 { // 192.168.2.89:5984/dbname 无用户名 密码
			info := auths[0]
			split := strings.Split(info, "/")
			if len(split) < 2 {
				dbname = DeFaultName
			} else {
				dbname = split[1]
			}
			db, err = couchdb.NewGoCouchDB(dbname, split[0],nil)
			viper.Set(DBName, dbname)
		} else { // admin:password@192.168.2.89:5984/dbname
			info := auths[0] // admin:password
			userandpass := strings.Split(info, ":")
			if len(userandpass) < 2 {
				hostandpath := auths[1]
				split := strings.Split(hostandpath, "/")
				if len(split) < 2 {
					dbname = DeFaultName
				} else {
					dbname = split[1]
				}
				db, err = couchdb.NewGoCouchDB(dbname, split[0],nil)
			} else {
				auth := &couchdb.BasicAuth{Username: userandpass[0], Password: userandpass[1]}
				hostandpath := auths[1]
				split := strings.Split(hostandpath, "/")
				if len(split) < 2 {
					dbname = DeFaultName
				} else {
					dbname = split[1]
				}
				db, err = couchdb.NewGoCouchDB(dbname, split[0], auth)
				viper.Set(DBAuthUser, userandpass[0])
				viper.Set(DBAuthPwd, userandpass[1])
				viper.Set(DBName, dbname)
			}
		}
		return
	}
}