package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
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

type QureyAppResponse struct {
	Code      uint32          `json:"code"`
	FormatData string		  `json:"format_data,omitempty"`
	Data      string          `json:"data,omitempty"`
	Log       string		  `json:"log,omitempty"`
	GasWanted uint64          `json:"gas_wanted,omitempty"`
	GasUsed   uint64          `json:"gas_used,omitempty"`
	Codespace string          `json:"codespace,omitempty"`
}

// TODO: In the future, more codes may be OK.
func (res Result) IsOK() bool {
	return res.Code.IsOK()
}

// ABCIMessageLogs represents a slice of ABCIMessageLog.
type ABCIMessageLogs []ABCIMessageLog

// ABCIMessageLog defines a structure containing an indexed tx ABCI message log.
type ABCIMessageLog struct {
	MsgIndex uint32 `json:"msg_index"`
	Success  bool   `json:"success"`
	Log      string `json:"log"`
	Events StringEvents `json:"events"`
}

func NewABCIMessageLog(i uint32, success bool, log string, events Events) ABCIMessageLog {
	return ABCIMessageLog{
		MsgIndex: i,
		Success: success,
		Log:      log,
		Events:   StringifyEvents(events.ToABCIEvents()),
	}
}

// String implements the fmt.Stringer interface for the ABCIMessageLogs types.
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
	Height    int64           `json:"height,omitempty"`
	TxHash    string          `json:"txhash"`
	Index     uint32          `json:"index"`
	Code      uint32          `json:"code"`
	FormatData string		  `json:"format_data,omitempty"`
	Data      string          `json:"data,omitempty"`
	RawLog    string          `json:"raw_log,omitempty"`
	Logs	  ABCIMessageLogs  `json:"logs"`
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
	return r.TxHash == "" && r.RawLog == ""
}

// NewResponseResultTx returns a TxResponse given a ResultTx from tendermint
func NewResponseResultTx(res *ctypes.ResultTx, tx Tx, timestamp string) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	parsedLogs, _ := ParseABCILogs(res.TxResult.Log)
	var code = res.TxResult.Code
	//if code == 0 {
	//	code = 1
	//}

	return TxResponse{
		TxHash:    res.Hash.String(),
		Height:    res.Height,
		Code:      code,
		Index:     res.Index,
		FormatData: string(res.TxResult.Data),
		Data:      strings.ToUpper(hex.EncodeToString(res.TxResult.Data)),
		//FormatData: string(res.TxResult.Data),
		RawLog:    res.TxResult.Log,
		Logs:  		parsedLogs,
		Info:      res.TxResult.Info,
		GasWanted: res.TxResult.GasWanted,
		GasUsed:   res.TxResult.GasUsed,
		Events:    StringifyEvents(res.TxResult.Events),
		Tx:        tx,
		Timestamp: timestamp,
	}
}

func NewResponseFormatBroadcastTxCommit(res *ctypes.ResultBroadcastTxCommit) TxResponse {
	if res == nil {
		return TxResponse{}
	}
	if !res.CheckTx.IsOK() {
		return newTxResponseCheckTx(res)
	}
	return newTxResponseDeliverTx(res)
}

func NewResponseFormatBroadcastTx(res *ctypes.ResultBroadcastTx) TxResponse {
	if res == nil {
		return TxResponse{}
	}
	parsedLogs, _ := ParseABCILogs(res.Log)
	return TxResponse{
		Code:	res.Code,
		Data:	strings.ToUpper(hex.EncodeToString(res.Data)),
		TxHash:	strings.ToUpper(hex.EncodeToString(res.Hash)),
		RawLog: res.Log,
		Logs:   parsedLogs,
	}
}

func newTxResponseCheckTx(res *ctypes.ResultBroadcastTxCommit) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	var txHash string
	if res.Hash != nil {
		txHash = res.Hash.String()
	}

	parsedLogs, _ := ParseABCILogs(res.CheckTx.Log)

	return TxResponse{
		Height:    res.Height,
		TxHash:    txHash,
		Code:      res.CheckTx.Code,
		Data:      strings.ToUpper(hex.EncodeToString(res.CheckTx.Data)),
		FormatData:   string(res.CheckTx.Data),
		RawLog:    res.CheckTx.Log,
		Logs:   	parsedLogs,
		Info:      res.CheckTx.Info,
		GasWanted: res.CheckTx.GasWanted,
		GasUsed:   res.CheckTx.GasUsed,
		Events:    StringifyEvents(res.CheckTx.Events),
		Codespace: res.CheckTx.Codespace,
	}
}


