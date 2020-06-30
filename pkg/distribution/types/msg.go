package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
)


type MsgFundCommunityPool struct {
	transaction.CommonTx
	Amount       sdk.Coin          `json:"amount"`
	Depositor    sdk.AccAddress    `json:"depositor"`
}


func NewMsgFundCommunityPool(from sdk.AccAddress,amount sdk.Coin, gas, nonce uint64, depositor sdk.AccAddress) MsgFundCommunityPool {
	return MsgFundCommunityPool{
		CommonTx:transaction.CommonTx{
			From:      from,
			Nonce:     nonce,
			Gas:       gas,
		},
		Amount:    amount,
		Depositor: depositor,
	}
}

// Route returns the MsgFundCommunityPool message route.
func (msg *MsgFundCommunityPool) Route() string { return RouteKey }

// GetSigners returns the signer addresses that are expected to sign the result
// of GetSignBytes.
func (msg *MsgFundCommunityPool) Bytes() []byte{
	bytes, err := DistributionCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgFundCommunityPool) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgFundCommunityPool) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

// GetSignBytes returns the raw bytes for a MsgFundCommunityPool message that
// the expected signer needs to sign.
func (msg *MsgFundCommunityPool) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	signBytes := tmsg.Bytes()
	return signBytes
}

// ValidateBasic performs basic MsgFundCommunityPool message validation.
func (msg *MsgFundCommunityPool) ValidateBasic() error {
	if !msg.Amount.IsValid() {
		return ErrInvalidCoin(DefaultCodespace, msg.Amount.String())
	}
	if msg.Depositor.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.Depositor.String())
	}

	return nil
}

func (msg *MsgFundCommunityPool) GetNonce() uint64 { return msg.Nonce}
func (msg *MsgFundCommunityPool) GetGas() uint64 { return msg.Gas}
func (msg *MsgFundCommunityPool) GetFromAddress() sdk.AccAddress { return msg.From}