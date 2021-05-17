package cosmos_gravity

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/cmd/rpc"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ChainStatusEnum string

const (
	MOVING ChainStatusEnum = "moving"
	SYNCING ChainStatusEnum = "syncing"
	WAITING_TO_START ChainStatusEnum = "waiting_to_start"
)

type Contact struct {
	rpc string
}

func NewContact(rpc string) Contact {
	rpc = strings.TrimSuffix(rpc, "/")
	return Contact{rpc: rpc}
}

type ChainStatus struct {
	BlockHeight uint64 `json:"block_height,omitempty"`
	Status ChainStatusEnum `json:"status"`
}

func (c Contact) GetChainStatus() (ChainStatus, error) {
	//need async?
	res, err := c.Get("/syncing")
	if err != nil {
		return ChainStatus{}, err
	}

	var sync rpc.SyncingResponse
	err = json.Unmarshal(res, &sync)
	if err != nil {
		return ChainStatus{}, err
	}

	var latestBlock *coretypes.ResultBlock
	if sync.Syncing {
		return ChainStatus{
			Status:      SYNCING,
		}, nil
	} else {
		res, err := c.Get("/blocks/latest")
		if err != nil {
			return ChainStatus{}, err
		}
		var blockRes rest.Response
		err = json.Unmarshal(res, &blockRes)
		if err != nil {
			return ChainStatus{}, err
		}

		var blockBz []byte
		err = json.Unmarshal(blockRes.Data, &blockBz)
		if err != nil {
			fmt.Println(err)
			return ChainStatus{}, err
		}

		err = json.Unmarshal(blockBz, &latestBlock)
		if err != nil {
			fmt.Println(err)
			return ChainStatus{}, err
		}

		if latestBlock.Block == nil {
			return ChainStatus{
				Status:      WAITING_TO_START,
			}, nil
		}

		if latestBlock.Block.LastCommit == nil {
			return ChainStatus{}, errors.New("No commit in block?")
		}
	}

	return ChainStatus{
		BlockHeight: uint64(latestBlock.Block.LastCommit.Height),
		Status:      MOVING,
	}, nil
}

func (c Contact) GetNonce(address string) uint64{
	data := url.Values{}
	data.Add("address", address)
	nonceRes, _ := c.Post("/account/nonce", data)
	var res rest.Response
	var queryRes rest.QueryRes
	json.Unmarshal(nonceRes, &res)
	json.Unmarshal(res.Data, &queryRes)
	nonce := queryRes.Value.(map[string]interface{})["nonce"].(float64)
	return uint64(nonce)
}

func (c Contact) BroadcastTx(txBz []byte) []byte {
	data := url.Values{}
	data.Add("tx_byte", hex.EncodeToString(txBz))
	txRes := gravity_utils.Exec(func() interface{} {
		res, err := c.Post("/tx/broadcast", data)
		if err != nil {
			return err
		}
		return res
	}).Await()
	res, ok := txRes.([]byte)
	if !ok {
		fmt.Println(txRes.(error).Error())
		return nil
	}
	return res
}

func (c Contact) Get(url string) ([]byte, error) {
	res, err := GetRequest(c.rpc + url)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c Contact) Post(url string, data url.Values) ([]byte, error) {
	res, err := PostRequest(c.rpc + url, data)
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
func PostRequest(url string, data url.Values) ([]byte, error) {
	res, err := http.PostForm(url, data) // nolint:gosec
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