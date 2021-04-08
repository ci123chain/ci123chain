package app

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	r "github.com/ci123chain/ci123chain/pkg/redis"
	"github.com/go-redis/redis/v8"
	sdk "github.com/tendermint/tendermint/abci/types"
	"strings"

	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"path/filepath"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
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
		//cdb, err := GetCDB(statedb)
		rdb, err := GetRDB(statedb)
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
				return nil,sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
			}
		}
		app := appFn(logger, ldb, rdb, traceStoreWriter)
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
		//cdb, err := GetCDB(statedb)
		//if err != nil {
		//	return nil, nil, abci.ErrInternal("GetCDB failed")
		//}
		//RedisDB
		rdb, err := GetRDB(statedb)
		if err != nil {
			return nil, nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "get db failed")
		}
		var traceStoreWriter io.Writer
		if traceStore != "" {
			traceStoreWriter, err = os.OpenFile(
				traceStore,
				os.O_WRONLY | os.O_APPEND | os.O_CREATE,
				0666,
				)
			if err != nil {
				return nil, nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
			}
		}
		return appFn(logger, ldb, rdb, traceStoreWriter)
		//return appFn(logger, ldb, cdb, traceStoreWriter)
	}
}

func GetRDB(stateDB string) (db dbm.DB, err error) {

	opt, err := getOption(stateDB)
	if err != nil {
		return nil, err
	}
	db = r.NewRedisDB(opt)
	err = r.DBIsValid(db.(*r.RedisDB))
	return
}

func getOption(statedb string) (*redis.Options, error) {
	// redisdb://admin:password@192.168.2.89:11001@tls
	// redisdb://192.168.2.89:11001@tls
	s := strings.Split(statedb, "://")
	if len(s) < 2 {
		return nil, errors.New(fmt.Sprintf("redisdb format error: %s", statedb))
	}
	if s[0] != "redisdb" {
		return nil, errors.New(fmt.Sprintf("redisdb format error: %s", statedb))
	}
	auths := strings.Split(s[1], "@")

	if len(auths) < 2 { // 192.168.2.89:11001#tls 无用户名 密码
		pd := strings.Split(auths[0], "#")
		opt := &redis.Options{
			Addr: pd[0],
			DB:   0,
		}
		if len(pd) == 1 {
			return opt, nil
		}
		if len(pd) == 2 {
			if pd[1] == "tls" {
				opt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
				return opt, nil
			}else {
				return nil, errors.New(fmt.Sprintf("unexpected tls setting: %s", statedb))
			}
		}else {
			return nil, errors.New(fmt.Sprintf("unexpected setting of db tls: %v", statedb))
		}
	} else { // admin:password@192.168.2.89:5984#tls
		info := auths[0] // admin:password
		userandpass := strings.Split(info, ":")
		if len(userandpass) < 2 {
			return nil, errors.New(fmt.Sprintf("unexpected setting of username and password %s", statedb))
		} else {
			pd := strings.Split(auths[1], "#")
			opt := &redis.Options{
				Addr:               pd[0],
				Username:           userandpass[0],
				Password:           userandpass[1],
				DB:                 0,
			}
			if len(pd) == 1 {
				return opt, nil
			}
			if len(pd) == 2 {
				if pd[1] == "tls" {
					opt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
					return opt, nil
				}else {
					return nil, errors.New(fmt.Sprintf("unexpected tls setting: %s", statedb))
				}
			}else {
				return nil, errors.New(fmt.Sprintf("unexpected setting of db tls: %v", statedb))
			}
		}
	}
}