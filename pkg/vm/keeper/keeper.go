package keeper

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	types2 "github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/params"
	keeper2 "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	vmmodule "github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"github.com/wasmerio/go-ext-wasm/wasmer"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strings"
)

const UINT_MAX uint64 = ^uint64(0)
const INIT string = "init"
const INVOKE string = "invoke"
type Keeper struct {
	storeKey    		sdk.StoreKey
	cdc         		*codec.Codec
	wasmer     	 		Wasmer
	homeDir				string
	AccountKeeper 		account.AccountKeeper
	StakingKeeper       keeper2.StakingKeeper
	// Ethermint concrete implementation on the EVM StateDB interface
	CommitStateDB 		*evmtypes.CommitStateDB
	// Transaction counter in a block. Used on StateSB's Prepare function.
	// It is reset to 0 every block on BeginBlock so there's no point in storing the counter
	// on the KVStore or adding it as a field on the EVM genesis state.
	TxCount 			int
	Bloom   			*big.Int
	cdb					dbm.DB
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, homeDir string, wasmConfig types.WasmConfig, paramSpace params.Subspace, accountKeeper account.AccountKeeper, stakingKeeper keeper2.StakingKeeper, cdb dbm.DB) Keeper {
	wasmer, err := NewWasmer(homeDir, wasmConfig)
	if err != nil {
		panic(err)
	}

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(evmtypes.ParamKeyTable())
	}

	wk := Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		wasmer:        *wasmer,
		homeDir:       homeDir,
		AccountKeeper: accountKeeper,
		StakingKeeper: stakingKeeper,
		CommitStateDB: evmtypes.NewCommitStateDB(sdk.Context{}, storeKey, paramSpace, accountKeeper),
		TxCount:       0,
		Bloom:         big.NewInt(0),
		cdb:		   cdb,
	}

	return wk
}

//
func (k *Keeper) Upload(ctx sdk.Context, wasmCode []byte, creator sdk.AccAddress) (codeHash []byte, err error) {
	ccstore := ctx.KVStore(k.storeKey)
	var wasmer Wasmer
	wasmerBz := ccstore.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return nil, sdk.ErrInternal("empty wasmer")
		}
		k.wasmer = Wasmer{
			FilePathMap: mapFromSortMaps(wasmer.SortMaps),
			SortMaps:    nil,
			LastFileID:  wasmer.LastFileID,
		}
	}
	codeHash, isExist, err := k.create(ctx, creator, wasmCode)
	if err != nil {
		return nil, err
	}
	//store code in local
	if !isExist {
		err = ioutil.WriteFile(k.homeDir + WASMDIR + k.wasmer.FilePathMap[fmt.Sprintf("%x", codeHash)], wasmCode, types.ModePerm)
		if err != nil {
			return nil, err
		}
	} else {
		//todo somebody else has already uploaded it
	}
	return codeHash, nil
}

