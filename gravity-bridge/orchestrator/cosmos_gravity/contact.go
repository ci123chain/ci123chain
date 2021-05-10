package cosmos_gravity

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Contact struct {
	rpc string
}

func NewContact(rpc string) Contact {
	rpc = strings.TrimSuffix(rpc, "/")
	return Contact{rpc: rpc}
}

type ChainStatusEnum string

const (
	MOVING ChainStatusEnum = "moving"
	SYNCING ChainStatusEnum = "syncing"
	WAITING_TO_START ChainStatusEnum = "waiting_to_start"
)

type ChainStatus struct {
	BlockHeight uint64
	Status ChainStatusEnum
}

func (c Contact) GetChainStatus() (ChainStatus, error) {
	return ChainStatus{
		BlockHeight: 0,
		Status:      "",
	}, nil
}

func (c Contact) Get(url string) ([]byte, error) {
	res, err := GetRequest(c.rpc + url)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c Contact) Post(url, contentType string, data []byte) ([]byte, error) {
	res, err := PostRequest(c.rpc + url, contentType, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetRequest defines a wrapper around an HTTP GET request with a provided URL.
// An error is returned if the request or reading the body fails.
func GetRequest(url string) ([]byte, error) {
	res, err := http.Get(url) // nolint:gosec
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = res.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}

// PostRequest defines a wrapper around an HTTP POST request with a provided URL and data.
// An error is returned if the request or reading the body fails.
func PostRequest(url string, contentType string, data []byte) ([]byte, error) {
	res, err := http.Post(url, contentType, bytes.NewBuffer(data)) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("error while sending post request: %w", err)
	}

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if err = res.Body.Close(); err != nil {
		return nil, err
	}

	return bz, nil
}