package validator

import (
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func NewValidatorKey() (validatorKey, pubKey string, err error) {
	var valKey secp256k1.PrivKeySecp256k1
	validator := secp256k1.GenPrivKey()
	cdc := amino.NewCodec()
	keyByte, err := cdc.MarshalJSON(validator)
	if err != nil {
		return "","", err
	}
	validatorKey = string(keyByte[1:len(keyByte)-1])
	privStr := fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PrivKeyAminoName, validatorKey)
	cdc = app.MakeCodec()
	err = cdc.UnmarshalJSON([]byte(privStr), &valKey)
	if err != nil {
		return "","", err
	}
	pubKey = hex.EncodeToString(cdc.MustMarshalJSON(valKey.PubKey()))
	return
}