func (k *Keeper) Instantiate(ctx sdk.Context, codeHash []byte, invoker sdk.AccAddress, args utils.CallData, name, version, author, email, describe string, genesisContractAddress sdk.AccAddress, gasWanted uint64) (sdk.AccAddress, error) {
	// 如果是官方合约，不限制gas数量
	runtimeCfg := &runtimeConfig{
		GasUsed: 0,
		PreCaller:   invoker,
		Invoker:     invoker,
		Creator:     invoker,
		SelfAddress: sdk.AccAddress{},
		Keeper:      k,
		Context:     &ctx,
	}

	isGenesis, ok := ctx.Value(types.SystemContract).(bool)
	if ok && isGenesis {
		runtimeCfg.SetGasWanted(UINT_MAX)
	} else {
		runtimeCfg.SetGasWanted(gasWanted)
	}

	var codeInfo types.CodeInfo
	var wasmer Wasmer
	var code []byte
	ccstore := ctx.KVStore(k.storeKey)
	bz := ccstore.Get(types.GetCodeKey(codeHash))
	if bz == nil {
		return sdk.AccAddress{}, sdk.ErrInternal("codeHash not found")
	}
	var contractAddress sdk.AccAddress
	if isGenesis {
		contractAddress = genesisContractAddress
	}else {
		nonce := k.AccountKeeper.GetAccount(ctx, invoker).GetSequence()
		contractAddress = k.generateContractAddress(codeHash, invoker, args, nonce)
	}

	existingAcct := k.AccountKeeper.GetAccount(ctx, contractAddress)
	if existingAcct != nil {
		//return sdk.AccAddress{}, sdk.ErrInternal("Contract account exists")
		return contractAddress, nil
	}
	runtimeCfg.SetSelfAddr(contractAddress)

	var contractAccount exported.Account
	contractAccount = k.AccountKeeper.NewAccountWithAddress(ctx, contractAddress)
	contractAccount.SetContractType(types2.WasmContractType)
	k.AccountKeeper.SetAccount(ctx, contractAccount)

	wasmerBz := ccstore.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return sdk.AccAddress{}, sdk.ErrInternal("empty wasmer")
		}
		k.wasmer = Wasmer{
			FilePathMap: mapFromSortMaps(wasmer.SortMaps),
			SortMaps:    nil,
			LastFileID:  wasmer.LastFileID,
		}
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &codeInfo)

	wc, err := k.wasmer.GetWasmCode(k.homeDir, codeHash)
	if err != nil {
		wc = ccstore.Get(codeHash)

		fileName := k.wasmer.FilePathMap[strings.ToLower(codeInfo.CodeHash)]
		err = ioutil.WriteFile(k.homeDir + WASMDIR + fileName, wc, types.ModePerm)
		if err != nil {
			return sdk.AccAddress{}, err
		}
	}
	code = wc
	//create store

	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	runtimeCfg.SetStore(prefixStore)

	input, err := types.ArgsToInput(args)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	sig, err := utils.ParseSignature(args.Method)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	wasmRuntime := new(wasmRuntime)
	_, err = wasmRuntime.Call(code, input, sig.Method, runtimeCfg)
	if err != nil {
		return sdk.AccAddress{}, err
	}

	initArgs, _ := json.Marshal(args)
	//save the contract info.
	createdAt := types.NewCreatedAt(ctx)
	contractInfo := types.NewContractInfo(codeInfo, initArgs, name, version, author, email, describe, createdAt)
	ccstore.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryBare(contractInfo))

	//save contractAddress into account
	if !isGenesis {
		contractAddrStr := contractAddress.String()
		accountAddr := k.AccountKeeper.GetAccount(ctx, invoker).GetAddress()

		var contractList []string
		contractListBytes := ccstore.Get(types.GetAccountContractListKey(accountAddr))
		if contractListBytes != nil {
			err := json.Unmarshal(contractListBytes, &contractList)
			if err != nil{
				return sdk.AccAddress{}, err
			}
		}
		contractList = append(contractList, contractAddrStr)
		contractListBytes, err = json.Marshal(contractList)
		if err != nil{
			return sdk.AccAddress{}, err
		}
		ccstore.Set(types.GetAccountContractListKey(accountAddr), contractListBytes)
	}
	ctx.GasMeter().ConsumeGas(sdk.Gas(runtimeCfg.GasUsed),"wasm cost")
	return contractAddress, nil
}

//
func (k *Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, invoker sdk.AccAddress, args utils.CallData, gasWanted uint64) (sdk.Result, error) {
	runtimeCfg := &runtimeConfig{
		GasUsed:     0,
		GasWanted: 	 gasWanted,
		PreCaller:   invoker,
		Invoker:     invoker,
		SelfAddress: contractAddress,
		Keeper:      k,
		Context:     &ctx,
	}

	contract := k.GetContractInfo(ctx, contractAddress)
	if contract == nil {
		return sdk.Result{}, errors.New("Cannot found this contract address")
	}
	runtimeCfg.SetCreator(contract.CodeInfo.Creator)
	codeInfo, err := k.contractInstance(ctx, contractAddress)
	if err != nil {
		return sdk.Result{}, err
	}
	var code []byte
	codeHash, _ := hex.DecodeString(codeInfo.CodeHash)
	wc, err := k.wasmer.GetWasmCode(k.homeDir, codeHash)
	ccstore := ctx.KVStore(k.storeKey)
	if err != nil {
		wc = ccstore.Get(codeHash)

		fileName := k.wasmer.FilePathMap[strings.ToLower(codeInfo.CodeHash)]
		err = ioutil.WriteFile(k.homeDir + WASMDIR + fileName, wc, types.ModePerm)
		if err != nil {
			return sdk.Result{}, err
		}
	}
	code = wc
	//get store
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	runtimeCfg.SetStore(prefixStore)
	input, err := types.ArgsToInput(args)
	if err != nil {
		return sdk.Result{}, err
	}
	sig, err := utils.ParseSignature(args.Method)
	if err != nil {
		return sdk.Result{}, err
	}
	wasmRuntime := new(wasmRuntime)
	res, err := wasmRuntime.Call(code, input, sig.Method, runtimeCfg)
	if err != nil {
		return sdk.Result{}, err
	}
	ctx.GasMeter().ConsumeGas(sdk.Gas(runtimeCfg.GasUsed),"wasm cost")
	return sdk.Result{
		Data:   []byte(fmt.Sprintf("%s", string(res))),
	}, nil
}

