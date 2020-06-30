package distribution

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
)

func SignFundCommunityPoolTx(from string, amount int64, gas, nonce uint64, priv string) ([]byte, error) {
	//
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	Amount := sdk.NewCoin(sdk.NewInt(amount))
	accountAddr := sdk.HexToAddress(from)
	tx := types.NewMsgFundCommunityPool(accountAddr, Amount, gas, nonce, accountAddr)

	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)
	tx.SetSignature(signature)

	return tx.Bytes(), nil
}
