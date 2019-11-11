package types

import sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"

type IBCMsg struct {
	BankAddress sdk.AccAddress	`json:"bank_address"`
	UniqueID []byte		`json:"unique_id"`
	ObserverID []byte	`json:"observer_id"`
	Raw 	[]byte		`json:"raw"`
	State 	 string 	`json:"state"`
}