func (k *Keeper) Migrate(ctx sdk.Context, codeHash []byte, invoker sdk.AccAddress, oldContract sdk.AccAddress, args utils.CallData, name, version, author, email, describe string, gasWanted uint64) (sdk.AccAddress, error) {
	newContract, err := k.Instantiate(ctx, codeHash, invoker, args, name, version, author, email, describe, types.EmptyAddress, gasWanted)

	if err != nil {
		return sdk.AccAddress{}, err
	}

	prefix := "s/k:" + k.storeKey.Name() + "/"
	oldKey := types.GetContractStorePrefixKey(oldContract)

	startKey := append([]byte(prefix), oldKey...)
	endKey := EndKey(startKey)

	iter := k.cdb.Iterator(startKey, endKey)
	defer iter.Close()

	prefixStoreKey := types.GetContractStorePrefixKey(newContract)
	prefixStore := NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	for iter.Valid() {
		key := string(iter.Key())
		realKey := strings.Split(key, string(startKey))
		value := iter.Value()
		prefixStore.Set([]byte(realKey[1]), value)
		iter.Next()
	}

	return newContract, nil
}

// queryContract
func (k Keeper) Query(ctx sdk.Context, contractAddress, invoker sdk.AccAddress, args utils.CallData) (types.ContractState, error) {
	runtimeCfg := &runtimeConfig{
		GasUsed:     0,
		GasWanted:   UINT_MAX,
		PreCaller:   invoker,
		Creator:     k.GetCreator(ctx, contractAddress),
		Invoker:     invoker,
		SelfAddress: contractAddress,
		Keeper:      &k,
		Context:     &ctx,
	}

	codeInfo, err := k.contractInstance(ctx, contractAddress)
	if err != nil {
		return types.ContractState{}, err
	}
	var code []byte
	store := ctx.KVStore(k.storeKey)
	codeHash, _ := hex.DecodeString(codeInfo.CodeHash)
	wc, err := k.wasmer.GetWasmCode(k.homeDir, codeHash)
	if err != nil {
		wc = store.Get(codeHash)

		fileName := k.wasmer.FilePathMap[strings.ToLower(codeInfo.CodeHash)]
		err = ioutil.WriteFile(k.homeDir + WASMDIR + fileName, wc, types.ModePerm)
		if err != nil {
			return types.ContractState{}, err
		}
	}
	code = wc

	//get store
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	runtimeCfg.SetStore(prefixStore)
	input, err := types.ArgsToInput(args)
	if err != nil {
		return types.ContractState{}, err
	}
	sig, err := utils.ParseSignature(args.Method)
	if err != nil {
		return types.ContractState{}, err
	}
	wasmRuntime := new(wasmRuntime)
	res, err := wasmRuntime.Call(code, input, sig.Method, runtimeCfg)
	if err != nil {
		return types.ContractState{}, err
	}

	contractState := types.ContractState{Result: string(res)}

	return contractState, nil
}


func (k *Keeper) contractInstance(ctx sdk.Context, contractAddress sdk.AccAddress) (types.CodeInfo, error) {
	var wasmer Wasmer
	ccstore := ctx.KVStore(k.storeKey)
	contractBz := ccstore.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return types.CodeInfo{}, sdk.ErrInternal(" get contract address failed")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)
	codeHash, _ := hex.DecodeString(contract.CodeInfo.CodeHash)
	bz := ccstore.Get(types.GetCodeKey(codeHash))
	if bz == nil {
		return types.CodeInfo{}, sdk.ErrInternal("get code key failed")
	}
	wasmerBz := ccstore.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return types.CodeInfo{}, sdk.ErrInternal("unexpected wasmer info")
		}
		k.wasmer = Wasmer{
			FilePathMap: mapFromSortMaps(wasmer.SortMaps),
			SortMaps:    nil,
			LastFileID:  wasmer.LastFileID,
		}
	}

	var codeInfo types.CodeInfo
	k.cdc.MustUnmarshalBinaryBare(bz, &codeInfo)
	return codeInfo, nil
}

