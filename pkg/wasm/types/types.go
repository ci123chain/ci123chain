package types

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"strings"
)

var (
	EmptyAddress = sdk.HexToAddress("0x0000000000000000000000000000000000000000")
)

type CodeInfo struct {
	CodeHash    string         `json:"code_hash"`
	Creator     sdk.AccAddress  `json:"creator"`
}

func NewCodeInfo(codeHash string, creator sdk.AccAddress) CodeInfo {

	return CodeInfo{
		CodeHash: codeHash,
		Creator:  creator,
	}
}

type ContractInfo struct {
	CodeInfo 	CodeInfo		`json:"code_info"`
	Name		string			`json:"name"`
	Version     string			`json:"version"`
	Author      string			`json:"author"`
	Email       string			`json:"email"`
	Describe	string			`json:"describe"`
	InitMsg     json.RawMessage  `json:"init_msg"`
	Created     *CreatedAt        `json:"created"`
}

type CreatedAt struct {
	BlockHeight    int64  `json:"block_height"`
	TxIndex        uint64 `json:"tx_index"`
}

func (a *CreatedAt) LessThan(b *CreatedAt) bool {
	if a == nil {
		return true
	}

	if b == nil {
		return false
	}

	return a.BlockHeight < b.BlockHeight || (a.BlockHeight == b.BlockHeight && a.TxIndex < b.TxIndex )
}

func NewCreatedAt(ctx sdk.Context) *CreatedAt {
	var index uint64
	meter := ctx.GasMeter()
	if meter != nil {
		index = meter.GasConsumed()
	}
	return &CreatedAt{
		BlockHeight:ctx.BlockHeight(),
		TxIndex:index,
	}
}

func NewContractInfo(CodeHash []byte, creator sdk.AccAddress, initMsg []byte, name, version, author, email, describe string, createdAt *CreatedAt) ContractInfo {

	return ContractInfo{
		CodeInfo:	NewCodeInfo(strings.ToUpper(hex.EncodeToString(CodeHash)),creator),
		Name:   	name,
		Version:	version,
		Author: 	author,
		Email:      email,
		Describe:   describe,
		InitMsg: 	initMsg,
		Created: 	createdAt,
	}
}

//func NewParams(ctx sdk.Context, creator sdk.AccAddress, deposit sdk.Coins, contractAcct account.BaseAccount) {}


type WasmConfig struct {}

