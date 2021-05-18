package collactor

import (
	"encoding/hex"
	cosmosSdkCtx "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	tx "github.com/cosmos/cosmos-sdk/client/tx"
	cosmosCrypto "github.com/cosmos/cosmos-sdk/crypto"
	keys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	cosmosSdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"context"
	"fmt"
	retry "github.com/avast/retry-go"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	sdkCtx "github.com/ci123chain/ci123chain/pkg/client/context"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	ibcexported "github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"golang.org/x/sync/errgroup"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	rtyAttNum = uint(5)
	rtyAtt    = retry.Attempts(rtyAttNum)
	rtyDel    = retry.Delay(time.Millisecond * 400)
	rtyErr    = retry.LastErrorOnly(true)

	AllowUpdateAfterExpiry       = true
	AllowUpdateAfterMisbehaviour = true
)

// Chain represents the necessary data for connecting to and indentifying a chain and its counterparites
type Chain struct {
	//Key            string  `yaml:"key" json:"key"`
	ChainID        string  `yaml:"chain-id" json:"chain-id"`
	RPCAddr        string  `yaml:"rpc-addr" json:"rpc-addr"`
	AccountPrefix  string  `yaml:"account-prefix" json:"account-prefix"`
	GasAdjustment  float64 `yaml:"gas-adjustment" json:"gas-adjustment"`
	GasPrices      string  `yaml:"gas-prices" json:"gas-prices"`
	TrustingPeriod string  `yaml:"trusting-period" json:"trusting-period"`
	PrivateKey     string  `yaml:"private-key" json:"private-key"`
	ChainType      string  `yaml:"chain-type" json:"chain-type"`

	// TODO: make these private
	HomePath string           `yaml:"-" json:"-"`
	PathEnd  *PathEnd         `yaml:"-" json:"-"`
	Keybase  keys.Keyring     `yaml:"-" json:"-"`
	Client   rpcclient.Client `yaml:"-" json:"-"`
	cdc  *codec.Codec `yaml:"-" json:"-"`
	Encoding types2.EncodingConfig `yaml:"-" json:"-"`
	CosmosEncoding params.EncodingConfig `yaml:"-" json:"-"`
	//KeyOutput *helper.KeyOutput
	address sdk.AccAddress
	logger  log.Logger
	timeout time.Duration
	debug   bool

	// stores facuet addresses that have been used reciently
	faucetAddrs map[string]time.Time
}

// Chains is a collection of Chain
type Chains []*Chain

// Get returns the configuration for a given chain
func (c Chains) Get(chainID string) (*Chain, error) {
	for _, chain := range c {
		if chainID == chain.ChainID {
			addr, _ := chain.GetAddress()
			chain.address = addr
			return chain, nil
		}
	}
	return &Chain{}, fmt.Errorf("chain with ID %s is not configured", chainID)
}
// Gets returns a map chainIDs to their chains
func (c Chains) Gets(chainIDs ...string) (map[string]*Chain, error) {
	out := make(map[string]*Chain)
	for _, cid := range chainIDs {
		chain, err := c.Get(cid)
		if err != nil {
			return out, err
		}
		out[cid] = chain
	}
	return out, nil
}

// Init initializes the pieces of a chain that aren't set when it parses a configs
// NOTE: All validation of the chain should happen here.
func (c *Chain) Init(homePath string, timeout time.Duration, logger log.Logger, debug bool) error {
	keybase, err := keys.New(c.ChainID, "test", keysDir(homePath, c.ChainID), nil)
	if err != nil {
		return err
	}


	client, err := newRPCClient(c.RPCAddr, timeout)
	if err != nil {
		return err
	}

	_, err = time.ParseDuration(c.TrustingPeriod)
	if err != nil {
		return fmt.Errorf("failed to parse trusting period (%s) for chain %s", c.TrustingPeriod, c.ChainID)
	}

	//_, err = sdk.ParseDecCoins(c.GasPrices)
	if err != nil {
		return fmt.Errorf("failed to parse gas prices (%s) for chain %s", c.GasPrices, c.ChainID)
	}

	c.cdc = types2.GetCodec()
	c.Encoding = *types2.GetEncodingConfig()
	c.CosmosEncoding = c.MakeCosmosEncodingConfig()
	c.Keybase = keybase
	err = c.ImportKey()
	if err != nil {
		return fmt.Errorf("failed to parse private key for chain %s by err: %s", c.ChainID, err.Error())
	}

	//c.KeyOutput = ko
	c.Client = client
	c.HomePath = homePath
	//c.Encoding = encodingConfig
	c.logger = logger
	c.timeout = timeout
	c.debug = debug
	c.faucetAddrs = make(map[string]time.Time)

	if c.logger == nil {
		c.logger = defaultChainLogger()
	}

	return nil
}

