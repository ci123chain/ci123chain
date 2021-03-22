package eth

import (
"fmt"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"math"
"math/big"
"strings"

"github.com/ethereum/go-ethereum/common/hexutil"
)

// BlockNumber represents decoding hex string to block values
type BlockNumber int64

const (
	// LatestBlockNumber mapping from "latest" to 0 for tm query
	LatestBlockNumber = BlockNumber(0)

	// EarliestBlockNumber mapping from "earliest" to 1 for tm query (earliest query not supported)
	EarliestBlockNumber = BlockNumber(1)
)

// NewBlockNumber creates a new BlockNumber instance.
func NewBlockNumber(n *big.Int) BlockNumber {
	return BlockNumber(n.Int64())
}

// UnmarshalJSON parses the given JSON fragment into a BlockNumber. It supports:
// - "latest", "earliest" or "pending" as string arguments
// - the block number
// Returned errors:
// - an invalid block number error when the given argument isn't a known strings
// - an out of range error when the given block number is either too little or too large
func (bn *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case "earliest":
		*bn = EarliestBlockNumber
		return nil
	case "latest":
		*bn = LatestBlockNumber
		return nil
	case "pending":
		*bn = LatestBlockNumber
		return nil
	}

	blckNum, err := hexutil.DecodeUint64(input)
	if err != nil {
		return err
	}
	if blckNum > math.MaxInt64 {
		return fmt.Errorf("blocknumber too high")
	}

	*bn = BlockNumber(blckNum)
	return nil
}

// Int64 converts block number to primitive type
func (bn BlockNumber) Int64() int64 {
	return int64(bn)
}

// TmHeight is a util function used for the Tendermint RPC clients. It returns
// nil if the block number is "latest". Otherwise, it returns the pointer of the
// int64 value of the height.
func (bn BlockNumber) TmHeight() *int64 {
	if bn == LatestBlockNumber {
		return nil
	}
	height := bn.Int64()
	return &height
}


// Account indicates the overriding fields of account during the execution of
// a message call.
// NOTE: state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type Account struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// EthHeaderFromTendermint is an util function that returns an Ethereum Header
// from a tendermint Header.
func EthHeaderFromTendermint(header tmtypes.Header) *ethtypes.Header {
	return &ethtypes.Header{
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		UncleHash:   common.Hash{},
		Coinbase:    common.Address{},
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      common.BytesToHash(header.DataHash),
		ReceiptHash: common.Hash{},
		Difficulty:  nil,
		Number:      big.NewInt(header.Height),
		Time:        uint64(header.Time.Unix()),
		Extra:       nil,
		MixDigest:   common.Hash{},
		Nonce:       ethtypes.BlockNonce{},
	}
}