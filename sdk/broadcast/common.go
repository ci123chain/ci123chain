package broadcast

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//异步
func httpBroadcastAsyncTx(tx, reqUrl string) {
	cli := &http.Client{}
	data := url.Values{}
	data.Set("txByte", tx)
	data.Set("proxy", "lb")
	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content types
	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	// request
	_, _ = cli.Do(req)
}

func SendTransaction(tx string, async bool, isIBC bool, requestURL string) ([]byte, retData, error) {
	//
	if async == false {
		if isIBC == false {
			res, err := httpBroadcastTx(tx, requestURL)
			if err != nil {
				return nil, retData{}, err
			}
			return res, retData{}, nil
		}else {
			retData := httpIBCBroadcastTx(tx, requestURL)
			return nil, retData, nil
		}
	}else {
		httpBroadcastAsyncTx(tx, requestURL)
		return nil, retData{}, nil
	}
}
