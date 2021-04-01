package evmtypes

import (
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Supported endpoints
const (
	QueryProtocolVersion = "protocolVersion"
	QueryBalance         = "balance"
	QueryBlockNumber     = "blockNumber"
	QueryStorage         = "storage"
	QueryCode            = "code"
	QueryNonce           = "nonce"
	QueryHashToHeight    = "hashToHeight"
	QueryTransactionLogs = "transactionLogs"
	QueryBloom           = "bloom"
	QueryLogs            = "logs"
	QueryAccount         = "account"
	QueryExportAccount   = "exportAccount"
)

// QueryResProtocolVersion is response types for protocol version query
type QueryResProtocolVersion struct {
	Version string `json:"version"`
}

func (q QueryResProtocolVersion) String() string {
	return q.Version
}

// QueryResBalance is response types for balance query
type QueryResBalance struct {
	Balance string `json:"balance"`
}

func (q QueryResBalance) String() string {
	return q.Balance
}

// QueryResBlockNumber is response types for block number query
type QueryResBlockNumber struct {
	Number int64 `json:"blockNumber"`
}

func (q QueryResBlockNumber) String() string {
	return fmt.Sprint(q.Number)
}

// QueryResStorage is response types for storage query
type QueryResStorage struct {
	Value []byte `json:"value"`
}

func (q QueryResStorage) String() string {
	return string(q.Value)
}

// QueryResCode is response types for code query
type QueryResCode struct {
	Code []byte
}

func (q QueryResCode) String() string {
	return string(q.Code)
}

// QueryResNonce is response types for Nonce query
type QueryResNonce struct {
	Nonce uint64 `json:"nonce"`
}

func (q QueryResNonce) String() string {
	return fmt.Sprint(q.Nonce)
}

// QueryETHLogs is response types for tx logs query
type QueryETHLogs struct {
	Logs []*ethtypes.Log `json:"logs"`
}

func (q QueryETHLogs) String() string {
	return fmt.Sprintf("%+v", q.Logs)
}

// QueryBloomFilter is response types for tx logs query
type QueryBloomFilter struct {
	Bloom ethtypes.Bloom `json:"bloom"`
}

func (q QueryBloomFilter) String() string {
	return string(q.Bloom.Bytes())
}

// QueryAccount is response types for querying Ethereum state objects
type QueryResAccount struct {
	Balance  string `json:"balance"`
	CodeHash []byte `json:"codeHash"`
	Nonce    uint64 `json:"nonce"`
}

type QueryResExportAccount = GenesisAccount