func (k Keeper) GetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) *types.ContractInfo {

	store := ctx.KVStore(k.storeKey)
	var contract types.ContractInfo
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)
	return &contract
}

func (k Keeper) SetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress, contract types.ContractInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryBare(contract))
}

func (k Keeper) GetCodeInfo(ctx sdk.Context, codeHash []byte) *types.CodeInfo {
	store := ctx.KVStore(k.storeKey)
	var codeInfo types.CodeInfo
	codeInfoBz := store.Get(types.GetCodeKey(codeHash))
	if codeInfoBz == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(codeInfoBz, &codeInfo)
	return &codeInfo
}

func (k Keeper) GetCreator(ctx sdk.Context, contractAddress sdk.AccAddress) sdk.AccAddress {

	store := ctx.KVStore(k.storeKey)
	var contract types.ContractInfo
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return sdk.AccAddress{}
	}
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)

	return contract.CodeInfo.Creator
}

func (k *Keeper) create(ctx sdk.Context, invokerAddr sdk.AccAddress, wasmCode []byte) (codeHash []byte, isExist bool, err error) {
	wasmCode, err = UnCompress(wasmCode)
	if err != nil {
		return nil, false, err
	}
	//checks if the file contents are of wasm binary
	ok := IsValidaWasmFile(wasmCode)
	if ok != nil {
		return nil, false, ok
	}
	ccstore := ctx.KVStore(k.storeKey)
	codeHash = MakeCodeHash(wasmCode)

	// addgas
	wasmCode, err = tryAddgas(wasmCode)
	if err != nil {
		return nil, false, err
	}
	//check if it has been saved in couchDB.

	codeByte := ccstore.Get(codeHash)
	if codeByte != nil {
		hash := fmt.Sprintf("%x", codeHash)
		filePath := path.Join(k.homeDir, WASMDIR, k.wasmer.FilePathMap[hash])
		if FileExist(filePath) {
			//the file content needs to be one
			localCode, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, false, err
			}
			localFileHash := MakeCodeHash(localCode)
			//the content if different, delete local file and save remote file.
			if !bytes.Equal(localFileHash, codeHash) {
				err = os.Remove(filePath)
				if err != nil {
					return nil, false, err
				}
			}
			err = ioutil.WriteFile(filePath, wasmCode, types.ModePerm)
			if err != nil {
				return nil, false, err
			}
			return codeHash, true,nil
		}else {
			err = ioutil.WriteFile(filePath, wasmCode, types.ModePerm)
			if err != nil {
				return nil, false, err
			}
			return codeHash, true,nil
		}
	}
	newWasmer, err := k.wasmer.Create(k.homeDir, fmt.Sprintf("%x", codeHash))
	if err != nil {
		return nil, false, err
	}
	bz := k.cdc.MustMarshalJSON(newWasmer)
	if bz == nil {
		return nil, false, sdk.ErrInternal("marshal json failed")
	}
	ccstore.Set(types.GetWasmerKey(), bz)
	ccstore.Set(codeHash, wasmCode)
	codeInfo := types.NewCodeInfo(strings.ToUpper(hex.EncodeToString(codeHash)), invokerAddr)
	ccstore.Set(types.GetCodeKey(codeHash), k.cdc.MustMarshalBinaryBare(codeInfo))

	return codeHash, false, nil
}

func (k Keeper) generateContractAddress(codeHash []byte, creatorAddr sdk.AccAddress, payload utils.CallData, nonce uint64) sdk.AccAddress {
	contract, _ := rlp.EncodeToBytes([]interface{}{codeHash, creatorAddr, payload, nonce})
	return sdk.ToAccAddress(crypto.Keccak256Hash(contract).Bytes()[12:])
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", vmmodule.ModuleName))
}

// ----------------------------------------------------------------------------
// Block hash mapping functions
// Required by Web3 API.
//  TODO: remove once tendermint support block queries by hash.
// ----------------------------------------------------------------------------

