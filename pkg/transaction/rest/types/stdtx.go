package types
//import (
//	"encoding/json"
//	"fmt"
//	sdk "CI123Chain/pkg/abci/types"
//	"github.com/tendermint/tendermint/crypto"
//	)
//
//
//type StdTx struct {
//	Msgs       []sdk.Msg      `json:"msg" yaml:"msg"`
//	Fee        StdFee         `json:"fee" yaml:"fee"`
//	Signatures []StdSignature `json:"signatures" yaml:"signatures"`
//	Memo       string         `json:"memo" yaml:"memo"`
//}
//
//type StdFee struct {
//	Amount sdk.Coins `json:"amount" yaml:"amount"`
//	Gas    uint64    `json:"gas" yaml:"gas"`
//}
//
//// StdSignature represents a sig
//type StdSignature struct {
//	crypto.PubKey `json:"pub_key" yaml:"pub_key"` // optional
//	Signature     []byte                          `json:"signature" yaml:"signature"`
//}