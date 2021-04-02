package app

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/app/module"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)


func setup(withGenesis bool, invCheckPeriod uint) (*Chain, GenesisState) {
	db := dbm.NewMemDB()
	app := NewChain(log.NewNopLogger(), db, nil, nil)
	if withGenesis {
		return app, NewDefaultGenesisState(encCdc.Marshaler)
	}
	return app, GenesisState{}
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(isCheckTx bool) *Chain {
	app, genesisState := setup(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}


// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	return module.ModuleBasics.DefaultGenesis()
}
