package types

import (
	"encoding/json"
	"errors"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)
type MsgInstantiateContract struct {
	FromAddress sdk.AccAddress		`json:"from_address"`
	Signature 	[]byte   			`json:"signature"`
	PubKey		[]byte				`json:"pub_key"`

	Code    	[]byte              `json:"code"`
	Name		string				`json:"name,omitempty"`
	Version     string				`json:"version,omitempty"`
	Author      string				`json:"author,omitempty"`
	Email       string				`json:"email,omitempty"`
	Describe	string				`json:"describe,omitempty"`
	Args      	json.RawMessage     `json:"args,omitempty"`
}

func NewMsgInstantiateContract(code []byte, from sdk.AccAddress, name, version, author, email, describe string,
	initMsg json.RawMessage) *MsgInstantiateContract {

		return &MsgInstantiateContract {
			FromAddress: from,
			Code:        code,
			Name:     	 name,
			Version: 	 version,
			Author: 	 author,
			Email:       email,
			Describe:    describe,
			Args:   	 initMsg,
		}
}

//TODO
func (msg *MsgInstantiateContract) ValidateBasic() sdk.Error {
	if msg.Code == nil {
		return ErrInvalidMsg(DefaultCodespace, errors.New("code is invalid"))
	}

	if msg.FromAddress.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("sender is invalid"))
	}

	return nil
}

func (msg *MsgInstantiateContract) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MsgInstantiateContract) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgInstantiateContract) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgInstantiateContract) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *MsgInstantiateContract) Route() string {
	return RouteKey
}

func (msg *MsgInstantiateContract) MsgType() string {
	return "instantiate"
}

func (msg *MsgInstantiateContract) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}

func (msg *MsgInstantiateContract) GetSignature() []byte {
	return msg.Signature
}

type MsgExecuteContract struct {
	FromAddress sdk.AccAddress		`json:"from_address"`
	Signature 	[]byte   			`json:"signature"`
	PubKey		[]byte				`json:"pub_key"`

	Contract         sdk.AccAddress      `json:"contract"`
	Args              json.RawMessage    `json:"args"`
}

func NewMsgExecuteContract(from sdk.AccAddress, contractAddress sdk.AccAddress, msg json.RawMessage) *MsgExecuteContract {

	return &MsgExecuteContract{
		FromAddress: from,
		Contract:    contractAddress,
		Args:        msg,
	}
}

//TODO
func (msg *MsgExecuteContract) ValidateBasic() sdk.Error {
	if msg.Contract.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("contractAddress is invalid"))
	}

	if msg.FromAddress.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("sender is invalid"))
	}

	if msg.Args == nil {
		return ErrInvalidMsg(DefaultCodespace, errors.New("args is invalid"))
	}
	return nil
}

func (msg *MsgExecuteContract) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MsgExecuteContract) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgExecuteContract) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgExecuteContract) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *MsgExecuteContract) Route() string {
	return RouteKey
}

func (msg *MsgExecuteContract) MsgType() string {
	return "execute"
}

func (msg *MsgExecuteContract) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}

func (msg *MsgExecuteContract) GetSignature() []byte {
	return msg.Signature
}

type MsgMigrateContract struct {
	FromAddress sdk.AccAddress		`json:"from_address"`
	Signature 	[]byte   			`json:"signature"`
	PubKey		[]byte				`json:"pub_key"`

	Code    	[]byte              `json:"code"`
	Contract	sdk.AccAddress		`json:"contract"`
	Name		string				`json:"name,omitempty"`
	Version     string				`json:"version,omitempty"`
	Author      string				`json:"author,omitempty"`
	Email       string				`json:"email,omitempty"`
	Describe	string				`json:"describe,omitempty"`
	Args      	json.RawMessage     `json:"args,omitempty"`
}

func NewMsgMigrateContract(code []byte, from sdk.AccAddress, name, version, author, email, describe string,
	contract sdk.AccAddress, initMsg json.RawMessage) *MsgMigrateContract{

	return &MsgMigrateContract{
		FromAddress: from,
		Code:        code,
		Contract:    contract,
		Name:        name,
		Version:     version,
		Author:      author,
		Email:       email,
		Describe:    describe,
		Args:        initMsg,
	}
}

func (msg *MsgMigrateContract) ValidateBasic() sdk.Error {
	if msg.Code == nil {
		return ErrInvalidMsg(DefaultCodespace, errors.New("code is invalid"))
	}

	if msg.FromAddress.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("sender is invalid"))
	}

	if msg.Contract.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("contract is invalid"))
	}
	return nil
}

func (msg *MsgMigrateContract) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MsgMigrateContract) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgMigrateContract) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgMigrateContract) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *MsgMigrateContract) Route() string {
	return RouteKey
}

func (msg *MsgMigrateContract) MsgType() string {
	return "migrate"
}

func (msg *MsgMigrateContract) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}

func (msg *MsgMigrateContract) GetSignature() []byte {
	return msg.Signature
}
