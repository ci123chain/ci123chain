package app

import (
	"encoding/json"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"path/filepath"
)

type (
	AppCreator func(home string, logger log.Logger, traceStore string) (abci.Application, error)

	AppExporter func(home string, logger log.Logger, traceStore string) (json.RawMessage, []tmtypes.GenesisValidator, error)

	AppCreatorInit func(logger log.Logger, db dbm.DB, writer io.Writer) abci.Application

	AppExporterInit func(logger log.Logger, db dbm.DB, writer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error)
)

func ConstructAppCreator(appFn AppCreatorInit, name string) AppCreator {
	return func(rootDir string, logger log.Logger, traceStore string) (abci.Application, error) {
		dataDir := filepath.Join(rootDir, "data")

		db, err := dbm.NewGoLevelDB(name, dataDir)
		if err != nil {
			return nil, err
		}
		var traceStoreWriter io.Writer
		if traceStore != "" {
			traceStoreWriter, err = os.OpenFile(
				traceStore,
				os.O_WRONLY | os.O_APPEND | os.O_CREATE,
				0666,
				)
			if err != nil {
				return nil, err
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
			return nil, nil, err
		}
		var traceStoreWriter io.Writer
		if traceStore != "" {
			traceStoreWriter, err = os.OpenFile(
				traceStore,
				os.O_WRONLY | os.O_APPEND | os.O_CREATE,
				0666,
				)
			if err != nil {
				return nil, nil, err
			}
		}
		return appFn(logger, db, traceStoreWriter)
	}
}