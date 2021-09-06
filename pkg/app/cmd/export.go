package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const (
	FlagHeight           = "height"
	FlagForZeroHeight    = "for-zero-height"
	FlagJailAllowedAddrs = "jail-allowed-addrs"
)

func ExportCmd(appExporter app.AppExporter, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export state to JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := app.NewDefaultContext()
			config := serverCtx.Config

			homeDir, _ := cmd.Flags().GetString(flags.FlagHome)
			config.SetRoot(homeDir)

			if _, err := os.Stat(config.GenesisFile()); os.IsNotExist(err) {
				return err
			}

			stateDB := ""

			dbType := viper.GetString(flagCiStateDBType)
			if dbType == "" {
				dbType = "redis"
			}
			dbHost := viper.GetString(flagCiStateDBHost)
			if dbHost == "" {
				var err error
				dbHost, err = util.GetDomain()
				if err != nil {
					return err
				}
			}
			dbTls := viper.GetBool(flagCiStateDBTls)
			dbPort := viper.GetUint64(flagCiStateDBPort)
			p := strconv.FormatUint(dbPort, 10)

			switch dbType {
			case "redis":
				stateDB = "redisdb://" + dbHost + ":" + p
				if dbTls {
					stateDB += "#tls"
				}
			default:
				return errors.New(fmt.Sprintf("types of db: %s, which is not reids not implement yet", dbType))
			}

			if appExporter == nil {
				if _, err := fmt.Fprintln(os.Stderr, "WARNING: App exporter not defined. Returning genesis file."); err != nil {
					return err
				}

				genesis, err := ioutil.ReadFile(config.GenesisFile())
				if err != nil {
					return err
				}

				fmt.Println(string(genesis))
				return nil
			}

			traceWriterFile, _ := cmd.Flags().GetString(flagTraceStore)
			traceWriter, err := openTraceWriter(traceWriterFile)
			if err != nil {
				return err
			}

			height, _ := cmd.Flags().GetInt64(FlagHeight)
			forZeroHeight, _ := cmd.Flags().GetBool(FlagForZeroHeight)
			jailAllowedAddrs, _ := cmd.Flags().GetStringSlice(FlagJailAllowedAddrs)

			exported, err := appExporter(serverCtx.Logger, stateDB, traceWriter, height, forZeroHeight, jailAllowedAddrs, nil)
			if err != nil {
				return fmt.Errorf("error exporting state: %v", err)
			}

			doc, err := tmtypes.GenesisDocFromFile(serverCtx.Config.GenesisFile())
			if err != nil {
				return err
			}

			var vals = make([]tmtypes.GenesisValidator, 0)
			for _, v := range exported.Validators {
				v := tmtypes.GenesisValidator{
					Address: v.GetConsPubKey().Address(),
					PubKey:  v.GetConsPubKey(),
					Power:   v.PotentialConsensusPower(),
					Name:    v.GetMoniker(),
				}
				vals = append(vals, v)
			}

			doc.AppState = exported.AppState
			doc.Validators = vals
			doc.InitialHeight = exported.Height + 1
			if exported.ConsensusParams != nil {
				//var block tmproto.BlockParams
				//var evidence tmproto.EvidenceParams
				//var validator tmproto.ValidatorParams
				//if exported.ConsensusParams.Block != nil {
				//	block.MaxBytes = exported.ConsensusParams.Block.MaxBytes
				//	block.MaxGas = exported.ConsensusParams.Block.MaxGas
				//	block.TimeIotaMs = doc.ConsensusParams.Block.TimeIotaMs
				//}else {
				//	block.MaxGas = 0
				//	block.MaxBytes = 0
				//	block.TimeIotaMs = 0
				//}
				//if exported.ConsensusParams.Evidence != nil {
				//	evidence.MaxAgeNumBlocks = exported.ConsensusParams.Evidence.MaxAgeNumBlocks
				//	evidence.MaxAgeDuration = exported.ConsensusParams.Evidence.MaxAgeDuration
				//	evidence.MaxBytes = exported.ConsensusParams.Evidence.MaxBytes
				//}else {
				//	evidence.MaxBytes = 0
				//	evidence.MaxAgeDuration = 0
				//	evidence.MaxAgeNumBlocks = 0
				//}
				//if exported.ConsensusParams.Validator != nil {
				//	validator.PubKeyTypes = exported.ConsensusParams.Validator.PubKeyTypes
				//}else {
				//	validator.PubKeyTypes = []string{vals[0].PubKey.Type()}
				//}
				doc.ConsensusParams = &tmproto.ConsensusParams{
					Block: tmproto.BlockParams{
						//MaxBytes:   exported.ConsensusParams.Block.MaxBytes,
						//MaxGas:     exported.ConsensusParams.Block.MaxGas,
						//TimeIotaMs: doc.ConsensusParams.Block.TimeIotaMs,
					},
					Evidence: tmproto.EvidenceParams{
						//MaxAgeNumBlocks: exported.ConsensusParams.Evidence.MaxAgeNumBlocks,
						//MaxAgeDuration:  exported.ConsensusParams.Evidence.MaxAgeDuration,
						//MaxBytes:        exported.ConsensusParams.Evidence.MaxBytes,
					},
					Validator: tmproto.ValidatorParams{
						PubKeyTypes: []string{vals[0].PubKey.Type()},//exported.ConsensusParams.Validator.PubKeyTypes,
					},
				}
			}
			//doc.ConsensusParams = &tmproto.ConsensusParams{
			//	Block: tmproto.BlockParams{
			//		MaxBytes:   exported.ConsensusParams.Block.MaxBytes,
			//		MaxGas:     exported.ConsensusParams.Block.MaxGas,
			//		TimeIotaMs: doc.ConsensusParams.Block.TimeIotaMs,
			//	},
			//	Evidence: tmproto.EvidenceParams{
			//		MaxAgeNumBlocks: exported.ConsensusParams.Evidence.MaxAgeNumBlocks,
			//		MaxAgeDuration:  exported.ConsensusParams.Evidence.MaxAgeDuration,
			//		MaxBytes:        exported.ConsensusParams.Evidence.MaxBytes,
			//	},
			//	Validator: tmproto.ValidatorParams{
			//		PubKeyTypes: exported.ConsensusParams.Validator.PubKeyTypes,
			//	},
			//}

			// NOTE: Tendermint uses a custom JSON decoder for GenesisDoc
			// (except for stuff inside AppState). Inside AppState, we're free
			// to encode as protobuf or amino.
			encoded, err := tmjson.MarshalIndent(doc, "", "")
			if err != nil {
				return err
			}
			_ = ioutil.WriteFile(filepath.Join(homeDir, "logs/exportFile.json"), encoded, 0755)
			cmd.Println(string(MustSortJSON(encoded)))
			return nil
		},
	}
	cmd.SetOut(cmd.OutOrStdout())
	cmd.SetErr(cmd.ErrOrStderr())
	cmd.Flags().String(flags.FlagHome, "", "The application home directory")
	cmd.Flags().Int64(FlagHeight, -1, "Export state from a particular height (-1 means latest height)")
	cmd.Flags().Bool(FlagForZeroHeight, false, "Export state to start at height zero (perform preproccessing)")
	cmd.Flags().StringSlice(FlagJailAllowedAddrs, []string{}, "Comma-separated list of operator addresses of jailed validators to unjail")
	cmd.Flags().String(flagCiStateDBType, "redis", "database types")
	cmd.Flags().String(flagCiStateDBHost, "localhost", "db host")
	cmd.Flags().Uint64(flagCiStateDBPort, 7443, "db port")
	cmd.Flags().Bool(flagCiStateDBTls, true, "use tls")

	return cmd
}


func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile == "" {
		return
	}
	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666,
	)
}

func SortJSON(toSortJSON []byte) ([]byte, error) {
	var c interface{}
	err := json.Unmarshal(toSortJSON, &c)
	if err != nil {
		return nil, err
	}
	js, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// MustSortJSON is like SortJSON but panic if an error occurs, e.g., if
// the passed JSON isn't valid.
func MustSortJSON(toSortJSON []byte) []byte {
	js, err := SortJSON(toSortJSON)
	if err != nil {
		panic(err)
	}
	return js
}