func (c *Chain) ImportKey() error{
	if c.ChainType == ChainTypeWeelink {
		return nil
	}
	if c.KeyExists(c.PrivateKey) {
		return nil
	}
	privKey, err := hex.DecodeString(c.PrivateKey)
	if err != nil {
		return err
	}
	private:=secp256k1.PrivKey{Key: privKey}
	//todo: passphrase需要改
	s := cosmosCrypto.EncryptArmorPrivKey(&private, "","")
	err = c.Keybase.ImportPrivKey(c.PrivateKey, s, "")
	return err
}

// KeyExists returns true if there is a specified key in chain's keybase
func (c *Chain) KeyExists(name string) bool {
	k, err := c.Keybase.Key(name)
	if err != nil {
		return false
	}

	return k.GetName() == name
}

// GetAddress returns the sdk.AccAddress associated with the configred key
func (c *Chain) GetAddress() (sdk.AccAddress, error) {
	switch c.ChainType {
	case ChainTypeWeelink:
		if !c.address.Empty()  {
			return c.address, nil
		}
		privateKey, err := crypto.HexToECDSA(c.PrivateKey)
		if err != nil {
			return sdk.AccAddress{}, errors.Errorf("error format privateKey: %s", c.PrivateKey)
		}
		address := crypto.PubkeyToAddress(privateKey.PublicKey)
		return sdk.ToAccAddress(address[:]), nil
	case ChainTypeCosmos:
		srcAddr, err := c.Keybase.Key(c.PrivateKey)
		if err != nil {
			return sdk.AccAddress{}, err
		}
		return sdk.ToAccAddress(srcAddr.GetAddress()[:]), nil
	default:
		return sdk.AccAddress{}, errors.New("unknow chain type")
	}
}

func (c *Chain) GetAddressString() (string, error) {
	address, err :=  c.GetAddress()
	if err != nil {
		return "",  err
	}
	switch c.ChainType {
	case ChainTypeWeelink:
		return address.Hex(), nil
	case ChainTypeCosmos:
		return bech32.ConvertAndEncode("cosmos", address.Address.Bytes())
	default:
		return "", errors.New("unknow chain type")
	}
}

func (c *Chain) MustGetAddressString() string {
	address, err := c.GetAddressString()
	if err != nil {
		panic(err)
	}
	return address
}

// CLIContext returns an instance of client.Context derived from Chain
func (c *Chain) CLIContext(height int64) sdkCtx.Context {
	addr, _ := c.GetAddress()
	return sdkCtx.Context{}.
		WithChainID(c.ChainID).
		WithFrom(addr).
		WithHeight(height).
		WithCodec(c.cdc).
		//WithJSONMarshaler(newContextualStdCodec(c.Encoding.Marshaler, c.UseSDKContext)).
		WithInterfaceRegistry(c.Encoding.InterfaceRegistry).
		//WithTxConfig(c.Encoding.TxConfig).
		//WithLegacyAmino(c.Encoding.Amino).
		//WithInput(os.Stdin).
		//WithNodeURI(c.RPCAddr).
		WithClient(c.Client)
		//WithAccountRetriever(authTypes.AccountRetriever{}).
		//WithBroadcastMode(flags.BroadcastBlock).
		//WithKeyring(c.Keybase).
		//WithOutputFormat("json").
		//WithFrom(c.Key).
		//WithFromName(c.Key).
		//WithFromAddress(c.MustGetAddress()).
		//WithSkipConfirmation(true).
		//WithNodeURI(c.RPCAddr).
		//WithHeight(height)
}

