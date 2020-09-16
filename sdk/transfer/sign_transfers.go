package transfer

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//off line
func SignMsgTransfer(from, to sdk.AccAddress, amount uint64, priv string, isfabric bool) (sdk.Msg, error) {

	var signature []byte
	privPub, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	msg := transfer.NewMsgTransfer(from, to, sdk.NewUInt64Coin(amount), isfabric)

	if isfabric {
		fab := cryptosuit.NewFabSignIdentity()
		pubkey, err := fab.GetPubKey(privPub)
		if err != nil {
			return nil, err
		}
		msg.SetPubKey(pubkey)
		signature, err = fab.Sign(msg.GetSignBytes(), privPub)
		if err != nil {
			return nil, err
		}
	} else {
		eth := cryptosuit.NewETHSignIdentity()
		signature, err = eth.Sign(msg.GetSignBytes(), privPub)
		if err != nil {
			return nil, err
		}
	}
	msg.SetSignature(signature)
	return msg, nil
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
