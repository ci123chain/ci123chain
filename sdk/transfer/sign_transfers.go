package transfer

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//off line
func SignTransferMsg(from, to string, amount, gas, nonce uint64, priv string, isfabric bool) ([]byte, error) {

	var signature []byte
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	privPub, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	tx := transfer.NewTransferTx(fromAddr, toAddr, gas, nonce, sdk.NewUInt64Coin(amount), isfabric)

	if isfabric {
		fab := cryptosuit.NewFabSignIdentity()
		pubkey, err := fab.GetPubKey(privPub)
		if err != nil {
			return nil, err
		}
		tx.SetPubKey(pubkey)
		signature, err = fab.Sign(tx.GetSignBytes(), privPub)
		if err != nil {
			return nil, err
		}
	} else {
		eth := cryptosuit.NewETHSignIdentity()
		signature, err = eth.Sign(tx.GetSignBytes(), privPub)
		if err != nil {
			return nil, err
		}
	}
	tx.SetSignature(signature)
	return tx.Bytes(), nil
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
