package keeper

import (
	"crypto/ecdsa"
	"crypto/md5"
	"encoding/hex"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"strings"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

func GenerateUniqueID(b []byte) string {
	hSum := md5.Sum([]byte(b))
	hexString := hex.EncodeToString(hSum[:])
	return strings.ToUpper(hexString)
}



func getPrivateKey() ([]byte, error) {
	priKey, err := cryptoutil.DecodePriv([]byte(Priv))
	if err != nil {
		return nil, transaction.ErrBadPrivkey(types.DefaultCodespace, err)
	}
	priBz := cryptoutil.MarshalPrivateKey(priKey)
	return priBz, nil
}

func getPublicKey() ([]byte, error){
	priKey, err := cryptoutil.DecodePriv([]byte(Priv))
	if err != nil {
		return nil, transaction.ErrBadPrivkey(types.DefaultCodespace, err)
	}
	pubKey := priKey.Public().(*ecdsa.PublicKey)
	pubketBz := cryptoutil.MarshalPubkey(pubKey)
	return pubketBz, nil
}

func getBankAddress() (sdk.AccAddress, error) {
	privBz, err := getPrivateKey()
	if err != nil {
		return sdk.AccAddress{}, transaction.ErrBadPrivkey(types.DefaultCodespace, err)
	}
	privKey, err := cryptoutil.UnMarshalPrivateKey(privBz)
	if err != nil {
		return sdk.AccAddress{}, transaction.ErrBadPrivkey(types.DefaultCodespace, err)
	}
	pubKey := privKey.Public().(*ecdsa.PublicKey)
	address, err := cryptoutil.PublicKeyToAddress(pubKey)
	if err != nil {
		return sdk.AccAddress{}, transaction.ErrBadPubkey(types.DefaultCodespace, err)
	}
	return sdk.HexToAddress(address), nil
}