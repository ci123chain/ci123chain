package staking

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func SignRedelegateMsg(from sdk.AccAddress, amount uint64, priv string, validatorSrcAddress,
	validatorDstAddress, delegatorAddress sdk.AccAddress) (sdk.Msg, error) {
	amt := sdk.NewUInt64Coin(amount)
	msg := staking.NewRedelegateMsg(from, delegatorAddress, validatorSrcAddress, validatorDstAddress, amt)

	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err = eth.Sign(msg.GetSignBytes(), privPub)
	if err != nil {
		return nil, err
	}
	msg.SetSignature(signature)
	return msg, nil
}

func HttpRedelegateTx(from, gas, nonce, amount, priv, validatorSrcAddr, validatorDstAddr, delegatorAddr, proxy, reqUrl string) {
	cli := &http.Client{}
	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("amount", amount)
	data.Set("privateKey", priv)
	data.Set("validatorSrcAddr", validatorSrcAddr)
	data.Set("validatorDstAddr", validatorDstAddr)
	data.Set("delegatorAddr", delegatorAddr)
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