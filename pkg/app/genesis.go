package app

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
	"github.com/tanhuiya/ci123chain/pkg/account"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	genesisBalance = 100
)

// State to Unmarshal
type GenesisState struct {
	Accounts []account.Account `json:"accounts"`
}

func AppGenStateJSON(cdc *amino.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	genesisState, err := AppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = json.Marshal(genesisState)
	return
}

func AppGenState(cdc *amino.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {
	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// get genesis flag account information
	accounts := make([]account.Account, 0, len(appGenTxs))
	accountm := make(map[common.Address]struct{})

	for _, appGenTx := range appGenTxs {

		var genTx AppGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}

		addr := common.HexToAddress(genTx.Address)
		if _, ok := accountm[addr]; !ok && genTx.Address != "" {
			accounts = append(accounts, account.Account{
				Address: addr,
				Amount:  genesisBalance,
			})
			accountm[addr] = struct{}{}
		}
	}
	// create the final app state
	genesisState = GenesisState{
		Accounts: accounts,
	}
	return
}

func GetInitChainer(am account.AccountMapper) func(types.Context, abci.RequestInitChain) abci.ResponseInitChain {
	return func(ctx types.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		stateJSON := req.AppStateBytes

		var genesisState GenesisState
		err := json.Unmarshal(stateJSON, &genesisState)
		if err != nil {
			panic(err)
			// return sdk.ErrGenesisParse("").TraceCause(err, "")
		}

		for _, acc := range genesisState.Accounts {
			if _, err := am.AddBalance(ctx, acc.Address, acc.Amount); err != nil {
				panic(err)
			}
			fmt.Printf("addr=%v amount=%v\n", acc.Address.Hex(), acc.Amount)
		}
		return abci.ResponseInitChain{}
	}
}