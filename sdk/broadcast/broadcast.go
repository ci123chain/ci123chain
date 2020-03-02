package broadcast

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//同步
func httpBroadcastTx(tx string) ([]byte, error) {
	cli := &http.Client{}
	reqUrl := "http://ciChain:3030/tx/broadcast"
	data := url.Values{}
	data.Set("txByte", tx)
	data.Set("proxy", "lb")
	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content type
	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	// request
	rep2, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(rep2.Body)
	if err != nil {
		return nil, err
	}
	return body, nil

}
