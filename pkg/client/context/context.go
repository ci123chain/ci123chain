package context

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/util"
	rpclient "github.com/tendermint/tendermint/rpc/client"
)

type Context struct {
	HomeDir 	string
	NodeURI 	string
	InputAddressed []types.AccAddress
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

func (ctx *Context) GetInputAddresses() ([]types.AccAddress, error) {
	return ctx.InputAddressed, nil
}

func (ctx *Context) GetBalanceByAddress(addr types.AccAddress) (uint64, error) {
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


func (ctx *Context) SignAndBroadcastTx(tx transaction.Transaction, addr types.AccAddress) (types.TxResponse, error) {
	sig, err := ctx.Sign(tx.GetSignBytes(), addr)
	if err != nil {
		return types.TxResponse{}, err
	}
	tx.SetSignature(sig)
	res, err := ctx.BroadcastTx(tx.Bytes())

	if err != nil {
		return res, err
	}
	if ctx.Verbose {
		fmt.Printf("txHash=%v BlockHeight=%v\n", res.TxHash, res.Height)
	}
	return res, nil
}

//func (ctx *Context) SignTx(tx transaction.Transaction, addr types.AccAddress) (transaction.Transaction, error) {
//	sig, err := ctx.Sign(tx.GetSignBytes(), addr)
//	if err != nil {
//		return nil, err
//	}
//	tx.SetSignature(sig)
//	return tx, nil
//}

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

func (ctx *Context) BroadcastSignedData(data []byte) (types.TxResponse, error) {
	res, err := ctx.BroadcastTx(data)
	if err != nil {
		return types.TxResponse{}, err
	}
	if ctx.Verbose {
		fmt.Printf("txHash=%v BlockHeight=%v\n", res.TxHash, res.Height)
	}
	return res, nil
}


//func (ctx *Context) SignTx2(tx transaction.Transaction, priKey string) (transaction.Transaction, error) {
//	pubkey, err := ctx.CryptoSuit.GetPubKey([]byte(priKey))
//	if err != nil {
//		return nil, err
//	}
//	tx.SetPubKey(pubkey)
//	sig, err := ctx.CryptoSuit.Sign(tx.GetSignBytes(), []byte(priKey))
//	if err != nil {
//		return nil, err
//	}
//	tx.SetSignature(sig)
//	return tx, nil
//}