// GetBlockHash gets block height from block consensus hash
func (k *Keeper) GetBlockHash(ctx sdk.Context, hash []byte) (int64, bool) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixBlockHash)
	bz := store.Get(hash)
	if len(bz) == 0 {
		return 0, false
	}

	height := binary.BigEndian.Uint64(bz)
	return int64(height), true
}

// SetBlockHash sets the mapping from block consensus hash to block height
func (k *Keeper) SetBlockHash(ctx sdk.Context, hash []byte, height int64) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixBlockHash)
	bz := sdk.Uint64ToBigEndian(uint64(height))
	store.Set(hash, bz)
}

// ----------------------------------------------------------------------------
// Block bloom bits mapping functions
// Required by Web3 API.
// ----------------------------------------------------------------------------

// GetBlockBloom gets bloombits from block height
func (k *Keeper) GetBlockBloom(ctx sdk.Context, height int64) (ethtypes.Bloom, bool) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixBloom)
	has := store.Has(evmtypes.BloomKey(height))
	if !has {
		return ethtypes.Bloom{}, false
	}

	bz := store.Get(evmtypes.BloomKey(height))
	return ethtypes.BytesToBloom(bz), true
}

// SetBlockBloom sets the mapping from block height to bloom bits
func (k *Keeper) SetBlockBloom(ctx sdk.Context, height int64, bloom ethtypes.Bloom) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixBloom)
	store.Set(evmtypes.BloomKey(height), bloom.Bytes())
}

// GetAllTxLogs return all the transaction logs from the store.
func (k *Keeper) GetAllTxLogs(ctx sdk.Context) []evmtypes.TransactionLogs {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, evmtypes.KeyPrefixLogs)
	defer iterator.Close()

	txsLogs := []evmtypes.TransactionLogs{}
	for ; iterator.Valid(); iterator.Next() {
		hash := common.BytesToHash(iterator.Key())
		var logs []*ethtypes.Log
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &logs)

		// add a new entry
		txLog := evmtypes.NewTransactionLogs(hash, logs)
		txsLogs = append(txsLogs, txLog)
	}
	return txsLogs
}

// GetAccountStorage return state storage associated with an account
func (k *Keeper) GetAccountStorage(ctx sdk.Context, address common.Address) (evmtypes.Storage, error) {
	storage := evmtypes.Storage{}
	err := k.ForEachStorage(ctx, address, func(key, value common.Hash) bool {
		storage = append(storage, evmtypes.NewState(key, value))
		return false
	})
	if err != nil {
		return evmtypes.Storage{}, err
	}

	return storage, nil
}

// GetChainConfig gets block height from block consensus hash
func (k *Keeper) GetChainConfig(ctx sdk.Context) (evmtypes.ChainConfig, bool) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixChainConfig)
	// get from an empty key that's already prefixed by KeyPrefixChainConfig
	bz := store.Get([]byte{})
	if len(bz) == 0 {
		return evmtypes.ChainConfig{}, false
	}

	var config evmtypes.ChainConfig
	k.cdc.MustUnmarshalBinaryBare(bz, &config)
	return config, true
}

// SetChainConfig sets the mapping from block consensus hash to block height
func (k *Keeper) SetChainConfig(ctx sdk.Context, config evmtypes.ChainConfig) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixChainConfig)
	bz := k.cdc.MustMarshalBinaryBare(config)
	// get to an empty key that's already prefixed by KeyPrefixChainConfig
	store.Set([]byte{}, bz)
}

// GetParams returns the total set of evm parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params evmtypes.Params) {
	return k.CommitStateDB.WithContext(ctx).GetParams()
}

// SetParams sets the evm parameters to the param space.
func (k *Keeper) SetParams(ctx sdk.Context, params evmtypes.Params) {
	k.CommitStateDB.WithContext(ctx).SetParams(params)
}

func EndKey(startKey []byte) (endKey []byte){
	key := string(startKey)
	length := len(key)
	last := []rune(key[length-1:])
	end := key[:length-1] + string(last[0] + 1)
	endKey = []byte(end)
	return
}


func IsValidaWasmFile(code []byte) error {
	if !IsWasm(code) {
		return errors.New("it is not a wasm file")
	}else {
		_, err := wasmer.Compile(code)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsWasm checks if the file contents are of wasm binary
func IsWasm(input []byte) bool {
	return bytes.Equal(input[:4], types.WasmIdent)
}
