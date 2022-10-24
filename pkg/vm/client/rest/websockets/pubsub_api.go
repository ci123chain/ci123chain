package websockets

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/tendermint/tendermint/libs/log"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ci123chain/ci123chain/pkg/client/context"

	rpctypes "github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/eth"
	rpcfilters "github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/eth/filters"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
)

// PubSubAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec
type PubSubAPI struct {
	clientCtx context.Context
	events    *rpcfilters.EventSystem
	filtersMu sync.Mutex
	filters   map[rpc.ID]*wsSubscription
	logger    log.Logger
}

// NewAPI creates an instance of the ethereum PubSub API.
func NewAPI(clientCtx context.Context) *PubSubAPI {
	return &PubSubAPI{
		clientCtx: clientCtx,
		events:    rpcfilters.NewEventSystem(clientCtx.Client),
		filters:   make(map[rpc.ID]*wsSubscription),
		logger:    log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "websocket-clients"),
	}
}

func (api *PubSubAPI) checkTendermint(clientCtx context.Context) {
	var ticker = time.NewTicker(time.Second * 15)

	select {
	case _ = <-ticker.C:
		api.events.Check(clientCtx.Context())
	}
}

func (api *PubSubAPI) subscribe(conn *websocket.Conn, params []interface{}) (rpc.ID, error) {
	method, ok := params[0].(string)
	if !ok {
		return "0", fmt.Errorf("invalid parameters")
	}

	switch method {
	case "newHeads":
		// TODO: handle extra params
		return api.subscribeNewHeads(conn)
	case "logs":
		if len(params) > 1 {
			return api.subscribeLogs(conn, params[1])
		}

		return api.subscribeLogs(conn, nil)
	case "newPendingTransactions":
		return api.subscribePendingTransactions(conn)
	case "syncing":
		return api.subscribeSyncing(conn)
	default:
		return "0", fmt.Errorf("unsupported method %s", method)
	}
}

func (api *PubSubAPI) unsubscribe(id rpc.ID) bool {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()

	if api.filters[id] == nil {
		return false
	}
	if api.filters[id].sub != nil {
		api.filters[id].sub.Unsubscribe(api.events)
	}
	close(api.filters[id].unsubscribed)
	delete(api.filters, id)
	return true
}

func (api *PubSubAPI) subscribeNewHeads(conn *websocket.Conn) (rpc.ID, error) {
	sub, _, err := api.events.SubscribeNewHeads()
	if err != nil {
		return "", fmt.Errorf("error creating block filter: %s", err.Error())
	}

	unsubscribed := make(chan struct{})
	api.filtersMu.Lock()
	api.filters[sub.ID()] = &wsSubscription{
		sub:          sub,
		conn:         conn,
		unsubscribed: unsubscribed,
	}
	api.filtersMu.Unlock()

	go func(headersCh <-chan coretypes.ResultEvent, errCh <-chan error) {
		for {
			select {
			case event := <-headersCh:
				data, _ := event.Data.(tmtypes.EventDataNewBlockHeader)
				header := rpctypes.EthHeaderFromTendermint(data.Header)

				api.filtersMu.Lock()
				if f, found := api.filters[sub.ID()]; found {
					// write to ws conn
					res := &SubscriptionNotification{
						Jsonrpc: "2.0",
						Method:  "eth_subscription",
						Params: &SubscriptionResult{
							Subscription: sub.ID(),
							Result:       header,
						},
					}

					err = f.conn.WriteJSON(res)
				}
				api.filtersMu.Unlock()
				if err != nil {
					api.logger.Warn("writing header", "Error", err.Error())
					api.unsubscribe(sub.ID())
				}
			case <-errCh:
				api.unsubscribe(sub.ID())
				return
			case <-unsubscribed:
				return
			}
		}
	}(sub.Event(), sub.Err())

	return sub.ID(), nil
}

