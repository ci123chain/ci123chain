package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
)

type MsgUndelegate struct {
	FromAddress		  types.AccAddress	 `json:"from_address"`
	DelegatorAddress  types.AccAddress   `json:"delegator_address"`
	ValidatorAddress  types.AccAddress	 `json:"validator_address"`
	Amount            types.Coin		 `json:"amount"`
}

func NewMsgUndelegate(from types.AccAddress, delegatorAddr, validatorAddr types.AccAddress,
	amount types.Coin) *MsgUndelegate {
	//
	return &MsgUndelegate{
		FromAddress: 	  from,
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Amount:           amount,
	}
}

func (msg *MsgUndelegate) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidParam(  "emtpy delegator address")
	}
	if msg.ValidatorAddress.Empty() {
		return ErrInvalidParam( "emtpy validator address")
	}
	if !msg.Amount.Amount.IsPositive() {
		return ErrInvalidParam( "amount can not be negative")
	}
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return ErrInvalidParam(  fmt.Sprintf("expected %s, got %s", msg.FromAddress.String(), msg.DelegatorAddress.String()))
	}
	return nil
}

func (msg *MsgUndelegate) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}
func (msg *MsgUndelegate) Route() string {return RouteKey}
func (msg *MsgUndelegate) MsgType() string {return "undelegate"}
func (msg *MsgUndelegate) GetFromAddress() types.AccAddress { return msg.FromAddress}