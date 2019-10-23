package types

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Result is the union of ResponseFormat and ResponseCheckTx.
type Result struct {
	// Code is the response code, is stored back on the chain.
	Code CodeType

	// Codespace is the string referring to the domain of an error
	Codespace CodespaceType

	// Data is any data returned from the app.
	// Data has to be length prefixed in order to separate
	// results from multiple msgs executions
	Data []byte

	// Log contains the txs log information. NOTE: nondeterministic.
	Log string

	// GasWanted is the maximum units of work we allow this tx to perform.
	GasWanted uint64

	// GasUsed is the amount of gas actually consumed. NOTE: unimplemented
	GasUsed uint64

	// Events contains a slice of Event objects that were emitted during some
	// execution.
	Events Events
}

// TODO: In the future, more codes may be OK.
func (res Result) IsOK() bool {
	return res.Code.IsOK()
}

// ABCIMessageLogs represents a slice of ABCIMessageLog.
type ABCIMessageLogs []ABCIMessageLog

// ABCIMessageLog defines a structure containing an indexed tx ABCI message log.
type ABCIMessageLog struct {
	MsgIndex uint16 `json:"msg_index"`
	Success  bool   `json:"success"`
	Log      string `json:"log"`
}

// String implements the fmt.Stringer interface for the ABCIMessageLogs type.
func (logs ABCIMessageLogs) String() (str string) {
	if logs != nil {
		raw, err := json.Marshal(logs)
		if err == nil {
			str = string(raw)
		}
	}

	return str
}

type TxResponse struct {
	Height    int64           `json:"height"`
	TxHash    string          `json:"txhash"`
	Code      uint32          `json:"code,omitempty"`
	Data      string          `json:"data,omitempty"`
	RawLog    string          `json:"raw_log,omitempty"`
	Logs      ABCIMessageLogs `json:"logs,omitempty"`
	Info      string          `json:"info,omitempty"`
	GasWanted int64           `json:"gas_wanted,omitempty"`
	GasUsed   int64           `json:"gas_used,omitempty"`
	Events    StringEvents    `json:"events,omitempty"`
	Codespace string          `json:"codespace,omitempty"`
	Tx        Tx              `json:"tx,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
}

// Empty returns true if the response is empty
func (r TxResponse) Empty() bool {
	return r.TxHash == "" && r.Logs == nil
}

// NewResponseResultTx returns a TxResponse given a ResultTx from tendermint
func NewResponseResultTx(res *ctypes.ResultTx, tx Tx, timestamp string) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	parsedLogs, _ := ParseABCILogs(res.TxResult.Log)

	return TxResponse{
		TxHash:    res.Hash.String(),
		Height:    res.Height,
		Code:      res.TxResult.Code,
		Data:      strings.ToUpper(hex.EncodeToString(res.TxResult.Data)),
		RawLog:    res.TxResult.Log,
		Logs:      parsedLogs,
		Info:      res.TxResult.Info,
		GasWanted: res.TxResult.GasWanted,
		GasUsed:   res.TxResult.GasUsed,
		Events:    StringifyEvents(res.TxResult.Events),
		Tx:        tx,
		Timestamp: timestamp,
	}
}



// ParseABCILogs attempts to parse a stringified ABCI tx log into a slice of
// ABCIMessageLog types. It returns an error upon JSON decoding failure.
func ParseABCILogs(logs string) (res ABCIMessageLogs, err error) {
	err = json.Unmarshal([]byte(logs), &res)
	return res, err
}
