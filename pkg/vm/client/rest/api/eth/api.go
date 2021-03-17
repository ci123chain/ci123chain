package eth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	accounttypes "github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	clientcontext "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	vmmodule "github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/big"
	"os"
)

const (
	DefaultRPCGasLimit = 10000000
	ChainID = 999
)

var cdc = types.MakeCodec()

type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`
}

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From     *common.Address `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Uint64 `json:"gas_price"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	Value    *hexutil.Big    `json:"value"`
	Data     *hexutil.Bytes  `json:"data"`
}

// PublicEthereumAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicEthereumAPI struct {
	ctx          context.Context
	clientCtx    clientcontext.Context
	chainIDEpoch *big.Int
	logger       log.Logger
	cdc 		 *codec.Codec
	ks             *keystore.KeyStore
}

// Transaction represents a transaction returned to RPC clients.
type Transaction struct {
	BlockHash        *common.Hash    `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Gas              hexutil.Uint64  `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            hexutil.Bytes   `json:"input"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex *hexutil.Uint64 `json:"transactionIndex"`
	Value            *hexutil.Big    `json:"value"`
	V                *hexutil.Big    `json:"v"`
	R                *hexutil.Big    `json:"r"`
	S                *hexutil.Big    `json:"s"`
}

func NewTransaction(tx *types.MsgEthereumTx, txHash, blockHash common.Hash, blockNumber, index uint64) (*Transaction, error) {
	from, _ := tx.VerifySig(big.NewInt(ChainID))
	rpcTx := &Transaction{
		From:     from,
		Gas:      hexutil.Uint64(tx.Data.GasLimit),
		GasPrice: (*hexutil.Big)(tx.Data.Price),
		Hash:     txHash,
		Input:    hexutil.Bytes(tx.Data.Payload),
		Nonce:    hexutil.Uint64(tx.Data.AccountNonce),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Data.Amount),
		V:        (*hexutil.Big)(tx.Data.V),
		R:        (*hexutil.Big)(tx.Data.R),
		S:        (*hexutil.Big)(tx.Data.S),
	}

	if blockHash != (common.Hash{}) {
		rpcTx.BlockHash = &blockHash
		rpcTx.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		rpcTx.TransactionIndex = (*hexutil.Uint64)(&index)
	}

	return rpcTx, nil
}

// NewAPI creates an instance of the public ETH Web3 API.
func NewAPI(clientCtx clientcontext.Context, ks *keystore.KeyStore) *PublicEthereumAPI {

	cdc := types.MakeCodec()
	epoch := big.NewInt(ChainID)
	api := &PublicEthereumAPI{
		ctx:          context.Background(),
		clientCtx:    clientCtx.WithBlocked(false),
		chainIDEpoch: epoch,
		logger:       log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "json-rpc", "namespace", "eth"),
		cdc: 		  cdc,
		ks: 		  ks,
	}
	return api
}

// ClientCtx returns the Cosmos SDK client context.
func (api *PublicEthereumAPI) ClientCtx() clientcontext.Context {
	return api.clientCtx
}

// ChainId returns the chain's identifier in hex format
func (api *PublicEthereumAPI) ChainId() (hexutil.Uint, error) { // nolint
	api.logger.Debug("eth_chainId")
	return hexutil.Uint(uint(api.chainIDEpoch.Uint64())), nil
}

// Syncing returns whether or not the current node is syncing with other peers. Returns false if not, or a struct
// outlining the state of the sync if it is.
func (api *PublicEthereumAPI) Syncing() (interface{}, error) {
	api.logger.Debug("eth_syncing")

	status, err := api.clientCtx.Client.Status()
	if err != nil {
		return false, err
	}

	if !status.SyncInfo.CatchingUp {
		return false, nil
	}

	return map[string]interface{}{
		// "startingBlock": nil, // NA
		"currentBlock": hexutil.Uint64(status.SyncInfo.LatestBlockHeight),
		// "highestBlock":  nil, // NA
		// "pulledStates":  nil, // NA
		// "knownStates":   nil, // NA
	}, nil
}

