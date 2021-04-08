package collactor

import (
	"fmt"
	retry "github.com/avast/retry-go"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkCtx "github.com/ci123chain/ci123chain/pkg/client/context"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	ibcexported "github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"golang.org/x/sync/errgroup"
	"os"
	"path"
	"strings"
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
	Key            string  `yaml:"key" json:"key"`
	ChainID        string  `yaml:"chain-id" json:"chain-id"`
	RPCAddr        string  `yaml:"rpc-addr" json:"rpc-addr"`
	AccountPrefix  string  `yaml:"account-prefix" json:"account-prefix"`
	GasAdjustment  float64 `yaml:"gas-adjustment" json:"gas-adjustment"`
	GasPrices      string  `yaml:"gas-prices" json:"gas-prices"`
	TrustingPeriod string  `yaml:"trusting-period" json:"trusting-period"`

	// TODO: make these private
	HomePath string           `yaml:"-" json:"-"`
	PathEnd  *PathEnd         `yaml:"-" json:"-"`
	Keybase  keys.Keyring     `yaml:"-" json:"-"`
	Client   rpcclient.Client `yaml:"-" json:"-"`
	//Encoding params.EncodingConfig `yaml:"-" json:"-"`

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

// Init initializes the pieces of a chain that aren't set when it parses a config
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

	//encodingConfig := c.MakeEncodingConfig()

	c.Keybase = keybase
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

// GetAddress returns the sdk.AccAddress associated with the configred key
func (c *Chain) GetAddress() (sdk.AccAddress, error) {
	if !c.address.Empty()  {
		return c.address, nil
	}

	// Signing key for c chain
	srcAddr, err := c.Keybase.Key(c.Key)
	if err != nil {
		return nil, err
	}

	return srcAddr.GetAddress(), nil
}


// CLIContext returns an instance of client.Context derived from Chain
func (c *Chain) CLIContext(height int64) sdkCtx.Context {
	addr, _ := c.GetAddress()
	return sdkCtx.Context{}.
		WithChainID(c.ChainID).
		WithFrom(addr).
		WithHeight(height)
		//WithCodec()
		//WithJSONMarshaler(newContextualStdCodec(c.Encoding.Marshaler, c.UseSDKContext)).
		//WithInterfaceRegistry(c.Encoding.InterfaceRegistry).
		//WithTxConfig(c.Encoding.TxConfig).
		//WithLegacyAmino(c.Encoding.Amino).
		//WithInput(os.Stdin).
		//WithNodeURI(c.RPCAddr).
		//WithClient(c.Client).
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



func defaultChainLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout))
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
	// Instantiate the client context
	ctx := c.CLIContext(0)

	// Query account details
	txf, err := tx.PrepareFactory(ctx, c.TxFactory(0))
	if err != nil {
		return nil, false, err
	}

	// If users pass gas adjustment, then calculate gas
	_, adjusted, err := tx.CalculateGas(ctx.QueryWithData, txf, msgs...)
	if err != nil {
		return nil, false, err
	}

	// Set the gas amount on the transaction factory
	txf = txf.WithGas(adjusted)

	// Build the transaction builder
	txb, err := tx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, false, err
	}

	// Attach the signature to the transaction
	err = tx.Sign(txf, c.Key, txb, false)
	if err != nil {
		return nil, false, err
	}

	// Generate the transaction bytes
	txBytes, err := ctx.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return nil, false, err
	}

	// Broadcast those bytes
	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		return nil, false, err
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

	clientState = clientStateRes.ClientState

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