func newTxResponseDeliverTx(res *ctypes.ResultBroadcastTxCommit) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	var txHash string
	if res.Hash != nil {
		txHash = res.Hash.String()
	}

	parsedLogs, _ := ParseABCILogs(res.DeliverTx.Log)

	return TxResponse{
		Height:    res.Height,
		TxHash:    txHash,
		Code:      res.DeliverTx.Code,
		FormatData:   string(res.DeliverTx.Data),
		Data:      strings.ToUpper(hex.EncodeToString(res.DeliverTx.Data)),
		RawLog:    res.DeliverTx.Log,
		Logs:  		parsedLogs,
		Info:      res.DeliverTx.Info,
		GasWanted: res.DeliverTx.GasWanted,
		GasUsed:   res.DeliverTx.GasUsed,
		Events:    StringifyEvents(res.DeliverTx.Events),
		Codespace: res.DeliverTx.Codespace,
	}
}


func (r TxResponse) String() string {
	var sb strings.Builder
	sb.WriteString("Response:\n")

	if r.Height > 0 {
		sb.WriteString(fmt.Sprintf("  Height: %d\n", r.Height))
	}

	if r.TxHash != "" {
		sb.WriteString(fmt.Sprintf("  TxHash: %s\n", r.TxHash))
	}

	if r.Code > 0 {
		sb.WriteString(fmt.Sprintf("  Code: %d\n", r.Code))
	}

	if r.Data != "" {
		sb.WriteString(fmt.Sprintf("  Data: %s\n", r.Data))
	}

	if r.FormatData != "" {
		sb.WriteString(fmt.Sprintf("  FormatData: %s\n", r.FormatData))
	}

	if r.RawLog != "" {
		sb.WriteString(fmt.Sprintf("  Log: %s\n", r.RawLog))
	}

	if r.Logs != nil {
		sb.WriteString(fmt.Sprintf("  Logs: %s\n", r.Logs))
	}

	if r.Info != "" {
		sb.WriteString(fmt.Sprintf("  Info: %s\n", r.Info))
	}

	if r.GasWanted != 0 {
		sb.WriteString(fmt.Sprintf("  GasWanted: %d\n", r.GasWanted))
	}

	if r.GasUsed != 0 {
		sb.WriteString(fmt.Sprintf("  GasUsed: %d\n", r.GasUsed))
	}

	if r.Codespace != "" {
		sb.WriteString(fmt.Sprintf("  Codespace: %s\n", r.Codespace))
	}

	if r.Timestamp != "" {
		sb.WriteString(fmt.Sprintf("  Timestamp: %s\n", r.Timestamp))
	}

	if len(r.Events) > 0 {
		sb.WriteString(fmt.Sprintf("  Events: \n%s\n", r.Events.String()))
	}

	return strings.TrimSpace(sb.String())
}


// ParseABCILogs attempts to parse a stringified ABCI tx log into a slice of
// ABCIMessageLog types. It returns an error upon JSON decoding failure.
func ParseABCILogs(logs string) (res ABCIMessageLogs, err error) {
	err = json.Unmarshal([]byte(logs), &res)
	return res, err
}


// WrapServiceResult wraps a result from a protobuf RPC service method call in
// a Result object or error. This method takes care of marshaling the res param to
// protobuf and attaching any events on the ctx.EventManager() to the Result.
func WrapServiceResult(ctx Context, res interface{}, err error) (*Result, error) {
	if err != nil {
		return nil, err
	}

	var data []byte
	if res != nil {
		data, err = json.Marshal(res)
		if err != nil {
			return nil, err
		}
	}

	var events Events
	if evtMgr := ctx.EventManager(); evtMgr != nil {
		events = evtMgr.Events()
	}

	return &Result{
		Data:   data,
		Events: events,
	}, nil
}