package infrastructure

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	infrastructure "github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

func SignInfrastructureStoreContent(from sdk.AccAddress, gas, nonce uint64, priv string, key, content string) ([]byte, error) {
		value := infrastructure.NewStoredContent(key, content)
		valueByte, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		msg := infrastructure.NewMsgStoreContent(from, key, valueByte)
		txByte, err := types.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
		if err != nil {
			return nil, err
		}
		return txByte, nil
}


var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
	sdk.RegisterCodec(cdc)
	cdc.RegisterConcrete(&infrastructure.MsgStoreContent{}, "ci123chain/StoreContentTx", nil)
}