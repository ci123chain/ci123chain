package context

import (
	"CI123Chain/pkg/abci/codec"
	"CI123Chain/pkg/transaction"
	"CI123Chain/pkg/util"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	rpclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Context struct {
	HomeDir 	string
	NodeURI 	string
	InputAddressed []common.Address
	Client 		rpclient.Client
	Verbose 	bool
	Height 		int64
	Cdc 		*codec.Codec
}

func (ctx *Context) GetNode() (rpclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("must define node URL")
	}
	return ctx.Client, nil
}

func (ctx Context) WithCodec(cdc *codec.Codec) Context {
	ctx.Cdc = cdc
	return ctx
}

func (ctx Context) WithHeight(height int64) Context {
	ctx.Height = height
	return ctx
}

// Broadcast the transaction bytes to Tendermint
func (ctx *Context) BroadcastTx(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return res, err
	}

	if res.CheckTx.Code != uint32(0) {
		return res, errors.Errorf("CheckTx failed: (%d) %s",
			res.CheckTx.Code, res.CheckTx.Log)
	}
	if res.DeliverTx.Code != uint32(0) {
		return res, errors.Errorf("DeliverTx failed: (%d) %s",
			res.DeliverTx.Code, res.DeliverTx.Log)
	}
	return res, err
}

func (ctx *Context) GetInputAddresses() ([]common.Address, error) {
	return ctx.InputAddressed, nil
}

func (ctx *Context) GetBalanceByAddress(addr common.Address) (uint64, error) {
	res, _, err := ctx.query("/store/main/key", addr.Bytes())

	if err != nil {
		return 0, err
	}

	balance, err := util.BytesToUint64(res)
	if err != nil && balance == 0 {
		return 0, nil
	}
	return balance, nil
}

func (ctx *Context) SignAndBroadcastTx(tx transaction.Transaction, addr common.Address) (string, error) {
	sig, err := ctx.Sign(tx.GetSignBytes(), addr)
	if err != nil {
		return "", err
	}
	tx.SetSignature(sig)
	res, err := ctx.BroadcastTx(tx.Bytes())

	if err != nil {
		return "", err
	}
	if ctx.Verbose {
		fmt.Printf("txHash=%v BlockHeight=%v\n", res.Hash.String(), res.Height)
	}
	return res.Hash.String(), nil
}

func (ctx *Context) SignTx(tx transaction.Transaction, addr common.Address) (transaction.Transaction, error) {
	sig, err := ctx.Sign(tx.GetSignBytes(), addr)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(sig)
	return tx, nil
}

func (ctx *Context) BroadcastSignedData(data []byte) (string, error) {
	res, err := ctx.BroadcastTx(data)
	if err != nil {
		return "", err
	}
	if ctx.Verbose {
		fmt.Printf("txHash=%v BlockHeight=%v\n", res.Hash.String(), res.Height)
	}
	return res.Hash.String(), nil
}