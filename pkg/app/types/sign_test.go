package types

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	gravity_types "github.com/ci123chain/ci123chain/pkg/gravity/types"
	"testing"
)

func TestValidateMsgSetOrchestratorAddress(t *testing.T) {
	msg := gravity_types.MsgSendToEth{
		Sender: "0x55CA7bfdE29227D166b719Bc0FA7C0c7D2650528",
		EthDest: "0x67bdF1F80EbAC989765b3F6944f2db82130E6bcA",
		Amount: sdk.Coin{
			Denom: "WLK",
			Amount: sdk.NewInt(1000),
		},
		BridgeFee: sdk.Coin{
			Denom: "WLK",
			Amount: sdk.NewInt(50),
		},
		TokenType: 1,
	}
	//bz,_ := GetCodec().MarshalBinaryBare(msg)
	//fmt.Println(bz)
	//sEnc := base64.StdEncoding.EncodeToString(bz)
	//fmt.Println(sEnc)

	txByte, err := SignCommonTx(sdk.HexToAddress(msg.Sender), 1550, 4000, []sdk.Msg{msg}, "2cd9367f83f81341d17c6f126ea771e7c588c7fa902aca79284ec60ddc52de3a", GetCodec())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("signdata:", txByte)
	fmt.Println("signdataHex:", hex.EncodeToString(txByte))
}
