package keeper

import (
	"crypto/ecdsa"
	"crypto/md5"
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"strings"
)

func GenerateUniqueID(b []byte) string {
	hSum := md5.Sum([]byte(b))
	hexString := hex.EncodeToString(hSum[:])
	return strings.ToUpper(hexString)
}



func getPrivateKey() ([]byte, error) {
	priKey, err := cryptoutil.DecodePriv([]byte(Priv))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, "bad privatekey")
	}
	priBz := cryptoutil.MarshalPrivateKey(priKey)
	return priBz, nil
}

func getPublicKey() ([]byte, error){
	priKey, err := cryptoutil.DecodePriv([]byte(Priv))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, "bad privatekey")
	}
	pubKey := priKey.Public().(*ecdsa.PublicKey)
	pubketBz := cryptoutil.MarshalPubkey(pubKey)
	return pubketBz, nil
}

func getBankAddress() (sdk.AccAddress, error) {
	privBz, err := getPrivateKey()
	if err != nil {
		return sdk.AccAddress{}, sdkerrors.Wrap(sdkerrors.ErrParams, "bad privatekey")
	}
	privKey, err := cryptoutil.UnMarshalPrivateKey(privBz)
	if err != nil {
		return sdk.AccAddress{},sdkerrors.Wrap(sdkerrors.ErrParams, "bad privatekey")
	}
	pubKey := privKey.Public().(*ecdsa.PublicKey)
	address, err := cryptoutil.PublicKeyToAddress(pubKey)
	if err != nil {
		return sdk.AccAddress{},sdkerrors.Wrap(sdkerrors.ErrParams, "bad pubkey")
	}
	return sdk.HexToAddress(address), nil
}