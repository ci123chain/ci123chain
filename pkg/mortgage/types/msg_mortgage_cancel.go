package types

import sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"

type MsgMortgageCancel struct {
	FromAddress  		sdk.AccAddress  `json:"from_address"`
	UniqueID  			[]byte			`json:"unique_id"`
}

func (MsgMortgageCancel) Route() string {
	return RouterKey
}

func (msg MsgMortgageCancel) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if len(msg.UniqueID) < 1 {
		return sdk.ErrInternal("param mortgageRecord missing")
	}
	return nil
}
