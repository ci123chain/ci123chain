package types

import "github.com/ci123chain/ci123chain/pkg/params/types"

const (
	DefaultMaxMemoCharacters 	uint64 = 256
	DefaultTxSizeCostPerByte 	uint64 = 10
	DefaultSigVerifyCostED25519 uint64 = 590
	DefaultSigVerifyCostSecp256k1 uint64 = 1000
)

// Parameter keys
var (
	KeyMaxMemoCharacters      = []byte("MaxMemoCharacters")
	KeyTxSizeCostPerByte      = []byte("TxSizeCostPerByte")
	KeySigVerifyCostED25519   = []byte("SigVerifyCostED25519")
	KeySigVerifyCostSecp256k1 = []byte("SigVerifyCostSecp256k1")
	//KeyTxSigLimit             = []byte("TxSigLimit")
)


type Params struct {
	MaxMemoCharacters		uint64 	`json:"max_memo_characters" yaml:"max_memo_characters"`
	TxSizeCostPerByte 		uint64	`json:"tx_size_cost_per_byte" yaml:"tx_size_cost_per_byte"`
	SigVerifyCostED25519 	uint64  `json:"sig_verify_cost_ed_25519" yaml:"sig_verify_cost_ed_25519"`
	SigVerifyCostSecp256k1 	uint64 	`json:"sig_verify_cost_secp_256_k_1" yaml:"sig_verify_cost_secp_256_k_1"`
}

func DefaultParams() Params {
	return Params{
		MaxMemoCharacters: DefaultMaxMemoCharacters,
		TxSizeCostPerByte: DefaultTxSizeCostPerByte,
		SigVerifyCostED25519: DefaultSigVerifyCostED25519,
		SigVerifyCostSecp256k1: DefaultSigVerifyCostSecp256k1,
	}
}

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{Params: params}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams())
}

func ParamKeyTable() types.KeyTable {
	return types.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() types.ParamSetPairs {
	return types.ParamSetPairs{
		{KeyMaxMemoCharacters, &p.MaxMemoCharacters, nil},
		{KeyTxSizeCostPerByte, &p.TxSizeCostPerByte, nil},
		{KeySigVerifyCostED25519,  &p.SigVerifyCostED25519, nil},
		{KeySigVerifyCostSecp256k1, &p.SigVerifyCostSecp256k1, nil},
	}
}