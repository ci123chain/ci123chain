package test

import (
	"fmt"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}

func Test6(t *testing.T) {
	var p = "ArZqPNdDQyzqYdBi+j7/E+j75PiGExyZsDMDVrPWakBT"
	pubStr := fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PubKeyAminoName, p)
	var pk secp256k1.PubKeySecp256k1
	cdc.UnmarshalJSON([]byte(pubStr), &pk)
	pp := crypto.PubKey(pk)
	fmt.Println(pp.Address())
	fmt.Println(pk.Address())
}
