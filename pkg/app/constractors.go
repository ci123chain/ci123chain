package app

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"github.com/ci123chain/ci123chain/pkg/abci"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/couchdb"
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
const DBAuthUser = "DBAuthUser"
const DBAuthPwd = "DBAuthPwd"
const DBAuth = "DBAuth"
const DBName = "DBName"
type (
	AppCreator func(home string, logger log.Logger, statedb, traceStore string) (sdk.Application, error)

	AppExporter func(home string, logger log.Logger, statedb, traceStore string) (json.RawMessage, []tmtypes.GenesisValidator, error)

	AppCreatorInit func(logger log.Logger, ldb dbm.DB, cdb dbm.DB, writer io.Writer) sdk.Application

	AppExporterInit func(logger log.Logger, ldb dbm.DB, cdb dbm.DB, writer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error)
)

func ConstructAppCreator(appFn AppCreatorInit, name string) AppCreator {

	return func(rootDir string, logger log.Logger, statedb, traceStore string) (sdk.Application, error) {
		dataDir := filepath.Join(rootDir, "data")
		ldb, err := dbm.NewGoLevelDB(name, dataDir)
		if err != nil {
			return nil, err
		}
		var sdb dbm.DB
		sdbType := strings.Split(statedb, ":")[0]
		switch sdbType {
		case "leveldb":
			sdb, err = GetLDB(statedb, name)
			if err != nil {
				return nil, err
			}
		case "couchdb":
			sdb, err = GetCDB(statedb)
			if err != nil {
				return nil, types.ErrNewDB(types.DefaultCodespace, err)
			}
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
		app := appFn(logger, ldb, sdb, traceStoreWriter)
		return app, nil
	}
}

func ConstructAppExporter(appFn AppExporterInit, name string) AppExporter {
	return func(rootDir string, logger log.Logger, statedb,traceStore string) (json.RawMessage, []tmtypes.GenesisValidator, error) {
		dataDir := filepath.Join(rootDir, "data")

		ldb, err := dbm.NewGoLevelDB(name, dataDir)
		if err != nil {
			return nil, nil, types.ErrNewDB(types.DefaultCodespace, err)
		}
		cdb, err := GetCDB(statedb)
		if err != nil {
			return nil, nil, abci.ErrInternal("GetCDB failed")
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
		return appFn(logger, ldb, cdb, traceStoreWriter)
	}
}

func GetCDB(statedb string) (db dbm.DB, err error) {
	// couchdb://admin:password@192.168.2.89:5984/dbname
	// couchdb://192.168.2.89:5984/dbname
	var dbname string
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
			dbname = DefaultDBName
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
				dbname = DefaultDBName
			} else {
				dbname = split[1]
			}
			db, err = couchdb.NewGoCouchDB(dbname, split[0],nil)
		} else {
			auth := &couchdb.BasicAuth{Username: userandpass[0], Password: userandpass[1]}
			hostandpath := auths[1]
			split := strings.Split(hostandpath, "/")
			if len(split) < 2 {
				dbname = DefaultDBName
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

func GetLDB(statedb, name string) (dbm.DB, error) {
	//leveldb:/Users/hy/sharedDB
	root := strings.Split(statedb, ":")[1]
	ldb, err := dbm.NewGoLevelDB(name, root)
	if err != nil {
		return nil, err
	}
	return ldb, nil
}