// CLIContext returns an instance of client.Context derived from Chain
func (c *Chain) CLICosmosContext(height int64) cosmosSdkCtx.Context {
	return cosmosSdkCtx.Context{}.
		WithChainID(c.ChainID).
		WithJSONMarshaler(newContextualStdCodec(c.CosmosEncoding.Marshaler, c.UseSDKContext)).
		WithInterfaceRegistry(c.CosmosEncoding.InterfaceRegistry).
		WithTxConfig(c.CosmosEncoding.TxConfig).
		WithLegacyAmino(c.CosmosEncoding.Amino).
		WithInput(os.Stdin).
		WithNodeURI(c.RPCAddr).
		WithClient(c.Client).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithKeyring(c.Keybase).
		WithOutputFormat("json").
		WithFrom(c.PrivateKey).
		WithFromName(c.PrivateKey).
		WithFromAddress(c.MustGetAddress().Bytes()).
		WithSkipConfirmation(true).
		WithNodeURI(c.RPCAddr).
		WithHeight(height)
}

// TxFactory returns an instance of tx.Factory derived from
func (c *Chain) TxFactory(height int64) tx.Factory {
	ctx := c.CLICosmosContext(height)
	return tx.Factory{}.
		WithAccountRetriever(ctx.AccountRetriever).
		WithChainID(c.ChainID).
		WithTxConfig(ctx.TxConfig).
		WithGasAdjustment(c.GasAdjustment).
		WithGasPrices(c.GasPrices).
		WithKeybase(c.Keybase).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
}

var sdkContextMutex sync.Mutex

func (c *Chain) UseSDKContext() func() {
	// Ensure we're the only one using the global context,
	// lock context to begin function
	sdkContextMutex.Lock()

	// Mutate the sdkConf
	sdkConf := cosmosSdk.GetConfig()
	sdkConf.SetBech32PrefixForAccount(c.AccountPrefix, c.AccountPrefix+"pub")
	sdkConf.SetBech32PrefixForValidator(c.AccountPrefix+"valoper", c.AccountPrefix+"valoperpub")
	sdkConf.SetBech32PrefixForConsensusNode(c.AccountPrefix+"valcons", c.AccountPrefix+"valconspub")

	// Return the unlock function, caller must lock and ensure that lock is released
	// before any other function needs to use c.UseSDKContext
	return sdkContextMutex.Unlock
}

// KeysDir returns the path to the keys for this chain
func keysDir(home, chainID string) string {
	return path.Join(home, "keys", chainID)
}

func defaultChainLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

// KeysDir returns the path to the keys for this chain
func KeysDir(home, chainID string) string {
	return path.Join(home, "keys", chainID)
}


func newRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}

	httpClient.Timeout = timeout
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	return rpcClient, nil
}
func lightDir(home string) string {
	return path.Join(home, "light")
}


// GetTrustingPeriod returns the trusting period for the chain
func (c *Chain) GetTrustingPeriod() time.Duration {
	tp, _ := time.ParseDuration(c.TrustingPeriod)
	return tp
}

// Log takes a string and logs the data
func (c *Chain) Log(s string) {
	c.logger.Info(s)
}

// Error takes an error, wraps it in the chainID and logs the error
func (c *Chain) Error(err error) {
	c.logger.Error(fmt.Sprintf("%s: err(%s)", c.ChainID, err.Error()))
}


// MustGetAddress used for brevity
func (c *Chain) MustGetAddress() sdk.AccAddress {
	srcAddr, err := c.GetAddress()
	if err != nil {
		panic(err)
	}

	return srcAddr
}

// Start the client service
func (c *Chain) Start() error {
	return c.Client.Start()
}


// Subscribe returns channel of events given a query
func (c *Chain) Subscribe(query string) (<-chan ctypes.ResultEvent, context.CancelFunc, error) {
	suffix, err := GenerateRandomString(8)
	if err != nil {
		return nil, nil, err
	}

	subscriber := fmt.Sprintf("%s-subscriber-%s", c.ChainID, suffix)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	eventChan, err := c.Client.Subscribe(ctx, subscriber, query, 1000)
	return eventChan, cancel, err
}


// SendMsg wraps the msg in a stdtx, signs and sends it
func (c *Chain) SendMsg(datagram sdk.Msg) (*sdk.TxResponse, bool, error) {
	return c.SendMsgs([]sdk.Msg{datagram})
}

