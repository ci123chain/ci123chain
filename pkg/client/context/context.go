package context

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	acc_types "github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	rpclient "github.com/tendermint/tendermint/rpc/client"
)

type Context struct {
	HomeDir 	string
	NodeURI 	string
	InputAddressed []sdk.AccAddress
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

func (ctx *Context) GetInputAddresses() ([]sdk.AccAddress, error) {
	return ctx.InputAddressed, nil
}

func (ctx *Context) GetBalanceByAddress(addr sdk.AccAddress) (uint64, error) {
	addrByte := acc_types.AddressStoreKey(addr)
	res, _, err := ctx.Query("/store/main/types", addrByte)
	if res == nil{
		return 0, errors.New("The account does not exist")
	}
	if err != nil {
		return 0, err
	}
	var acc exported.Account
	err2 := ctx.Cdc.UnmarshalBinaryLengthPrefixed(res, &acc)
	if err2 != nil {
		return 0, err2
	}
	balance := uint64(acc.GetCoin().Amount.Int64())
	return balance, nil
}

func (ctx *Context) GetNonceByAddress(addr sdk.AccAddress) (uint64, error) {
	addrByte := acc_types.AddressStoreKey(addr)
	res, _, err := ctx.Query("/store/main/types", addrByte)
	if res == nil{
		return 0, errors.New("The account does not exist")
	}
	if err != nil {
		return 0, err
	}
	var acc exported.Account
	err2 := ctx.Cdc.UnmarshalBinaryLengthPrefixed(res, &acc)
	if err2 != nil {
		return 0, err2
	}
	nonce := acc.GetSequence()
	return nonce, nil
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

//
//func (ctx *Context) SignAndBroadcastTx(tx transaction.Transaction, addr types.AccAddress) (types.TxResponse, error) {
//	sig, err := ctx.Sign(tx.GetSignBytes(), addr)
//	if err != nil {
//		return types.TxResponse{}, err
//	}
//	tx.SetSignature(sig)
//	res, err := ctx.BroadcastTx(tx.Bytes())
//
//	if err != nil {
//		return res, err
//	}
//	if ctx.Verbose {
//		fmt.Printf("txHash=%v BlockHeight=%v\n", res.TxHash, res.Height)
//	}
//	return res, nil
//}

//func (ctx *Context) SignTx(tx transfer.Transaction, addr types.AccAddress) (transfer.Transaction, error) {
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

func (ctx *Context) BroadcastSignedData(data []byte) (sdk.TxResponse, error) {
	res, err := ctx.BroadcastTx(data)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	if ctx.Verbose {
		fmt.Printf("txHash=%v BlockHeight=%v\n", res.TxHash, res.Height)
	}
	return res, nil
}

func (ctx *Context) BroadcastSignedDataAsync(data []byte) (sdk.TxResponse, error) {
	_, _ = ctx.BroadcastTxAsync(data)

	return sdk.TxResponse{}, nil
}


//func (ctx *Context) SignTx2(tx transfer.Transaction, priKey string) (transfer.Transaction, error) {
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
