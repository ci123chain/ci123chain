package transfer

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var cdc = app.MakeCodec()
//off line
func SignMsgTransfer(from, to string, gas, nonce, amount uint64, priv string, isfabric bool) ([]byte, error) {
	fromAddr := sdk.HexToAddress(from)
	toAddr := sdk.HexToAddress(to)
	msg := transfer.NewMsgTransfer(fromAddr, toAddr, sdk.NewUInt64Coin(amount), isfabric)
	txByte, err := app.SignCommonTx(fromAddr, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
	//if isfabric {
	//	fab := cryptosuit.NewFabSignIdentity()
	//	pubkey, err := fab.GetPubKey(privPub)
	//	if err != nil {
	//		return nil, err
	//	}
	//	msg.SetPubKey(pubkey)
	//	signature, err = fab.Sign(msg.GetSignBytes(), privPub)
	//	if err != nil {
	//		return nil, err
	//	}
	//} else {
	//	eth := cryptosuit.NewETHSignIdentity()
	//	signature, err = eth.Sign(msg.GetSignBytes(), privPub)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
}

func NewMsgTransfer(from, to sdk.AccAddress, amount sdk.Coin, isfabric bool) []byte {
	msg := transfer.NewMsgTransfer(from, to, amount, isfabric)
	return msg.Bytes()
}

//on line
func HttpTransferTx(from, to, gas, nonce, amount, priv, proxy, reqUrl string) {
	cli := &http.Client{}
	data := url.Values{}
	data.Set("from", from)
	data.Set("to", to)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("amount", amount)
	data.Set("privateKey", priv)
	data.Set("proxy", proxy)


	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content type
	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	// request
	rep, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}