func (api *PubSubAPI) subscribeLogs(conn *websocket.Conn, extra interface{}) (rpc.ID, error) {
	crit := filters.FilterCriteria{}

	if extra != nil {
		params, ok := extra.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid criteria")
		}

		if params["address"] != nil {
			address, ok := params["address"].(string)
			addresses, sok := params["address"].([]interface{})
			if !ok && !sok {
				return "", fmt.Errorf("invalid address; must be address or array of addresses")
			}

			if ok {
				crit.Addresses = []common.Address{common.HexToAddress(address)}
			}

			if sok {
				crit.Addresses = []common.Address{}
				for _, addr := range addresses {
					address, ok := addr.(string)
					if !ok {
						return "", fmt.Errorf("invalid address")
					}

					crit.Addresses = append(crit.Addresses, common.HexToAddress(address))
				}
			}
		}

		if params["topics"] != nil {
			topics, ok := params["topics"].([]interface{})
			if !ok {
				return "", fmt.Errorf("invalid topics")
			}

			crit.Topics = [][]common.Hash{}
			for _, topic := range topics {
				//var isSlice bool
				tss, ok := topic.([]interface{})
				if !ok {
					var tstr string
					var ok bool
					tstr, ok = topic.(string)
					if !ok {
						return "", fmt.Errorf("invalid topics")
					}

					h := common.HexToHash(tstr)
					crit.Topics = append(crit.Topics, []common.Hash{h})
				} else {
					for _, t := range tss {
						var tstr string
						var ok bool
						tstr, ok = t.(string)
						if !ok {
							return "", fmt.Errorf("invalid topics")
						}

						h := common.HexToHash(tstr)
						crit.Topics = append(crit.Topics, []common.Hash{h})
					}
				}
			}
		}
	}

	sub, _, err := api.events.SubscribeLogs(crit)
	api.logger.Info("subscribe logs:", "address", crit.Addresses)
	if err != nil {
		return rpc.ID(""), err
	}

	unsubscribed := make(chan struct{})
	api.filtersMu.Lock()
	api.filters[sub.ID()] = &wsSubscription{
		sub:          sub,
		conn:         conn,
		unsubscribed: unsubscribed,
	}
	api.filtersMu.Unlock()

	go func(ch <-chan coretypes.ResultEvent, errCh <-chan error) {
		for {
			select {
			case event := <-ch:
				dataTx, ok := event.Data.(tmtypes.EventDataTx)
				if !ok {
					err = fmt.Errorf("invalid event data %T, expected EventDataTx", event.Data)
					return
				}

				var resultData evmtypes.ResultData
				resultData, err = evmtypes.DecodeResultData(dataTx.TxResult.Result.Data)
				if err != nil {
					return
				}
				api.logger.Info("get events:", "data", resultData)

				logs := rpcfilters.FilterLogs(resultData.Logs, crit.FromBlock, crit.ToBlock, crit.Addresses, crit.Topics)
				if logs == nil {
					continue
				}
				api.logger.Info("get logs:", "logs", logs)

				api.filtersMu.Lock()
				if f, found := api.filters[sub.ID()]; found {
					// write to ws conn
					for i := range logs {
						res := &SubscriptionNotification{
							Jsonrpc: "2.0",
							Method:  "eth_subscription",
							Params: &SubscriptionResult{
								Subscription: sub.ID(),
								Result:       logs[i],
							},
						}
						err = f.conn.WriteJSON(res)
						api.logger.Info("write logs:", "logs", logs[i])
						if err != nil {
							api.logger.Warn("failed to write header", "Error", err.Error())
							api.unsubscribe(sub.ID())
							return
						}
					}
				}
				api.filtersMu.Unlock()
			case <-errCh:
				api.unsubscribe(sub.ID())
				return
			case <-unsubscribed:
				return
			}
		}
	}(sub.Event(), sub.Err())

	return sub.ID(), nil
}

func (api *PubSubAPI) subscribePendingTransactions(conn *websocket.Conn) (rpc.ID, error) {
	sub, _, err := api.events.SubscribePendingTxs()
	if err != nil {
		return "", fmt.Errorf("error creating block filter: %s", err.Error())
	}

	unsubscribed := make(chan struct{})
	api.filtersMu.Lock()
	api.filters[sub.ID()] = &wsSubscription{
		sub:          sub,
		conn:         conn,
		unsubscribed: unsubscribed,
	}
	api.filtersMu.Unlock()

	go func(txsCh <-chan coretypes.ResultEvent, errCh <-chan error) {
		for {
			select {
			case ev := <-txsCh:
				data, _ := ev.Data.(tmtypes.EventDataTx)
				txHash := common.BytesToHash(data.Tx)

				api.filtersMu.Lock()
				if f, found := api.filters[sub.ID()]; found {
					// write to ws conn
					res := &SubscriptionNotification{
						Jsonrpc: "2.0",
						Method:  "eth_subscription",
						Params: &SubscriptionResult{
							Subscription: sub.ID(),
							Result:       txHash,
						},
					}

					err = f.conn.WriteJSON(res)
				}
				api.filtersMu.Unlock()

				if err != nil {
					api.logger.Warn("failed to write header", "Error", err.Error())
					api.unsubscribe(sub.ID())
					return
				}
			case <-errCh:
				api.unsubscribe(sub.ID())
			}
		}
	}(sub.Event(), sub.Err())

	return sub.ID(), nil
}

func (api *PubSubAPI) subscribeSyncing(conn *websocket.Conn) (rpc.ID, error) {
	return "", nil
}
