package types

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type MsgApplyIBC struct {
	FromAddress sdk.AccAddress	`json:"from_address"`
	Signature 	[]byte   		`json:"signature"`
	PubKey		[]byte			`json:"pub_key"`

	UniqueID 	[]byte			`json:"unique_id"`
	ObserverID 	[]byte			`json:"observer_id"`
}

func NewMsgApplyIBC(from sdk.AccAddress, uniqueID, observerID []byte) *MsgApplyIBC {
	return &MsgApplyIBC{
		FromAddress: from,
		UniqueID: uniqueID,
		ObserverID: observerID,
	}
}

func (msg *MsgApplyIBC) ValidateBasic() error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	if len(msg.UniqueID) < 1 {
		return sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("UniqueID is invalid " + hex.EncodeToString(msg.UniqueID)))
	}
	if len(msg.ObserverID) < 1 {
		return sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("ObserverID is invalid " + hex.EncodeToString(msg.ObserverID)))
	}
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty from address")
	}
	return nil
	//return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}

func (msg *MsgApplyIBC) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MsgApplyIBC) SetSignature(sig []byte) {
	msg.SetSignature(sig)
}

func (msg *MsgApplyIBC) GetSignature()[]byte {
	return msg.Signature
}

func (msg *MsgApplyIBC) Bytes() []byte {
	bytes, err := IbcCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgApplyIBC) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *MsgApplyIBC) Route() string {
	return RouterKey
}

func (msg *MsgApplyIBC) MsgType() string {
	return "apply_IBC"
}

func (msg *MsgApplyIBC) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}