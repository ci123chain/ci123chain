package types

import (
	"encoding/json"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

type CodeInfo struct {
	CodeHash    []byte          `json:"code_hash"`
	Creator     sdk.AccAddress  `json:"creator"`
}

func NewCodeInfo(codeHash []byte, creator sdk.AccAddress) CodeInfo {

	return CodeInfo{
		CodeHash: codeHash,
		Creator:  creator,
	}
}

type ContractInfo struct {

	CodeID      uint64   `json:"code_id"`
	Creator     sdk.AccAddress  `json:"creator"`
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

func NewContractInfo(CodeID uint64, creator sdk.AccAddress, initMsg []byte, label string, createdAt *CreatedAt) ContractInfo {

	return ContractInfo{
		CodeID:  CodeID,
		Creator: creator,
		Label:   label,
		InitMsg: initMsg,
		Created: createdAt,
	}
}

//func NewParams(ctx sdk.Context, creator sdk.AccAddress, deposit sdk.Coins, contractAcct account.BaseAccount) {}


type WasmConfig struct {}

