package types

type Msg interface {
	Route() string
	MsgType() string
	ValidateBasic() Error
	GetFromAddress() AccAddress
	Bytes() []byte
}

type Tx interface {
	GetMsgs() []Msg
	ValidateBasic() Error
	GetSignBytes() []byte
	GetSignature() []byte
	SetSignature([]byte)
	Bytes() []byte
	SetPubKey([]byte)
	GetGas() uint64
	GetNonce() uint64
	GetFromAddress() AccAddress
}

// TxDecoder unmarshals transfer bytes
type TxDecoder func(txBytes []byte) (Tx, Error)