// Coinbase is the address that staking rewards will be send to (alias for Etherbase).
func (api *PublicEthereumAPI) Coinbase() (common.Address, error) {
	api.logger.Debug("eth_coinbase")

	node, err := api.clientCtx.GetNode()
	if err != nil {
		return common.Address{}, err
	}

	status, err := node.Status()
	if err != nil {
		return common.Address{}, err
	}

	return common.BytesToAddress(status.ValidatorInfo.Address.Bytes()), nil
}

// Mining returns whether or not this node is currently mining. Always false.
func (api *PublicEthereumAPI) Mining() bool {
	api.logger.Debug("eth_mining")
	return false
}

// Hashrate returns the current node's hashrate. Always 0.
func (api *PublicEthereumAPI) Hashrate() hexutil.Uint64 {
	api.logger.Debug("eth_hashrate")
	return 0
}

// GasPrice returns the current gas price based on Eth's gas price oracle.
func (api *PublicEthereumAPI) GasPrice() *hexutil.Big {
	api.logger.Debug("eth_gasPrice")
	out := big.NewInt(1)
	return (*hexutil.Big)(out)
}

func (api *PublicEthereumAPI) BlockNumber() (hexutil.Uint64, error) {
	api.logger.Debug("eth_blockNumber")
	info, err := api.clientCtx.Client.ABCIInfo()
	if err != nil {
		return 0, err
	}
	api.logger.Debug(fmt.Sprintf("%d", info.Response.LastBlockHeight))
	return hexutil.Uint64(info.Response.LastBlockHeight), nil
}

func (api *PublicEthereumAPI) GetBlockByNumber(blockNum BlockNumber, fullTx bool) (map[string]interface{}, error) {
	api.logger.Debug("eth_getBlockByNumber", "number", blockNum, "full", fullTx)
	height := blockNum.Int64()
	if height <= 0 {
		// get latest block height
		num, err := api.BlockNumber()
		if err != nil {
			return nil, err
		}
		height = int64(num)
	}

	resBlock, err := api.clientCtx.Client.Block(&height)
	if err != nil {
		return nil, err
	}

	var transactions []common.Hash
	for _, tx := range resBlock.Block.Txs {
		transactions = append(transactions, common.BytesToHash(tx.Hash()))
	}

	return EthBlockFromTendermint(api.clientCtx, resBlock.Block, transactions)
}

