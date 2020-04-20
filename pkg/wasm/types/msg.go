package types

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
)

type StoreCodeTx struct {
	transaction.CommonTx
	Sender      sdk.AccAddress    `json:"sender"`
	WASMByteCode []byte           `json:"wasm_byte_code"`
}

func NewStoreCodeTx(from sdk.AccAddress, gas, nonce uint64, sender sdk.AccAddress, wasmCode []byte) StoreCodeTx{

	return StoreCodeTx{
		CommonTx:     transaction.CommonTx{
			From:  from,
			Gas:   gas,
			Nonce: nonce,
		},
		Sender:       sender,
		WASMByteCode: wasmCode,
	}
}

//TODO
func (msg *StoreCodeTx) ValidateBasic() sdk.Error {

	/*if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	return msg.VerifySignature(msg.GetSignBytes(), true)*/

	return nil
}

func (msg *StoreCodeTx) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	signBytes := tmsg.Bytes()
	return signBytes
}

func (msg *StoreCodeTx) SetSignature(sig []byte) {}

func (msg *StoreCodeTx) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *StoreCodeTx) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *StoreCodeTx) Route() string {
	return RouteKey
}

func (msg *StoreCodeTx)GetGas() uint64 {
	return msg.Gas
}

func (msg *StoreCodeTx)GetNonce() uint64 {
	return msg.Nonce
}

func (msg *StoreCodeTx) GetFromAddress() sdk.AccAddress {
	return msg.From
}

type InstantiateContractTx struct {
	transaction.CommonTx
	Sender       sdk.AccAddress       `json:"sender"`
	CodeHash      []byte              `json:"code_hash"`
	Label        string               `json:"label"`
	Args      json.RawMessage         `json:"args"`
}

func NewInstantiateContractTx(from sdk.AccAddress, gas, nonce uint64, codeHash []byte, sender sdk.AccAddress, label string,
	initMsg json.RawMessage) InstantiateContractTx{

		return InstantiateContractTx{
			CommonTx: transaction.CommonTx{
				From:  from,
				Gas:   gas,
				Nonce: nonce,
			},
			Sender:    sender,
			CodeHash:    codeHash,
			Label:     label,
			Args:   initMsg,
		}
}


//TODO
func (msg *InstantiateContractTx) ValidateBasic() sdk.Error {
	/*if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	return msg.VerifySignature(msg.GetSignBytes(), true)*/

	return nil
}

func (msg *InstantiateContractTx) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	signBytes := tmsg.Bytes()
	return signBytes
}

func (msg *InstantiateContractTx) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *InstantiateContractTx) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *InstantiateContractTx) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *InstantiateContractTx) Route() string {
	return RouteKey
}

func (msg *InstantiateContractTx) GetGas() uint64 {
	return msg.Gas
}

func (msg *InstantiateContractTx) GetNonce() uint64 {
	return msg.Nonce
}

func (msg *InstantiateContractTx) GetFromAddress() sdk.AccAddress {
	return msg.From
}


type ExecuteContractTx struct {
	transaction.CommonTx
	Sender           sdk.AccAddress      `json:"sender"`
	Contract         sdk.AccAddress      `json:"contract"`
	Args              json.RawMessage    `json:"args"`
}

func NewExecuteContractTx(from sdk.AccAddress, gas, nonce uint64, sender sdk.AccAddress,
	contractAddress sdk.AccAddress, msg json.RawMessage) ExecuteContractTx {

	return ExecuteContractTx{
		CommonTx:  transaction.CommonTx{
			From:      from,
			Nonce:     nonce,
			Gas:       gas,
		},
		Sender:    sender,
		Contract:  contractAddress,
		Args:       msg,
	}
}

//TODO
func (msg *ExecuteContractTx) ValidateBasic() sdk.Error {

	/*if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	return msg.VerifySignature(msg.GetSignBytes(), true)*/
	return nil
}

func (msg *ExecuteContractTx) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	signBytes := tmsg.Bytes()
	return signBytes
}

func (msg *ExecuteContractTx) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *ExecuteContractTx) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *ExecuteContractTx) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *ExecuteContractTx) Route() string {
	return RouteKey
}

func (msg *ExecuteContractTx) GetGas() uint64 {
	return msg.Gas
}

func (msg *ExecuteContractTx) GetNonce() uint64 {
	return msg.Nonce
}

func (msg *ExecuteContractTx) GetFromAddress() sdk.AccAddress {
	return msg.From
}
