package transfer

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)

const RouteKey = "Transfer"

type MsgTransfer struct {
	FromAddress sdk.AccAddress  `json:"from"`
	To     		sdk.AccAddress  `json:"to"`
	Amount 		sdk.Coin        `json:"amount"`
	FabricMode 	bool         	`json:"fabric_mode"`
	Signature 	[]byte   		`json:"signature"`
	PubKey 	    []byte			`json:"pub_key"`
}

func NewMsgTransfer(from, to sdk.AccAddress, amount sdk.Coin, isFabric bool ) *MsgTransfer {
	msg := &MsgTransfer{
		FromAddress: 	from,
		To: 			to,
		Amount: 		amount,
		FabricMode: 	isFabric,
	}
	return msg
}

func (msg *MsgTransfer) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgTransfer) GetSignature() []byte{
	return msg.Signature
}

func (msg *MsgTransfer) ValidateBasic() sdk.Error {
	if msg.Amount.IsEqual(sdk.NewCoin(sdk.NewInt(0)))  {
		return types.ErrBadAmount(types.DefaultCodespace, errors.New("amount = 0"))
	}
	if msg.To.Empty() {
		return types.ErrBadReceiver(types.DefaultCodespace, errors.New("empty to address"))
	}
	return nil
}

func (msg *MsgTransfer) Route() string { return RouteKey }

func (msg *MsgTransfer) MsgType() string { return "transfer"}

func (msg *MsgTransfer) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MsgTransfer) Bytes() []byte {
	bytes, err := transferCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgTransfer) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}

func (msg *MsgTransfer) SetPubKey(pub []byte) {
	msg.PubKey = pub
}