// SendMsgs wraps the msgs in a StdTx, signs and sends it. An error is returned if there
// was an issue sending the transaction. A successfully sent, but failed transaction will
// not return an error. If a transaction is successfully sent, the result of the execution
// of that transaction will be logged. A boolean indicating if a transaction was successfully
// sent and executed successfully is returned.
func (c *Chain) SendMsgs(msgs []sdk.Msg) (*sdk.TxResponse, bool, error) {

	switch c.ChainType {
	case ChainTypeWeelink:
		ctx := c.CLIContext(0)

		pbMsgs := []sdk.PbMsg{}
		for _, v := range msgs {
			pdMsg := v.(sdk.PbMsg)
			pbMsgs = append(pbMsgs, pdMsg)
		}

		nonce, err := c.QueryNonce()
		if err != nil {
			return nil, false, err
		}
		gas, err := strconv.ParseInt(c.GasPrices, 10, 64)
		if err != nil {
			return nil, false, err
		}
		txByte, err := types2.SignPbTx(c.MustGetAddress(), nonce, uint64(gas), pbMsgs, c.PrivateKey, c.Encoding.Marshaler)
		if err != nil{
			panic(err)
		}
		res, err := ctx.BroadcastSignedData(txByte)
		if res.Code != 0 {
			c.LogFailedTx(&res, err, msgs)
			return &res, false, nil
		}
		c.LogSuccessTx(&res, msgs)
		return &res, true, nil
	case ChainTypeCosmos:
		ctx := c.CLICosmosContext(0)

		cosmosMsgs := []cosmosSdk.Msg{}
		for _, v := range msgs {
			cosmosMsg := v.(cosmosSdk.Msg)
			cosmosMsgs = append(cosmosMsgs, cosmosMsg)
		}

		// Query account details
		txf, err := tx.PrepareFactory(ctx, c.TxFactory(0))
		if err != nil {
			return nil, false, err
		}

		// If users pass gas adjustment, then calculate gas
		_, adjusted, err := tx.CalculateGas(ctx.QueryWithData, txf, cosmosMsgs...)
		if err != nil {
			return nil, false, err
		}

		// Set the gas amount on the transaction factory
		txf = txf.WithGas(adjusted)

		// Build the transaction builder
		txb, err := tx.BuildUnsignedTx(txf, cosmosMsgs...)
		if err != nil {
			return nil, false, err
		}

		// Attach the signature to the transaction
		err = tx.Sign(txf, c.PrivateKey, txb, false)
		if err != nil {
			return nil, false, err
		}

		// Generate the transaction bytes
		txBytes, err := ctx.TxConfig.TxEncoder()(txb.GetTx())
		if err != nil {
			return nil, false, err
		}

		// Broadcast those bytes
		cosmosRes, err := ctx.BroadcastTx(txBytes)
		if err != nil {
			return nil, false, err
		}

		var tx sdk.Tx
		err = c.CosmosEncoding.InterfaceRegistry.UnpackAny(cosmosRes.Tx, &tx)
		if err != nil{
			return nil, false, err
		}

		var logs sdk.ABCIMessageLogs
		for _, v := range cosmosRes.Logs {
			events := sdk.StringEvents{}
			for _, e := range v.Events {
				attributes := []sdk.Attribute{}
				for _, a := range e.Attributes {
					attributes = append(attributes, sdk.Attribute{
						Key:   []byte(a.Key),
						Value: []byte(a.Value),
						Index: false,
					})
				}
				events =append(events, sdk.StringEvent{
					Type:       e.Type,
					Attributes: attributes,
				})
			}
			logs = append(logs, sdk.ABCIMessageLog{
				MsgIndex: v.MsgIndex,
				Success:  false,
				Log:      v.Log,
				Events:   events,
			})
		}

		res := &sdk.TxResponse{
			Height:     cosmosRes.Height,
			TxHash:     cosmosRes.TxHash,
			//Index:      cosmosRes.in,
			Code:       cosmosRes.Code,
			//FormatData: cosmosRes.a,
			Data:       cosmosRes.Data,
			RawLog:     cosmosRes.RawLog,
			Logs:       logs,
			Info:       cosmosRes.Info,
			GasWanted:  cosmosRes.GasWanted,
			GasUsed:    cosmosRes.GasUsed,
			Events:     nil,
			Codespace:  cosmosRes.Codespace,
			Tx:         tx,
			Timestamp:  cosmosRes.Timestamp,
		}
		// transaction was executed, log the success or failure using the tx response code
		// NOTE: error is nil, logic should use the returned error to determine if the
		// transaction was successfully executed.
		if res.Code != 0 {
			c.LogFailedTx(res, err, msgs)
			return res, false, nil
		}

		c.LogSuccessTx(res, msgs)
		return res, true, nil
	default:
		return nil, false, errors.New("unknow chain type")
	}
}


