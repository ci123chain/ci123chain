package staking

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func SignUndelegateMsg(from sdk.AccAddress, gas, nonce uint64, amount sdk.Coin, priv string, validatorAddress, delegatorAddress sdk.AccAddress) ([]byte, error) {
	//amt := sdk.NewUInt64Coin(amount)
	msg := staking.NewUndelegateMsg(from, validatorAddress, delegatorAddress, amount)
	txByte, err := types.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func NewUndelegateMsg(from, validatorAddress, delegatorAddress sdk.AccAddress, amt sdk.Coin) []byte {
	msg := staking.NewUndelegateMsg(from, validatorAddress, delegatorAddress, amt)
	return msg.Bytes()
}

func HttpUndelegateTx(from, gas, nonce, amount, priv, validatorAddr, delegatorAddr, proxy, reqUrl string) {
	cli := &http.Client{}
	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("amount", amount)
	data.Set("privateKey", priv)
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
