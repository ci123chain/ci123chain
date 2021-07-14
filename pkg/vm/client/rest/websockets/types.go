package websockets

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/gorilla/websocket"

	"github.com/ethereum/go-ethereum/rpc"

	rpcfilters "github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/eth/filters"
)

type SubscriptionResponseJSON struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	ID      float64     `json:"id"`
}

type SubscriptionNotification struct {
	Jsonrpc string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  *SubscriptionResult `json:"params"`
}

type SubscriptionResult struct {
	Subscription rpc.ID      `json:"subscription"`
	Result       interface{} `json:"result"`
}

type ErrorResponseJSON struct {
	Jsonrpc string            `json:"jsonrpc"`
	Error   *ErrorMessageJSON `json:"error"`
	ID      *big.Int          `json:"id"`
}

type ErrorMessageJSON struct {
	Code    *big.Int `json:"code"`
	Message string   `json:"message"`
}

type wsSubscription struct {
	sub          *rpcfilters.Subscription
	unsubscribed chan struct{} // closed when unsubscribing
	conn         *websocket.Conn
}

type JsonrpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (err *jsonError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("json-rpc error %d", err.Code)
	}
	return err.Message
}

func (err *jsonError) ErrorCode() int {
	return err.Code
}

func (err *jsonError) ErrorData() interface{} {
	return err.Data
}