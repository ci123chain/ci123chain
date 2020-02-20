package account

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)
const (
	proxy = "lb"
)

//off line
func NewAccountOffLine() (string, string, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("Error: ", err.Error());
	}
	if key == nil {
		return "","", errors.New("the key generated is nil")
	}

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privKey := hex.EncodeToString(key.D.Bytes())
	return address, privKey, nil
}

//on line
func NewAccountOnLine() ([]byte, error) {
	cli := &http.Client{}
	reqUrl := "http://ciChain:3030/account/new"
	data := url.Values{}
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
		return nil, err
	}
	return b, nil
}