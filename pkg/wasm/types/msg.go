package types

import (
	"encoding/json"
	"errors"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type MsgUploadContract struct {
	FromAddress sdk.AccAddress		`json:"from_address"`
	Code    	[]byte              `json:"code"`
}

func NewMsgUploadContract(code []byte, from sdk.AccAddress) *MsgUploadContract {
	return &MsgUploadContract {
		FromAddress: from,
		Code:        code,
	}
}

//TODO
func (msg *MsgUploadContract) ValidateBasic() sdk.Error {
	if msg.Code == nil {
		return ErrInvalidMsg(DefaultCodespace, errors.New("code is invalid"))
	}

	if msg.FromAddress.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("sender is invalid"))
	}

	return nil
}


func (msg *MsgUploadContract) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgUploadContract) Route() string {
	return RouteKey
}

func (msg *MsgUploadContract) MsgType() string {
	return "upload"
}

func (msg *MsgUploadContract) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}

type MsgInstantiateContract struct {
	FromAddress sdk.AccAddress		`json:"from_address"`
	CodeHash    []byte              `json:"code"`
	Name		string				`json:"name,omitempty"`
	Version     string				`json:"version,omitempty"`
	Author      string				`json:"author,omitempty"`
	Email       string				`json:"email,omitempty"`
	Describe	string				`json:"describe,omitempty"`
	Args      	json.RawMessage     `json:"args,omitempty"`
}

func NewMsgInstantiateContract(codeHash []byte, from sdk.AccAddress, name, version, author, email, describe string,
	initMsg json.RawMessage) *MsgInstantiateContract {

		return &MsgInstantiateContract {
			FromAddress: from,
			CodeHash:    codeHash,
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
	if len(msg.CodeHash) == 0 {
		return ErrInvalidMsg(DefaultCodespace, errors.New("code is invalid"))
	}

	if msg.FromAddress.Empty() {
		return ErrInvalidMsg(DefaultCodespace, errors.New("sender is invalid"))
	}

	return nil
}

func (msg *MsgInstantiateContract) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
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

type MsgExecuteContract struct {
	FromAddress sdk.AccAddress		`json:"from_address"`
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

func (msg *MsgExecuteContract) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
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

type MsgMigrateContract struct {
	FromAddress sdk.AccAddress		`json:"from_address"`
	CodeHash    []byte              `json:"code"`
	Contract	sdk.AccAddress		`json:"contract"`
	Name		string				`json:"name,omitempty"`
	Version     string				`json:"version,omitempty"`
	Author      string				`json:"author,omitempty"`
	Email       string				`json:"email,omitempty"`
	Describe	string				`json:"describe,omitempty"`
	Args      	json.RawMessage     `json:"args,omitempty"`
}

func NewMsgMigrateContract(codeHash []byte, from sdk.AccAddress, name, version, author, email, describe string,
	contract sdk.AccAddress, initMsg json.RawMessage) *MsgMigrateContract{
	return &MsgMigrateContract{
		FromAddress: from,
		CodeHash:    codeHash,
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
	if len(msg.CodeHash) == 0 {
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

func (msg *MsgMigrateContract) Bytes() []byte {
	bytes, err := WasmCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
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