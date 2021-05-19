package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
)

type MsgDelegate struct {
	FromAddress			types.AccAddress	`json:"from_address"`
	DelegatorAddress  	types.AccAddress    `json:"delegator_address"`
	ValidatorAddress  	types.AccAddress    `json:"validator_address"`
	Amount            	types.Coin		  	`json:"amount"`
}

func NewMsgDelegate(from types.AccAddress, delegatorAddr types.AccAddress, validatorAddr types.AccAddress,
	amount types.Coin) *MsgDelegate {

	return &MsgDelegate{
		FromAddress: 	  from,
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Amount:           amount,
	}
}

func (msg *MsgDelegate) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidParam
	}
	if msg.ValidatorAddress.Empty() {
		return ErrInvalidParam
	}
	if !msg.Amount.Amount.IsPositive() {
		return ErrInvalidParam
	}
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return ErrInvalidParam
	}
	return nil
}

func (msg *MsgDelegate) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgDelegate) Route() string {return RouteKey}
func (msg *MsgDelegate) MsgType() string {return "delegate"}
func (msg *MsgDelegate) GetFromAddress() types.AccAddress { return msg.FromAddress}

