package collactor

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkCtx "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/tendermint/tendermint/libs/log"
	retry "github.com/avast/retry-go"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"os"
	"path"
	"time"
)

var (
	rtyAttNum = uint(5)
	rtyAtt    = retry.Attempts(rtyAttNum)
	rtyDel    = retry.Delay(time.Millisecond * 400)
	rtyErr    = retry.LastErrorOnly(true)
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