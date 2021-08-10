package baseapp

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/tendermint/tendermint/abci/server"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	cmn "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

// nolint - Mostly for testing
func (app *BaseApp) Check(tx sdk.Tx) (result sdk.Result, err error) {
	return app.runTx(runTxModeCheck, nil, tx)
}

// nolint - full tx execution
func (app *BaseApp) Simulate(txByte []byte) (result sdk.Result, err error) {
	tx, err := app.txDecoder(txByte)
	if err != nil {
		return sdk.Result{}, err
	}
	return app.runTx(runTxModeSimulate, nil, tx)
}

// nolint
func (app *BaseApp) Deliver(tx sdk.Tx) (result sdk.Result, err error) {
	return app.runTx(runTxModeDeliver, nil, tx)
}

// RunForever - BasecoinApp execution and cleanup
func RunForever(app abci.Application) {

	// Start the ABCI server
	srv, err := server.NewServer("0.0.0.0:26658", "socket", app)
	if err != nil {
		cmn.Exit(err.Error())
		return
	}
	err = srv.Start()
	if err != nil {
		cmn.Exit(err.Error())
		return
	}

	// Wait forever
	cmn.TrapSignal(log.NewNopLogger(), func() {
		// Cleanup
		err := srv.Stop()
		if err != nil {
			cmn.Exit(err.Error())
		}
	})
}

func (app *BaseApp) NewUncachedContext(isCheckTx bool, header tmtypes.Header) sdk.Context {
	return sdk.NewContext(app.cms, header, isCheckTx, app.Logger)
}

