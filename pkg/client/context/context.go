package context

import (
	"context"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
	rpclient "github.com/tendermint/tendermint/rpc/client"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const DefaultShardPort = ":8545"

type Context struct {
	HomeDir 	string
	NodeURI 	string
	ChainID 	string
	FromAddr 	sdk.AccAddress
	Blocked     bool
	Client 		rpclient.Client
	Verbose 	bool
	Height 		int64
	Cdc 		*codec.Codec
	InterfaceRegistry codectypes.InterfaceRegistry
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

// WithInterfaceRegistry returns the context with an updated InterfaceRegistry
func (ctx Context) WithInterfaceRegistry(interfaceRegistry codectypes.InterfaceRegistry) Context {
	ctx.InterfaceRegistry = interfaceRegistry
	return ctx
}

func (ctx Context) WithCodec(cdc *codec.Codec) Context {
	ctx.Cdc = cdc
	return ctx
}

func (ctx Context) WithChainID(chainID string) Context {
	ctx.ChainID = chainID
	return ctx
}

func (ctx Context) WithHeight(height int64) Context {
	ctx.Height = height
	return ctx
}

func (ctx Context) WithClient(client rpclient.Client) Context {
	ctx.Client = client
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

func (ctx *Context) GetHistoryBalance(addr sdk.AccAddress, isProve bool, height string) (sdk.Coins, *crypto.ProofOps, error, []byte) {
	var h int64
	if height == "" {
		h = -1
	}else {
		var err error
		h, err = util.CheckInt64(height)
		if err != nil {
			return sdk.NewCoins(), nil, err, nil
		}
		if h <=0 {
			return sdk.NewCoins(), nil, errors.New(fmt.Sprintf("unexpected height: %v", height)), nil
		}
	}
	qparams := keeper.NewQueryAccountParams(addr, h)
	bz, err := ctx.Cdc.MarshalJSON(qparams)
	if err != nil {
		return sdk.NewCoins(), nil , err, nil
	}
	res, _, _, err := ctx.Query("/custom/" + types.ModuleName + "/" + types.QueryHistoryAccount, bz, isProve)
	if res == nil{
		return sdk.NewCoins(), nil, errors.New("The account does not exist"), nil
	}
	if err != nil {
		return sdk.NewCoins(), nil, err, nil
	}
	var result util.HistoryAccount
	err2 := ctx.Cdc.UnmarshalBinaryLengthPrefixed(res, &result)
	if err2 != nil {
		return sdk.NewCoins(), nil, err2, nil
	}

	var acc exported.Account
	err2 = ctx.Cdc.UnmarshalBinaryLengthPrefixed(result.Account, &acc)
	if err2 != nil {
		return sdk.NewCoins(), nil, err2, nil
	}

	if result.Proof == nil && result.Shard != "" {
		var params = make(map[string]string, 0)
		params["address"] = addr.String()
		params["height"] = height
		params["prove"] = "true"
		res, err := sendRequest(result.Shard+DefaultShardPort, params)
		if err != nil {
			return nil, nil, err, nil
		}
		return nil, nil, nil, res
	}
	fmt.Println(result.Coins)
	return acc.GetCoins(), result.Proof, nil, nil
}

func (ctx *Context) GetBalanceByAddress(addr sdk.AccAddress, isProve bool, height string) (sdk.Coins, *crypto.ProofOps, error) {
	var h int64
	if height == "" {
		h = -1
	}else {
		var err error
		h, err = util.CheckInt64(height)
		if err != nil {
			return sdk.NewCoins(), nil, err
		}
		if h <=0 {
			return sdk.NewCoins(), nil, errors.New(fmt.Sprintf("unexpected height: %v", height))
		}
	}
	qparams := keeper.NewQueryAccountParams(addr, h)
	bz, err := ctx.Cdc.MarshalJSON(qparams)
	if err != nil {
		return sdk.NewCoins(), nil , err
	}
	res, _, proof, err := ctx.Query("/custom/" + types.ModuleName + "/" + types.QueryAccount, bz, isProve)
	if res == nil{
		return sdk.NewCoins(), nil, errors.New("The account does not exist")
	}
	if err != nil {
		return sdk.NewCoins(), nil, err
	}
	var acc exported.Account
	err2 := ctx.Cdc.UnmarshalBinaryLengthPrefixed(res, &acc)
	if err2 != nil {
		return sdk.NewCoins(), nil, err2
	}
	balance := acc.GetCoins()

	return balance, proof, nil
}

func (ctx *Context) GetNonceByAddress(addr sdk.AccAddress, isProve bool) (uint64, *crypto.ProofOps, error) {
	qparams := keeper.NewQueryAccountParams(addr, 0)
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
	//var nonce uint64
	//err2 := ctx.Cdc.UnmarshalBinaryLengthPrefixed(res, &nonce)
	//if err2 != nil {
	//	return 0, nil, err2
	//}

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


func sendRequest(host string, RequestParams map[string]string) ([]byte, error) {
	cli := &http.Client{
		Transport:&http.Transport{DisableKeepAlives:true},
	}
	reqUrl := "http://" + host + "/bank/history/balance"
	data := url.Values{}
	for k, v := range RequestParams {
		data.Set(k, v)
	}

	req2, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}
	req2.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))
	defer req2.Body.Close()
	//not use one connection
	req2.Close = true

	// set request content types
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// request
	rep2, err := cli.Do(req2)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(rep2.Body)
	if err != nil {
		return nil, err
	}
	defer rep2.Body.Close()
	return b, nil
}

