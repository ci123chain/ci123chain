package types

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"strings"
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
	Label       string          `json:"label"` //标签
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

func NewContractInfo(CodeHash []byte, creator sdk.AccAddress, initMsg []byte, label string, createdAt *CreatedAt) ContractInfo {

	return ContractInfo{
		CodeInfo:	NewCodeInfo(strings.ToUpper(hex.EncodeToString(CodeHash)),creator),
		Label:   	label,
		InitMsg: 	initMsg,
		Created: 	createdAt,
	}
}

//func NewParams(ctx sdk.Context, creator sdk.AccAddress, deposit sdk.Coins, contractAcct account.BaseAccount) {}


type WasmConfig struct {}

