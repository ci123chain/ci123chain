package types

import (
	"encoding/base64"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"testing"
)

func TestValidateMsgSetOrchestratorAddress(t *testing.T) {
	msg := MsgSendToEth{
		Sender: "0x55CA7bfdE29227D166b719Bc0FA7C0c7D2650528",
		EthDest: "0x55CA7bfdE29227D166b719Bc0FA7C0c7D2650528",
		Amount: sdk.Coin{
			Denom: "wlk",
			Amount: sdk.NewInt(1000),
		},
		BridgeFee: sdk.Coin{
			Denom: "wlk",
			Amount: sdk.NewInt(5),
		},
		TokenType: 1,
	}
	bz,_ := GravityCodec.MarshalBinaryBare(msg)
	fmt.Println(bz)
	sEnc := base64.StdEncoding.EncodeToString(bz)
	fmt.Println(sEnc)
}
