package broadcast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type retData struct {
	Data string `json:"data"`
	RawLog  string `json:"raw_log"`
}

type ciRes struct{
	Ret 	uint32 	`json:"ret"`
	Data 	string	`json:"data"`
	Message	string	`json:"message"`
}

//同步
func httpIBCBroadcastTx(tx, reqUrl string) retData {
	cli := &http.Client{}
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
		panic(err)
	}
	var ret retData
	var ty ciRes
	err = json.Unmarshal(body, &ty)
	d := []byte(ty.Data)
	err = json.Unmarshal(d, &ret)

	if err != nil {
		fmt.Println(string(body))
	}

	if len(ret.Data) > 0 {
		fmt.Println(ret.Data)
	}
	if len(ret.RawLog) > 0 {
		fmt.Println(ret.RawLog)
	}
	return ret
}