// Accounts returns the list of accounts available to this node.
func (api *PublicEthereumAPI) Accounts() ([]common.Address, error) {
	api.logger.Debug("eth_accounts")
	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	accs := api.ks.Accounts()
	for i := range accs {
		addresses = append(addresses, accs[i].Address)
	}

	return addresses, nil
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (api *PublicEthereumAPI) GetTransactionCount(address common.Address, blockNum rpc.BlockNumber) (*hexutil.Uint64, error) {
	api.logger.Debug("eth_getTransactionCount", "address", address, "block number", blockNum)
	clientCtx := api.clientCtx.WithHeight(blockNum.Int64())

	// Get nonce (sequence) from account
	from := sdk.AccAddress{address}

	qparams := keeper.NewQueryAccountParams(from)
	bz, err := cdc.MarshalJSON(qparams)
	if err != nil {
		return nil, err
	}
	res, _, _, err := clientCtx.Query("/custom/" + account.ModuleName + "/" + accounttypes.QueryAccount, bz, false)
	if res == nil {
		return nil, errors.New("The account does not exist")
	}
	if err != nil {
		return nil, err
	}
	var acc exported.Account
	err2 := cdc.UnmarshalBinaryLengthPrefixed(res, &acc)
	if err2 != nil {
		return nil, err2
	}

	nonce := acc.GetSequence()
	n := hexutil.Uint64(nonce)
	return &n, nil
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (api *PublicEthereumAPI) GetStorageAt(address common.Address, key string, blockNum BlockNumber) (hexutil.Bytes, error) {
	api.logger.Debug("eth_getStorageAt", "address", address, "key", key, "block number", blockNum)
	clientCtx := api.clientCtx.WithHeight(blockNum.Int64())
	res, _, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/storage/%s/%s", evmtypes.ModuleName, address.Hex(), key), nil, false)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResStorage
	cdc.MustUnmarshalJSON(res, &out)
	return out.Value, nil
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (api *PublicEthereumAPI) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
	api.logger.Debug("eth_getTransactionReceipt", "hash", hash)
	tx, err := api.clientCtx.Client.Tx(hash.Bytes(), false)
	if err != nil {
		// Return nil for transaction when not found
		return nil, nil
	}

	var sdkTx sdk.Tx
	err = cdc.UnmarshalBinaryBare(tx.Tx, &sdkTx)
	if err != nil {
		return nil, err
	}

	ethTx, ok := sdkTx.(*types.MsgEthereumTx)
	if !ok {
		return nil, errors.New("not msg ethereumTx")
	}

	from, err := ethTx.VerifySig(ethTx.ChainID())
	if err != nil {
		return nil, err
	}

	// Query block for consensus hash
	block, err := api.clientCtx.Client.Block(&tx.Height)
	if err != nil {
		return nil, err
	}

	blockHash := common.BytesToHash(block.Block.Header.Hash())

	// Set status codes based on tx result
	var status hexutil.Uint
	if tx.TxResult.IsOK() {
		status = hexutil.Uint(1)
	} else {
		status = hexutil.Uint(0)
	}

	txData := tx.TxResult.GetData()

	data, err := evmtypes.DecodeResultData(txData)
	if err != nil {
		status = 0 // transaction failed
	}

	if len(data.Logs) == 0 {
		data.Logs = []*ethtypes.Log{}
	}

	receipt := map[string]interface{}{
		// Consensus fields: These fields are defined by the Yellow Paper
		"status":            status,
		"cumulativeGasUsed": nil, // ignore until needed
		"logsBloom":         data.Bloom,
		"logs":              data.Logs,

		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		"transactionHash": hash,
		"contractAddress": data.ContractAddress,
		"gasUsed":         hexutil.Uint64(tx.TxResult.GasUsed),

		// Inclusion information: These fields provide information about the inclusion of the
		// transaction corresponding to this receipt.
		"blockHash":        blockHash,
		"blockNumber":      hexutil.Uint64(tx.Height),
		"transactionIndex": hexutil.Uint64(tx.Index),

		// sender and receiver (contract or EOA) addreses
		"from": from,
		"to":   ethTx.To(),
	}

	return receipt, nil
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (api *PublicEthereumAPI) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	api.logger.Debug("eth_sign", "address", address, "data", data)
	// TODO: Change this functionality to find an unlocked account by address

	acc := accounts.Account{
		Address: address,
	}
	// Sign the requested hash with the wallet
	signature, err := api.ks.SignHash(acc, data)
	if err != nil {
		return nil, err
	}

	signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// SendTransaction sends an Ethereum transaction.
func (api *PublicEthereumAPI) SendTransaction(args SendTxArgs) (common.Hash, error) {
	api.logger.Debug("eth_sendTransaction", "args")
	// TODO: Change this functionality to find an unlocked account by address

	if args.Nonce == nil {
		nonce, _, err := api.clientCtx.GetNonceByAddress(sdk.AccAddress{args.From}, false)
		if err != nil {
			return common.Hash{}, err
		}
		x := hexutil.Uint64(nonce)
		args.Nonce = &x
	}
	// Assemble transaction from fields
	tx, err := api.generateFromArgs(args)
	if err != nil {
		api.logger.Debug("failed to generate tx", "error", err)
		return common.Hash{}, err
	}

	chainID := big.NewInt(ChainID)
	hash := tx.RLPSignBytes(chainID)
	sig, err := api.ks.SignHash(accounts.Account{Address: args.From}, hash.Bytes())
	if err != nil {
		return common.Hash{}, err
	}

	if len(sig) != 65 {
		return common.Hash{}, fmt.Errorf("wrong size for signature: got %d, want 65", len(sig))
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])

	var v *big.Int

	if chainID.Sign() == 0 {
		v = new(big.Int).SetBytes([]byte{sig[64] + 27})

	} else {
		v = big.NewInt(int64(sig[64] + 35))
		chainIDMul := new(big.Int).Mul(chainID, big.NewInt(2))

		v.Add(v, chainIDMul)
	}

	if len(sig) != 65 {
		return common.Hash{}, fmt.Errorf("wrong size for signature: got %d, want 65", len(sig))
	}
	tx.Data.V = v
	tx.Data.R = r
	tx.Data.S = s

	// Broadcast transaction in async mode (default)
	// NOTE: If error is encountered on the node, the broadcast will not return an error
	res, err := api.clientCtx.BroadcastSignedTx(tx.Bytes())
	if err != nil {
		return common.Hash{}, err
	}

	// Return transaction hash
	return common.HexToHash(res.TxHash), nil
}

// SendRawTransaction send a raw Ethereum transaction.
func (api *PublicEthereumAPI) SendRawTransaction(data hexutil.Bytes) (common.Hash, error) {
	api.logger.Debug("eth_sendRawTransaction", "data", data)
	tx := new(types.MsgEthereumTx)

	// RLP decode raw transaction bytes
	if err := rlp.DecodeBytes(data, tx); err != nil {
		// Return nil is for when gasLimit overflows uint64
		return common.Hash{}, nil
	}

	//// TODO: Possibly log the contract creation address (if recipient address is nil) or tx data
	//// If error is encountered on the node, the broadcast will not return an error
	res, err := api.clientCtx.BroadcastSignedTx(tx.Bytes())
	if err != nil {
		api.logger.Debug("eth_sendRawTransaction", "err", err)
		return common.Hash{}, err
	}
	api.logger.Debug("sendRawTransaction response log", "log", res.Log)

	// Return transaction hash
	return common.HexToHash(res.TxHash), nil
}

// GetCode returns the contract code at the given address and block number.
func (api *PublicEthereumAPI) GetCode(address common.Address, blockNumber BlockNumber) (hexutil.Bytes, error) {
	api.logger.Debug("eth_getCode", "address", address, "block number", blockNumber)
	clientCtx := api.clientCtx.WithHeight(blockNumber.Int64())
	res, _, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s/%s", vmmodule.ModuleName, evmtypes.QueryCode, address.Hex()), nil, false)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResCode
	cdc.MustUnmarshalJSON(res, &out)
	return out.Code, nil
}

func (api *PublicEthereumAPI) GetTransactionByHash(hash common.Hash) (*Transaction, error) {
	api.logger.Debug("eth_getTransactionByHash", "hash", hash)
	tx, err := api.clientCtx.Client.Tx(hash.Bytes(), false)
	if err != nil {
		// Return nil for transaction when not found
		return nil, nil
	}

	// Can either cache or just leave this out if not necessary
	block, err := api.clientCtx.Client.Block(&tx.Height)
	if err != nil {
		return nil, err
	}

	blockHash := common.BytesToHash(block.Block.Header.Hash())

	rawtx, err2 := types.DefaultTxDecoder(cdc)(tx.Tx)
	if err2 != nil {
		return nil, err2
	}

	height := uint64(tx.Height)
	s, err :=  NewTransaction(rawtx.(*types.MsgEthereumTx), common.BytesToHash(tx.Tx.Hash()), blockHash, height, uint64(tx.Index))
	return s, err
}

// GetBalance returns the provided account's balance up to the provided block number.
func (api *PublicEthereumAPI) GetBalance(address common.Address, blockNum BlockNumber) (*hexutil.Big, error) {
	api.logger.Debug("eth_getBalance", "address", address, "block number", blockNum)
	clientCtx := api.clientCtx.WithHeight(blockNum.Int64())

	qparams := keeper.NewQueryAccountParams(sdk.AccAddress{address})
	bz, err := cdc.MarshalJSON(qparams)
	if err != nil {
		return nil, err
	}
	res, _, _, err := clientCtx.Query("/custom/" + account.ModuleName + "/" + accounttypes.QueryAccount, bz, false)
	if res == nil {
		return nil, errors.New("The account does not exist")
	}
	if err != nil {
		return nil, err
	}
	var acc exported.Account
	err2 := cdc.UnmarshalBinaryLengthPrefixed(res, &acc)
	if err2 != nil {
		return nil, err2
	}

	val := acc.GetCoin().Amount.BigInt()
	api.logger.Debug("eth_getBalance", "balance", val)
	return (*hexutil.Big)(val), nil
}

// Call performs a raw contract call.
func (api *PublicEthereumAPI) Call(args CallArgs, _ BlockNumber, _ *map[common.Address]Account) (hexutil.Bytes, error) {
	api.logger.Debug("eth_call", "args")
	simRes, err := api.doCall(args, big.NewInt(DefaultRPCGasLimit))
	if err != nil {
		return []byte{}, err
	}

	data, err := evmtypes.DecodeResultData([]byte(simRes.FormatData))
	if err != nil {
		return []byte{}, err
	}

	return (hexutil.Bytes)(data.Ret), nil
}

// DoCall performs a simulated call operation through the evmtypes. It returns the
// estimated gas used on the operation or an error if fails.
func (api *PublicEthereumAPI) doCall(
	args CallArgs, globalGasCap *big.Int,
) (*sdk.QureyAppResponse, error) {
	// Set height for historical queries
	clientCtx := api.clientCtx

	// Set sender address or use a default if none specified
	var addr common.Address

	if args.From == nil {
		addrs, err := api.Accounts()
		if err == nil && len(addrs) > 0 {
			addr = addrs[0]
		}
	} else {
		addr = *args.From
	}

	if args.Nonce == nil {
		nonce, _, err := api.clientCtx.GetNonceByAddress(sdk.AccAddress{addr}, false)
		if err != nil {
			return nil, err
		}
		x := hexutil.Uint64(nonce)
		args.Nonce = &x
	}
	// Set default gas & gas price if none were set
	// Change this to uint64(math.MaxUint64 / 2) if gas cap can be configured
	gas := uint64(DefaultRPCGasLimit)
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	if globalGasCap != nil && globalGasCap.Uint64() < gas {
		api.logger.Debug("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap.Uint64()
	}

	// Set gas price using default or parameter if passed in
	gasPrice := big.NewInt(1)

	// Set value for transaction
	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}

	// Set Data if provided
	var data []byte
	if args.Data != nil {
		data = []byte(*args.Data)
	}

	// Create new call message
	msg := evmtypes.NewMsgEvmTx(sdk.AccAddress{addr}, uint64(*args.Nonce), args.To, value, gas, gasPrice, data)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	txBytes, err := api.cdc.MarshalBinaryBare(types.NewCommonTx(sdk.AccAddress{addr}, uint64(*args.Nonce), gas, []sdk.Msg{msg}))
	if err != nil {
		return nil, err
	}
	// Transaction simulation through query
	res, _, _, err := clientCtx.Query("app/simulate", txBytes, false)
	if err != nil {
		return nil, err
	}
	var simResponse sdk.QureyAppResponse
	if err := api.cdc.UnmarshalBinaryBare(res, &simResponse); err != nil {
		return nil, err
	}

	return &simResponse, nil
}
//// generateFromArgs populates tx message with args (used in RPC API)
//func (api *PublicEthereumAPI) generateFromArgs(args SendTxArgs) (*evmtypes.MsgEvmTx, error) {
//	amount := (*big.Int)(args.Value)
//	gasPrice := big.NewInt(1)
//	nonce := uint64(*args.Nonce)
//
//	var gasLimit uint64
//	if args.Gas == nil {
//		gasLimit = uint64(DefaultRPCGasLimit)
//	} else {
//		gasLimit = (uint64)(*args.Gas)
//	}
//
//	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
//		return nil, errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`)
//	}
//
//	// Sets input to either Input or Data, if both are set and not equal error above returns
//	var input []byte
//	if args.Input != nil {
//		input = *args.Input
//	} else if args.Data != nil {
//		input = *args.Data
//	}
//
//	if args.To == nil && len(input) == 0 {
//		// Contract creation
//		return nil, fmt.Errorf("contract creation without any data provided")
//	}
//
//	msg := evmtypes.NewMsgEvmTx(sdk.AccAddress{args.From}, nonce, args.To, amount, gasLimit, gasPrice, input)
//	return &msg, nil
//}


func (api *PublicEthereumAPI) generateFromArgs(args SendTxArgs) (*types.MsgEthereumTx, error) {
	amount := (*big.Int)(args.Value)
	gasPrice := big.NewInt(1)
	nonce := uint64(*args.Nonce)

	var gasLimit uint64
	if args.Gas == nil {
		gasLimit = uint64(DefaultRPCGasLimit)
	} else {
		gasLimit = (uint64)(*args.Gas)

	}
	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return nil, errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`)
	}

	// Sets input to either Input or Data, if both are set and not equal error above returns
	var input []byte
	if args.Input != nil {
		input = *args.Input
	} else if args.Data != nil {
		input = *args.Data
	}

	if args.To == nil && len(input) == 0 {
		// Contract creation
		return nil, fmt.Errorf("contract creation without any data provided")
	}

	msg := types.NewMsgEthereumTx(nonce, args.To, amount, gasLimit, gasPrice, input)
	return &msg, nil
}

// EthBlockFromTendermint returns a JSON-RPC compatible Ethereum blockfrom a given Tendermint block.
func EthBlockFromTendermint(clientCtx clientcontext.Context, block *tmtypes.Block, transactions []common.Hash) (map[string]interface{}, error) {
	gasLimit := DefaultRPCGasLimit
	gasUsed := big.NewInt(0)

	res, _, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s/%d", vmmodule.ModuleName, evmtypes.QueryBloom, block.Height), nil, false)
	if err != nil {
		return nil, err
	}

	var bloomRes evmtypes.QueryBloomFilter
	types.MakeCodec().MustUnmarshalJSON(res, &bloomRes)

	bloom := bloomRes.Bloom

	return formatBlock(block.Header, block.Size(), int64(gasLimit), gasUsed, transactions, bloom), nil
}

func formatBlock(
	header tmtypes.Header, size int, gasLimit int64,
	gasUsed *big.Int, transactions interface{}, bloom ethtypes.Bloom,
) map[string]interface{} {
	if len(header.DataHash) == 0 {
		header.DataHash = common.Hash{}.Bytes()
	}

	return map[string]interface{}{
		"number":           hexutil.Uint64(header.Height),
		"hash":             hexutil.Bytes(header.Hash()),
		"parentHash":       hexutil.Bytes(header.LastBlockID.Hash),
		"nonce":            hexutil.Uint64(0), // PoW specific
		"sha3Uncles":       common.Hash{},     // No uncles in Tendermint
		"logsBloom":        bloom,
		"transactionsRoot": hexutil.Bytes(header.DataHash),
		"stateRoot":        hexutil.Bytes(header.AppHash),
		"miner":            common.Address{},
		"mixHash":          common.Hash{},
		"difficulty":       0,
		"totalDifficulty":  0,
		"extraData":        hexutil.Uint64(0),
		"size":             hexutil.Uint64(size),
		"gasLimit":         hexutil.Uint64(gasLimit), // Static gas limit
		"gasUsed":          (*hexutil.Big)(gasUsed),
		"timestamp":        hexutil.Uint64(header.Time.Unix()),
		"transactions":     transactions.([]common.Hash),
		"uncles":           []string{},
		"receiptsRoot":     common.Hash{},
	}
}