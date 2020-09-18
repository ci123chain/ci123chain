package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type MsgUpgrade struct {
	FromAddress sdk.AccAddress	`json:"from_address"`
	Signature 	[]byte   		`json:"signature"`
	PubKey		[]byte			`json:"pub_key"`

	Type      	string   		`json:"type"`
	Height    	int64    		`json:"height"`
	Name      	string   		`json:"name"`
}

func NewMsgUpgrade(t, name string, height int64) *MsgUpgrade{
	return &MsgUpgrade{
		Type:t,
		Height:height,
		Name:name,
	}
}

func (msg *MsgUpgrade) Route() string { return RouteKey }

func (msg *MsgUpgrade) MsgType() string { return "upgrade"}

func (msg *MsgUpgrade) GetSignBytes() []byte{
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MsgUpgrade) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgUpgrade) ValidateBasic() sdk.Error{
	if len(msg.Type) == 0 {
		return ErrCheckParams(DefaultCodespace, "type is invalid")
	}
	if msg.Height < 0 {
		return ErrCheckParams(DefaultCodespace, "height is invalid")
	}
	if len(msg.Name) == 0 {
		return ErrCheckParams(DefaultCodespace, "name is invalid")
	}
	return nil
}

func (msg *MsgUpgrade) Bytes() []byte {
	bytes, err := ModuleCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgUpgrade) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}

func (msg *MsgUpgrade) SetPubKey(pubKey []byte) {
	msg.PubKey = pubKey
}

func (msg *MsgUpgrade) GetSignature() []byte {
	return msg.Signature
}