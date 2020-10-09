package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
)

type MsgRedelegate struct {
	FromAddress			 types.AccAddress	 `json:"from_address"`
	DelegatorAddress     types.AccAddress    `json:"delegator_address"`
	ValidatorSrcAddress  types.AccAddress	 `json:"validator_src_address"`
	ValidatorDstAddress  types.AccAddress	 `json:"validator_dst_address"`
	Amount               types.Coin	 		 `json:"amount"`
}

func NewMsgRedelegate(from types.AccAddress, delegateAddr types.AccAddress, validatorSrcAddr,
	validatorDstAddr types.AccAddress, amount types.Coin) *MsgRedelegate {
	//
	return &MsgRedelegate{

		DelegatorAddress:    delegateAddr,
		ValidatorSrcAddress: validatorSrcAddr,
		ValidatorDstAddress: validatorDstAddr,
		Amount:              amount,
	}
}

func (msg *MsgRedelegate) ValidateBasic() types.Error {
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("empty delegator address"))
	}
	if msg.ValidatorSrcAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("empty validator address"))
	}
	if msg.ValidatorDstAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("empty validator address"))
	}
	if !msg.Amount.Amount.IsPositive() {
		return types.ErrBadSharesAmount("bad shares amount")
	}
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", msg.FromAddress.String(), msg.DelegatorAddress.String()))
	}
	return nil
}

func (msg *MsgRedelegate) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}
func (msg *MsgRedelegate) Route() string {return RouteKey}
func (msg *MsgRedelegate) MsgType() string {return "redelegate"}
func (msg *MsgRedelegate) GetFromAddress() types.AccAddress { return msg.FromAddress}