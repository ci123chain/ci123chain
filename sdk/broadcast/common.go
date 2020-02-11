package broadcast

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//异步
func httpBroadcastAsyncTx(tx string) {
	cli := &http.Client{}
	reqUrl := "http://ciChain:3030/tx/broadcast_async"
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
	_, _ = cli.Do(req)
}

func SendTransaction(tx string, async bool, isIBC bool) ([]byte, retData, error) {
	//
	if async == false {
		if isIBC == false {
			res, err := httpBroadcastTx(tx)
			if err != nil {
				return nil, retData{}, err
			}
			return res, retData{}, nil
		}else {
			retData := httpIBCBroadcastTx(tx)
			return nil, retData, nil
		}
	}else {
		httpBroadcastAsyncTx(tx)
		return nil, retData{}, nil
	}
}
