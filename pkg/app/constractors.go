package app

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci"
	"github.com/tanhuiya/ci123chain/pkg/app/types"
	"github.com/tanhuiya/ci123chain/pkg/couchdb"
	sdk "github.com/tendermint/tendermint/abci/types"
	//"github.com/tanhuiya/ci123chain/pkg/couchdb"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"path/filepath"
)

type (
	AppCreator func(home string, logger log.Logger, traceStore string) (sdk.Application, error)

	AppExporter func(home string, logger log.Logger, traceStore string) (json.RawMessage, []tmtypes.GenesisValidator, error)

	AppCreatorInit func(logger log.Logger, db dbm.DB, writer io.Writer) sdk.Application

	AppExporterInit func(logger log.Logger, db dbm.DB, writer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error)
)

func ConstructAppCreator(appFn AppCreatorInit, name string) AppCreator {
	return func(rootDir string, logger log.Logger, traceStore string) (sdk.Application, error) {
		//dataDir := filepath.Join(rootDir, "data")

		db, err := couchdb.NewGoCouchDB(name, "172.31.0.2", 5984, &couchdb.BasicAuth{Username: "adminuser", Password: "password"})
		//db, err := dbm.NewGoLevelDB(name, dataDir)
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