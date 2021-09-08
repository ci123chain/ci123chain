package proxy

import (
	"context"
	"net/http"

	amino "github.com/tendermint/go-amino"

	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	cmn "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	//"github.com/tendermint/tendermint/rpc/collactor"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcserver "github.com/tendermint/tendermint/rpc/jsonrpc/server"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"github.com/tendermint/tendermint/types"
)

const (
	wsEndpoint = "/websocket"
)

// StartProxy will start the websocket manager on the clients,
// set up the rpc routes to proxy via the given clients,
// and start up an http/rpc server on the location given by bind (eg. :1234)
// NOTE: This function blocks - you may want to call it in a go-routine.
func StartProxy(c rpcclient.Client, listenAddr string, logger log.Logger, maxOpenConnections int) error {
	err := c.Start()
	if err != nil {
		return err
	}

	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	r := RPCRoutes(c)

	// build the handler...
	mux := http.NewServeMux()
	rpcserver.RegisterRPCFuncs(mux, r, logger)

	unsubscribeFromAllEvents := func(remoteAddr string) {
		if err := c.UnsubscribeAll(context.Background(), remoteAddr); err != nil {
			logger.Warn("Failed to unsubscribe from events", "err", err)
		}
	}
	wm := rpcserver.NewWebsocketManager(r, rpcserver.OnDisconnect(unsubscribeFromAllEvents))
	wm.SetLogger(logger)
	///*****TODO
	//collactor.SetLogger(logger)
	mux.HandleFunc(wsEndpoint, wm.WebsocketHandler)

	config := rpcserver.DefaultConfig()
	config.MaxOpenConnections = maxOpenConnections
	l, err := rpcserver.Listen(listenAddr, config)
	if err != nil {
		return err
	}
	return rpcserver.Serve(l, mux, logger, config)
}

// RPCRoutes just routes everything to the given clients, as if it were
// a tendermint fullnode.
//
// if we want security, the clients must implement it as a secure clients
func RPCRoutes(c rpcclient.Client) map[string]*rpcserver.RPCFunc {
	return map[string]*rpcserver.RPCFunc{
		// Subscribe/unsubscribe are reserved for websocket events.
		////*****TODO
		//"subscribe":       rpcserver.NewWSRPCFunc(c.(Wrapper).SubscribeWS, "query"),
		//"unsubscribe":     rpcserver.NewWSRPCFunc(c.(Wrapper).UnsubscribeWS, "query"),
		//"unsubscribe_all": rpcserver.NewWSRPCFunc(c.(Wrapper).UnsubscribeAllWS, ""),

		// info API
		"status":     rpcserver.NewRPCFunc(makeStatusFunc(c), ""),
		"blockchain": rpcserver.NewRPCFunc(makeBlockchainInfoFunc(c), "minHeight,maxHeight"),
		"genesis":    rpcserver.NewRPCFunc(makeGenesisFunc(c), ""),
		"block":      rpcserver.NewRPCFunc(makeBlockFunc(c), "height"),
		"commit":     rpcserver.NewRPCFunc(makeCommitFunc(c), "height"),
		"tx":         rpcserver.NewRPCFunc(makeTxFunc(c), "hash,prove"),
		"validators": rpcserver.NewRPCFunc(makeValidatorsFunc(c), "height"),

		// broadcast API
		"broadcast_tx_commit": rpcserver.NewRPCFunc(makeBroadcastTxCommitFunc(c), "tx"),
		"broadcast_tx_sync":   rpcserver.NewRPCFunc(makeBroadcastTxSyncFunc(c), "tx"),
		"broadcast_tx_async":  rpcserver.NewRPCFunc(makeBroadcastTxAsyncFunc(c), "tx"),

		// abci API
		"abci_query": rpcserver.NewRPCFunc(makeABCIQueryFunc(c), "path,data"),
		"abci_info":  rpcserver.NewRPCFunc(makeABCIInfoFunc(c), ""),
	}
}

func makeStatusFunc(c rpcclient.Client) func(ctx *rpctypes.Context) (*ctypes.ResultStatus, error) {
	return func(ctx *rpctypes.Context) (*ctypes.ResultStatus, error) {
		return c.Status(ctx.Context())
	}
}

func makeBlockchainInfoFunc(c rpcclient.Client) func(ctx *rpctypes.Context, minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	return func(ctx *rpctypes.Context, minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
		return c.BlockchainInfo(ctx.Context(), minHeight, maxHeight)
	}
}

func makeGenesisFunc(c rpcclient.Client) func(ctx *rpctypes.Context) (*ctypes.ResultGenesis, error) {
	return func(ctx *rpctypes.Context) (*ctypes.ResultGenesis, error) {
		return c.Genesis(ctx.Context())
	}
}

func makeBlockFunc(c rpcclient.Client) func(ctx *rpctypes.Context, height *int64) (*ctypes.ResultBlock, error) {
	return func(ctx *rpctypes.Context, height *int64) (*ctypes.ResultBlock, error) {
		return c.Block(ctx.Context(), height)
	}
}

func makeCommitFunc(c rpcclient.Client) func(ctx *rpctypes.Context, height *int64) (*ctypes.ResultCommit, error) {
	return func(ctx *rpctypes.Context, height *int64) (*ctypes.ResultCommit, error) {
		return c.Commit(ctx.Context(), height)
	}
}

func makeTxFunc(c rpcclient.Client) func(ctx *rpctypes.Context, hash []byte, prove bool) (*ctypes.ResultTx, error) {
	return func(ctx *rpctypes.Context, hash []byte, prove bool) (*ctypes.ResultTx, error) {
		return c.Tx(ctx.Context(), hash, prove)
	}
}

func makeValidatorsFunc(c rpcclient.Client) func(ctx *rpctypes.Context, height *int64) (*ctypes.ResultValidators, error) {
	return func(ctx *rpctypes.Context, height *int64) (*ctypes.ResultValidators, error) {
		///****TODO
		page := 10
		perPage := 10
		return c.Validators(ctx.Context(), height, &page, &perPage)
	}
}

func makeBroadcastTxCommitFunc(c rpcclient.Client) func(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return func(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
		return c.BroadcastTxCommit(ctx.Context(),tx)
	}
}

func makeBroadcastTxSyncFunc(c rpcclient.Client) func(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return func(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
		return c.BroadcastTxSync(ctx.Context(), tx)
	}
}

func makeBroadcastTxAsyncFunc(c rpcclient.Client) func(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return func(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
		return c.BroadcastTxAsync(ctx.Context(),tx)
	}
}

func makeABCIQueryFunc(c rpcclient.Client) func(ctx *rpctypes.Context, path string, data cmn.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return func(ctx *rpctypes.Context, path string, data cmn.HexBytes) (*ctypes.ResultABCIQuery, error) {
		return c.ABCIQuery(ctx.Context(), path, data)
	}
}

func makeABCIInfoFunc(c rpcclient.Client) func(ctx *rpctypes.Context) (*ctypes.ResultABCIInfo, error) {
	return func(ctx *rpctypes.Context) (*ctypes.ResultABCIInfo, error) {
		return c.ABCIInfo(ctx.Context())
	}
}
