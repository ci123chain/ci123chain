package shard

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ci123chain/ci123chain/pkg/order"
)

var cdc = app.MakeCodec()
//off line
func SignUpgradeMsg(t, name string, height int64, from string, gas, nonce uint64, priv string) ([]byte, error){
	fromAddr := sdk.HexToAddress(from)
	msg := order.NewMsgUpgrade(t, name, height)
	txByte, err := app.SignCommonTx(fromAddr, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func NewUpgradeMsg(t, name string, height int64) []byte {
	msg := order.NewMsgUpgrade(t, name, height)
	return msg.Bytes()
}

//on line
func HttpAddShardTx(from, gas, nonce, Type, name, height, priv, proxy, reqUrl string) {
	cli := &http.Client{}
	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("type", Type)
	data.Set("name", name)
	data.Set("height", height)
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