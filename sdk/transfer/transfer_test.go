package transfer

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)


func TestSignTransferMsg(t *testing.T) {

	signdata, err := SignTransferMsg("0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c",
		"0x505A74675dc9C71eF3CB5DF309256952917E801e",2, 20000, 2, "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70", false)


	assert.NoError(t, err)
	httpTransfer(hex.EncodeToString(signdata))
}
func httpTransfer(param string) {
	cli := &http.Client{}
	reqUrl := "http://localhost:3030/tx/broadcast"
	data := url.Values{}
	data.Set("txByte", param)
	data.Set("proxy", "lb")
	req2, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req2.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content type
	req2.Header.Set("Content-Type", "x-www-form-urlencoded")
	// request
	rep2, err := cli.Do(req2)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(rep2.Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(b)
}

func TestSendTransferMSg(t *testing.T) {
	from := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	to := "0x505A74675dc9C71eF3CB5DF309256952917E801e"
	amount := "2"
	gas := "20000"
	nonce := "2"
	priv := "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	proxy := "lb"
	httpSendTransferMSg(from, to, amount, gas, nonce, priv, proxy)
}

func httpSendTransferMSg(from, to, amount, gas, nonce, priv, proxy string) {
	//
	cli := &http.Client{}
	reqUrl := "http://localhost:3030/tx/transfers"
	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("amount", amount)
	data.Set("privateKey", priv)
	data.Set("proxy", proxy)


	req2, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req2.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content type
	req2.Header.Set("Content-Type", "x-www-form-urlencoded")
	// request
	rep2, err := cli.Do(req2)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(rep2.Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(b)
}