// ValidateClientPaths takes two chains and validates their clients
func ValidateClientPaths(src, dst *Chain) error {
	if err := src.PathEnd.Vclient(); err != nil {
		return err
	}
	if err := dst.PathEnd.Vclient(); err != nil {
		return err
	}
	return nil
}

// ValidateConnectionPaths takes two chains and validates the connections
// and underlying client identifiers
func ValidateConnectionPaths(src, dst *Chain) error {
	if err := src.PathEnd.Vclient(); err != nil {
		return err
	}
	if err := dst.PathEnd.Vclient(); err != nil {
		return err
	}
	if err := src.PathEnd.Vconn(); err != nil {
		return err
	}
	if err := dst.PathEnd.Vconn(); err != nil {
		return err
	}
	return nil
}


// ValidateChannelParams takes two chains and validates their respective channel params
func ValidateChannelParams(src, dst *Chain) error {
	if err := src.PathEnd.ValidateBasic(); err != nil {
		return err
	}
	if err := dst.PathEnd.ValidateBasic(); err != nil {
		return err
	}
	//nolint:staticcheck
	if strings.ToUpper(src.PathEnd.Order) != strings.ToUpper(dst.PathEnd.Order) {
		return fmt.Errorf("src and dst path ends must have same ORDER. got src: %s, dst: %s",
			src.PathEnd.Order, dst.PathEnd.Order)
	}
	return nil
}


// GenerateConnHandshakeProof generates all the proofs needed to prove the existence of the
// connection state on this chain. A counterparty should use these generated proofs.
func (c *Chain) GenerateConnHandshakeProof(height uint64) (clientState ibcexported.ClientState,
	clientStateProof []byte, consensusProof []byte, connectionProof []byte,
	connectionProofHeight clienttypes.Height, err error) {
	var (
		clientStateRes     *clienttypes.QueryClientStateResponse
		consensusStateRes  *clienttypes.QueryConsensusStateResponse
		connectionStateRes *connectiontypes.QueryConnectionResponse

		eg = new(errgroup.Group)
	)

	// query for the client state for the proof and get the height to query the consensus state at.
	clientStateRes, err = c.QueryClientState(int64(height))
	if err != nil {
		return nil, nil, nil, nil, clienttypes.Height{}, err
	}

	clientState, err = clienttypes.UnpackClientState(clientStateRes.ClientState)
	if err != nil {
		return nil, nil, nil, nil, clienttypes.Height{}, err
	}

	eg.Go(func() error {
		consensusStateRes, err = c.QueryClientConsensusState(int64(height), clientState.GetLatestHeight())
		return err
	})
	eg.Go(func() error {
		connectionStateRes, err = c.QueryConnection(int64(height))
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, nil, nil, nil, clienttypes.Height{}, err
	}

	return clientState, clientStateRes.Proof, consensusStateRes.Proof, connectionStateRes.Proof,
		connectionStateRes.ProofHeight, nil

}

// Update returns a new chain with updated values
func (c *Chain) Update(key, value string) (out *Chain, err error) {
	out = c
	switch key {
	case "private-key":
		out.PrivateKey = value
	//case "key":
	//	out.Key = value
	case "chain-id":
		out.ChainID = value
	case "rpc-addr":
		if _, err = rpchttp.New(value, "/websocket"); err != nil {
			return
		}
		out.RPCAddr = value
	case "gas-adjustment":
		if value == "" {
			out.GasAdjustment = 1
		} else {
			adj, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			out.GasAdjustment = adj
		}
	case "gas-prices":
		//_, err = sdk.ParseDecCoins(value)
		//_, err = sdk.ParseDecCoin(value)
		//if err != nil {
		//	return nil, err
		//}
		out.GasPrices = value
	case "account-prefix":
		out.AccountPrefix = value
	case "trusting-period":
		if _, err = time.ParseDuration(value); err != nil {
			return
		}
		out.TrustingPeriod = value
	case "chain-type":
		if value == "" {
			value = ChainTypeWeelink
		}
		out.ChainType = value
	default:
		return out, fmt.Errorf("key %s not found", key)
	}
	return out, err
}

const (
	ChainTypeWeelink = "weelink"
	ChainTypeCosmos  = "cosmos"
)