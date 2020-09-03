package types

import (
	"encoding/json"
	"errors"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/util"
)
type InstantiateContractTx struct {
	transaction.CommonTx
	Code    	[]byte              `json:"code"`
	Sender      sdk.AccAddress      `json:"sender"`
	Name		string				`json:"name,omitempty"`
	Version     string				`json:"version,omitempty"`
	Author      string				`json:"author,omitempty"`
	Email       string				`json:"email,omitempty"`
	Describe	string				`json:"describe,omitempty"`
	Args      	json.RawMessage     `json:"args,omitempty"`
}

func NewInstantiateContractTx(code []byte, from sdk.AccAddress, gas, nonce uint64, sender sdk.AccAddress, name, version, author, email, describe string,
	initMsg json.RawMessage) InstantiateContractTx{

		return InstantiateContractTx{
			CommonTx: transaction.CommonTx{
				From:  	from,
				Gas:   	gas,
				Nonce: 	nonce,
			},
			Code:       code,
			Sender:    	sender,
			Name:     	name,
			Version: 	version,
			Author: 	author,
			Email:      email,
			Describe:   describe,
			Args:   	initMsg,
		}
}


//TODO
func (msg *InstantiateContractTx) ValidateBasic() sdk.Error {
	if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}

	if msg.Code == nil {
		return ErrInvalidMsg(DefaultCodespace, errors.New("code is invalid"))
	}

	if msg.Sender.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("sender is invalid"))
	}

	return nil
}

func (msg *InstantiateContractTx) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
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

	if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}

	if msg.Contract.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("contractAddress is invalid"))
	}

	if msg.Sender.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("sender is invalid"))
	}

	if msg.Args == nil {
		return ErrInvalidMsg(DefaultCodespace, errors.New("args is invalid"))
	}
	return nil
}

func (msg *ExecuteContractTx) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
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

type MigrateContractTx struct {
	transaction.CommonTx
	Code    	[]byte              `json:"code"`
	Sender      sdk.AccAddress      `json:"sender"`
	Contract	sdk.AccAddress		`json:"contract"`
	Name		string				`json:"name,omitempty"`
	Version     string				`json:"version,omitempty"`
	Author      string				`json:"author,omitempty"`
	Email       string				`json:"email,omitempty"`
	Describe	string				`json:"describe,omitempty"`
	Args      	json.RawMessage     `json:"args,omitempty"`
}

func NewMigrateContractTx(code []byte, from sdk.AccAddress, gas, nonce uint64, sender sdk.AccAddress, name, version, author, email, describe string,
	contract sdk.AccAddress, initMsg json.RawMessage) MigrateContractTx{

	return MigrateContractTx{
		CommonTx: transaction.CommonTx{
			From:  	from,
			Gas:   	gas,
			Nonce: 	nonce,
		},
		Code:       code,
		Sender:    	sender,
		Contract: 	contract,
		Name:     	name,
		Version: 	version,
		Author: 	author,
		Email:      email,
		Describe:   describe,
		Args:   	initMsg,
	}
}

func (msg *MigrateContractTx) ValidateBasic() sdk.Error {
	if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}

	if msg.Code == nil {
		return ErrInvalidMsg(DefaultCodespace, errors.New("code is invalid"))
	}

	if msg.Sender.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("sender is invalid"))
	}

	if msg.Contract.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("contract is invalid"))
	}
	return nil
}

func (msg *MigrateContractTx) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MigrateContractTx) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MigrateContractTx) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MigrateContractTx) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *MigrateContractTx) Route() string {
	return RouteKey
}

func (msg *MigrateContractTx) GetGas() uint64 {
	return msg.Gas
}

func (msg *MigrateContractTx) GetNonce() uint64 {
	return msg.Nonce
}

func (msg *MigrateContractTx) GetFromAddress() sdk.AccAddress {
	return msg.From
}
