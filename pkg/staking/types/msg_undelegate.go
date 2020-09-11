package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type MsgUndelegate struct {
	FromAddress		  types.AccAddress	 `json:"from_address"`
	Signature 		  []byte  			 `json:"signature"`
	PubKey			  []byte			 `json:"pub_key"`

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

func (msg *MsgUndelegate) ValidateBasic() types.Error {
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("empty delegator address"))
	}
	if msg.ValidatorAddress.Empty() {
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

func (msg *MsgUndelegate) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	signBytes := tmsg.Bytes()
	return util.TxHash(signBytes)

}
func (msg *MsgUndelegate) SetSignature(sig []byte) {
	msg.Signature = sig
}
func (msg *MsgUndelegate) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}
func (msg *MsgUndelegate) SetPubKey(pub []byte) {
	msg.PubKey = pub
}
func (msg *MsgUndelegate) Route() string {return RouteKey}
func (msg *MsgUndelegate) MsgType() string {return "undelegate"}
func (msg *MsgUndelegate) GetFromAddress() types.AccAddress { return msg.FromAddress}
func (msg *MsgUndelegate) GetSignature() []byte {
	return msg.Signature
}