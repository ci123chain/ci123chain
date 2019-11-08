package types

import sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"

const (
	StateMortgaged = "StateMortgaged"
	StateSuccess = "StateSuccess"
	StateCancel = "StateCancel"
)

type MsgMortgage struct {
	FromAddress  sdk.AccAddress `json:"from_address"`
	ToAddress 	 sdk.AccAddress `json:"to_address"`
	UniqueID 	 []byte 		`json:"unique_id"`
	Coin 	 sdk.Coin			`json:"coin"`
}

func (msg MsgMortgage) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if len(msg.UniqueID) < 1 {
		return sdk.ErrInternal("param mortgageRecord missing")
	}
	if !msg.Coin.IsValid() {
		return sdk.ErrInvalidCoins("coin is invalid" + msg.Coin.String())
	}
	return nil
}

func NewMsgMortgage(from, to sdk.AccAddress, coin sdk.Coin, uniqueID []byte) *MsgMortgage {
	msg := &MsgMortgage{
		FromAddress: from,
		ToAddress: 	to,
		UniqueID: 	uniqueID,
		Coin: 		coin,
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil
	}
	return msg
}

func (MsgMortgage) Route() string {
	return RouterKey
}

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

type Mortgage struct {
	MsgMortgage

	State  string `json:"state"`
}