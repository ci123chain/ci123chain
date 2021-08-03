package cosmos_gravity

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/cmd/rpc"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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

type formatBlock struct {
	BlockID types.BlockID `json:"block_id"`
	Block   *struct {
		Header  struct {
			Version struct {
				Block string `protobuf:"varint,1,opt,name=block,proto3" json:"block,omitempty"`
				App   string `protobuf:"varint,2,opt,name=app,proto3" json:"app,omitempty"`
			} `json:"version"`
			ChainID string              `json:"chain_id"`
			Height  string               `json:"height"`
			Time    time.Time           `json:"time"`
			NumTxs   string             `json:"num_txs"`
			TotalTxs string             `json:"total_txs"`

			// prev block info
			LastBlockID types.BlockID `json:"last_block_id"`

			// hashes of block data
			LastCommitHash tmbytes.HexBytes `json:"last_commit_hash"` // commit from validators from the last block
			DataHash       tmbytes.HexBytes `json:"data_hash"`        // transactions

			// hashes from the app output from the prev block
			ValidatorsHash     tmbytes.HexBytes `json:"validators_hash"`      // validators for the current block
			NextValidatorsHash tmbytes.HexBytes `json:"next_validators_hash"` // validators for the next block
			ConsensusHash      tmbytes.HexBytes `json:"consensus_hash"`       // consensus params for current block
			AppHash            tmbytes.HexBytes `json:"app_hash"`             // state after txs from the previous block
			// root hash of all results from the txs from the previous block
			// see `deterministicResponseDeliverTx` to understand which parts of a tx is hashed into here
			LastResultsHash tmbytes.HexBytes `json:"last_results_hash"`

			// consensus info
			EvidenceHash    tmbytes.HexBytes `json:"evidence_hash"`    // evidence included in the block
			ProposerAddress types.Address          `json:"proposer_address"` // original proposer of the block

			Random types.VrfRandom `json:"vrf_random"`
		}   `json:"header"`
		Data  	   types.Data 		  `json:"data"`
		Evidence   types.EvidenceData `json:"evidence"`
		LastCommit *struct {
			Height     string       	 `json:"height"`
			Round      int32       		 `json:"round"`
			BlockID    types.BlockID     `json:"block_id"`
			Signatures []types.CommitSig `json:"signatures"`
		}      `json:"last_commit"`
	}  `json:"block"`
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

	var latestBlock formatBlock
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

	blockHeight, _ := strconv.ParseUint(latestBlock.Block.LastCommit.Height, 10, 64)
	return ChainStatus{
		BlockHeight: blockHeight,
		Status:      MOVING,
	}, nil
}

func (c Contact) GetNonce(address string) uint64 {
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

func (c Contact) GetBalance(address string) sdk.Coins {
	data := url.Values{}
	data.Add("address", address)
	balanceRes, _ := c.Post("/bank/balance", data)
	var res rest.Response
	var queryRes rest.QueryRes
	json.Unmarshal(balanceRes, &res)
	json.Unmarshal(res.Data, &queryRes)
	balanceList := queryRes.Value.(map[string]interface{})["balance_list"].([]interface{})

	var balance sdk.Coins
	for _, v := range balanceList {
		vMap := v.(map[string]interface{})
		amt := vMap["amount"].(string)
		x, _ :=new(big.Int).SetString(amt, 10)
		balance = append(balance, sdk.NewCoin(vMap["denom"].(string), sdk.NewIntFromBigInt(x)))
	}

	return balance
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