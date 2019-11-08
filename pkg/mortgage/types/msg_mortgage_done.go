package types

import sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"

type MsgMortgageDone struct {
	FromAddress  		sdk.AccAddress  `json:"from_address"`
	UniqueID  			[]byte			`json:"unique_id"`
}

func (MsgMortgageDone) Route() string {
	return RouterKey
}
func (msg MsgMortgageDone) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if len(msg.UniqueID) < 1 {
		return sdk.ErrInternal("param mortgageRecord missing")
	}
	return nil
}
