package context

import (
	"context"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/merkle"
	rpclient "github.com/tendermint/tendermint/rpc/client"
)

type Context struct {
	HomeDir 	string
	NodeURI 	string
	FromAddr 	sdk.AccAddress
	Blocked     bool
	Client 		rpclient.Client
	Verbose 	bool
	Height 		int64
	Cdc 		*codec.Codec
	Code 		int64
}

func (ctx *Context) Context() context.Context {
	return context.Background()
}

func (ctx *Context) GetNode() (rpclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("must define node URL")
	}
	return ctx.Client, nil
}

func (ctx Context) WithCode (code int64) Context {
	ctx.Code = code
	return ctx
}

func (ctx Context) WithCodec(cdc *codec.Codec) Context {
	ctx.Cdc = cdc
	return ctx
}

func (ctx Context) WithHeight(height int64) Context {
	ctx.Height = height
	return ctx
}

func (ctx Context) WithFrom(from sdk.AccAddress) Context {
	ctx.FromAddr = from
	return ctx
}

func (ctx Context) WithBlocked(blocked bool) Context {
	ctx.Blocked = blocked
	return ctx
}

func (ctx *Context) GetFromAddresses() (sdk.AccAddress) {
	return ctx.FromAddr
}

func (ctx *Context) GetBlocked() (bool) {
	return ctx.Blocked
}

func (ctx *Context) GetBalanceByAddress(addr sdk.AccAddress, isProve bool) (sdk.Coin, *merkle.Proof, error) {
	qparams := keeper.NewQueryAccountParams(addr)
	bz, err := ctx.Cdc.MarshalJSON(qparams)
	if err != nil {
		return sdk.NewEmptyCoin(), nil , err
	}
	res, _, proof, err := ctx.Query("/custom/" + types.ModuleName + "/" + types.QueryAccount, bz, isProve)
	if res == nil{
		return sdk.NewEmptyCoin(), nil, errors.New("The account does not exist")
	}
	if err != nil {
		return sdk.NewEmptyCoin(), nil, err
	}
	var acc exported.Account
	err2 := ctx.Cdc.UnmarshalBinaryLengthPrefixed(res, &acc)
	if err2 != nil {
		return sdk.NewEmptyCoin(), nil, err2
	}
	balance := acc.GetCoin()
	return balance, proof, nil
}

func (ctx *Context) GetNonceByAddress(addr sdk.AccAddress, isProve bool) (uint64, *merkle.Proof, error) {
	qparams := keeper.NewQueryAccountParams(addr)
	bz, err := ctx.Cdc.MarshalJSON(qparams)
	if err != nil {
		return 0, nil , err
	}
	res, _, proof, err := ctx.Query("/custom/" + types.ModuleName + "/" + types.QueryAccount, bz, isProve)
	if res == nil{
		return 0, nil, errors.New("The account does not exist")
	}
	if err != nil {
		return 0, nil, err
	}
	var acc exported.Account
	err2 := ctx.Cdc.UnmarshalBinaryLengthPrefixed(res, &acc)
	if err2 != nil {
		return 0, nil, err2
	}
	nonce := acc.GetSequence()
	return nonce, proof, nil
}

// PrintOutput prints output while respecting output and indent flags
// NOTE: pass in marshalled structs that have been unmarshaled
// because this function will panic on marshaling errors
func (ctx Context) PrintOutput(toPrint fmt.Stringer) (err error) {
	//var out []byte

	//switch ctx.OutputFormat {
	//case "text":
	//	out, err = yaml.Marshal(&toPrint)
	//
	//case "json":
	//	if ctx.Indent {
	//		out, err = ctx.Codec.MarshalJSONIndent(toPrint, "", "  ")
	//	} else {
	//		out, err = ctx.Cdc.MarshalJSON(toPrint)
	//	}
	//}
	//if err != nil {
	//	return
	//}

	fmt.Println(toPrint)
	return
}


func (ctx *Context) SignWithTx(tx transaction.Transaction, privKey []byte, fabricMode bool) (transaction.Transaction, error) {

	var signature []byte
	var err error

	if fabricMode {
		fab := cryptosuit.NewFabSignIdentity()
		pubkey, err := fab.GetPubKey(privKey)
		if err != nil {
			return nil, err
		}
		tx.SetPubKey(pubkey)
		signature, err = fab.Sign(tx.GetSignBytes(), privKey)
		if err != nil {
			return nil, err
		}
	} else {
		//cryptosuit.NewETHSignIdentity().Sign(tx.GetSignBytes(), addr)
		eth := cryptosuit.NewETHSignIdentity()
		signature, err = eth.Sign(tx.GetSignBytes(), privKey)
		if err != nil {
			return nil, err
		}
	}

	tx.SetSignature(signature)
	return tx, nil
}

// broadcastTx
func (ctx *Context) BroadcastSignedTx(data []byte) (sdk.TxResponse, error) {
	async := ctx.Blocked
	if async {
		return ctx.BroadcastSignedDataAsync(data)
	} else {
		return ctx.BroadcastSignedData(data)
	}
}

// 消息确认 同步
func (ctx *Context) BroadcastSignedData(data []byte) (sdk.TxResponse, error) {
	res, err := ctx.broadcastTx(data)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	if ctx.Verbose {
		fmt.Printf("txHash=%v BlockHeight=%v\n", res.TxHash, res.Height)
	}
	return res, nil
}


// 内部调用 BroadcastTxSync，同步等待checktx 消息，消息确认还是异步的
func (ctx *Context) BroadcastSignedDataAsync(data []byte) (sdk.TxResponse, error) {
	return ctx.broadcastTxSync(data)
}

// Broadcast the transfer bytes to Tendermint
func (ctx *Context) broadcastTx(tx []byte) (sdk.TxResponse, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}
	res, err := node.BroadcastTxCommit(ctx.Context(), tx)
	if err != nil {
		return sdk.NewResponseFormatBroadcastTxCommit(res), err
	}
	if res.CheckTx.Code != uint32(0) {
		return sdk.NewResponseFormatBroadcastTxCommit(res), err
	}
	if res.DeliverTx.Code != uint32(0) {
		return sdk.NewResponseFormatBroadcastTxCommit(res), err
	}

	return sdk.NewResponseFormatBroadcastTxCommit(res), nil
}

func (ctx *Context) broadcastTxSync(tx []byte) (sdk.TxResponse, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}
	res, err := node.BroadcastTxSync(ctx.Context(), tx)
	if err != nil {
		return sdk.NewResponseFormatBroadcastTx(res), err
	}

	return sdk.NewResponseFormatBroadcastTx(res), nil
}

