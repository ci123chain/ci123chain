package types

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
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

func NewContractInfo(codeInfo CodeInfo, initMsg []byte, name, version, author, email, describe string, createdAt *CreatedAt) ContractInfo {

	return ContractInfo{
		CodeInfo:	codeInfo,
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

// go to ../keeper/keeper.go
type WasmKeeperI interface {
	Upload(ctx sdk.Context, wasmCode []byte, creator sdk.AccAddress) (codeHash []byte, err error)

	Instantiate(ctx sdk.Context, codeHash []byte, invoker sdk.AccAddress, args json.RawMessage, name, version, author, email, describe string, genesisContractAddress sdk.AccAddress, gasWanted uint64) (sdk.AccAddress, error)

	Execute(ctx sdk.Context, contractAddress sdk.AccAddress, invoker sdk.AccAddress, args json.RawMessage, gasWanted uint64) (sdk.Result, error)
}
