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

func SignDelegateMsg(from sdk.AccAddress, amount uint64, priv string, validatorAddress, delegatorAddress sdk.AccAddress) (sdk.Msg, error) {
	amt := sdk.NewUInt64Coin(amount)
	msg := staking.NewDelegateMsg(from, delegatorAddress, validatorAddress, amt)

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

func HttpDelegateTx(from, gas, nonce, amount,priv, validatorAddr, delegatorAddr, proxy, reqUrl string) {
	cli := &http.Client{}
	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("privateKey", priv)
	data.Set("amount", amount)
	data.Set("validatorAddr", validatorAddr)
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