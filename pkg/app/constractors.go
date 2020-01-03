package app

import (
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/abci"
	"github.com/tanhuiya/ci123chain/pkg/app/types"
	"github.com/tanhuiya/ci123chain/pkg/couchdb"
	sdk "github.com/tendermint/tendermint/abci/types"
	"strings"

	//"github.com/tanhuiya/ci123chain/pkg/couchdb"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"path/filepath"
)

type (
	AppCreator func(home string, logger log.Logger, statedb, traceStore string) (sdk.Application, error)

	AppExporter func(home string, logger log.Logger, traceStore string) (json.RawMessage, []tmtypes.GenesisValidator, error)

	AppCreatorInit func(logger log.Logger, db dbm.DB, writer io.Writer) sdk.Application

	AppExporterInit func(logger log.Logger, db dbm.DB, writer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error)
)

func ConstructAppCreator(appFn AppCreatorInit, name string) AppCreator {

	return func(rootDir string, logger log.Logger, statedb, traceStore string) (sdk.Application, error) {
		dataDir := filepath.Join(rootDir, "data")
		db, err := getStateDB(name, dataDir, statedb)
		if err != nil {
			return nil, types.ErrNewDB(types.DefaultCodespace, err)
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
		app := appFn(logger, db, traceStoreWriter)
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
		return appFn(logger, db, traceStoreWriter)
	}
}

func getStateDB(name, path, statedb string) (db dbm.DB, err error) {
	if statedb == "leveldb" {
		db, err = dbm.NewGoLevelDB(name, path)
		return
	} else {
		// couchdb://admin:password@192.168.2.89:5984
		s := strings.Split(statedb, "://")
		if len(s) < 2 {
			return nil, errors.New("statedb format error")
		}
		if s[0] != "couchdb" {
			return nil, errors.New("statedb format error")
		}
		auths := strings.Split(s[1], "@")

		if len(auths) < 2 {
			db, err = couchdb.NewGoCouchDB(name, auths[0],nil)
		} else {
			info := auths[0]
			userpass := strings.Split(info, ":")
			if len(userpass) < 2 {
				db, err = couchdb.NewGoCouchDB(name, auths[1],nil)
			}
			auth := &couchdb.BasicAuth{Username: userpass[0], Password: userpass[1]}
			db, err = couchdb.NewGoCouchDB(name, auths[1], auth)
		}
		return
	}
	return nil, errors.New("statedb format error")
}