package keeper

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	types2 "github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/params"
	keeper2 "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/upgrade"
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
	"github.com/wasmerio/wasmer-go/wasmer"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strings"
)

const UINT_MAX uint64 = ^uint64(0)
const INIT string = "init"
const INVOKE string = "invoke"
const CAN_MIGRATE string = "canMigrate"
const CAN_MIGRATE_RESULT string = "true"
type Keeper struct {
	storeKey    		sdk.StoreKey
	cdc         		*codec.Codec
	wasmer     	 		Wasmer
	homeDir				string
	AccountKeeper 		account.AccountKeeper
	StakingKeeper       keeper2.StakingKeeper
	UpgradeKeeper 		upgrade.Keeper
	// Ethermint concrete implementation on the EVM StateDB interface
	CommitStateDB 		*evmtypes.CommitStateDB
	// Transaction counter in a block. Used on StateSB's Prepare function.
	// It is reset to 0 every block on BeginBlock so there's no point in storing the counter
	// on the KVStore or adding it as a field on the EVM genesis state.
	TxCount 			int
	Bloom   			*big.Int
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, homeDir string, paramSpace params.Subspace, accountKeeper account.AccountKeeper, stakingKeeper keeper2.StakingKeeper) Keeper {
	wasmer, err := NewWasmer(homeDir)
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
	}
	return wk
}

func (k Keeper) GetStoreKey() sdk.StoreKey {
	return k.storeKey
}

//
func (k *Keeper) Upload(ctx sdk.Context, wasmCode []byte, creator sdk.AccAddress) (codeHash []byte, err error) {
	ccstore := ctx.KVStore(k.storeKey)
	var wasmer Wasmer
	wasmerBz := ccstore.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "emppty wasmer")
		}
		k.wasmer = Wasmer{
			FilePathMap: mapFromSortMaps(wasmer.SortMaps),
			SortMaps:    nil,
			LastFileID:  wasmer.LastFileID,
		}
	}
	codeHash, err = k.create(ctx, creator, wasmCode)
	if err != nil {
		return nil, err
	}

	return codeHash, nil
}

func (k *Keeper) Instantiate(ctx sdk.Context, codeHash []byte, invoker sdk.AccAddress, args utils.WasmInput, name, version, author, email, describe string, genesisContractAddress sdk.AccAddress, gasWanted uint64) (sdk.AccAddress, error) {
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
	
	if args.Method != InstantiateFuncName {
		return sdk.AccAddress{}, errors.New("Instantiate function must be `init`")
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
		return sdk.AccAddress{}, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid code_hash")
	}
	var contractAddress sdk.AccAddress
	if isGenesis {
		if genesisContractAddress == types.EmptyAddress {
			contractAddress = k.generateContractAddress(codeHash, invoker, args, 1)
		}else {
			contractAddress = genesisContractAddress
		}
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
	_ = contractAccount.SetContractType(types2.WasmContractType)
	k.AccountKeeper.SetAccount(ctx, contractAccount)

	wasmerBz := ccstore.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return sdk.AccAddress{}, sdkerrors.Wrap(sdkerrors.ErrInternal, "empty wasmer")
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

	wasmRuntime := new(wasmRuntime)
	_, err = wasmRuntime.Call(code, args.Sink, args.Method, runtimeCfg)
	if err != nil {
		return sdk.AccAddress{}, err
	}

	initArgs, _ := json.Marshal(args)
	//save the contract info.
	createdAt := types.NewCreatedAt(ctx)
	contractInfo := types.NewContractInfo(codeInfo, initArgs, name, version, author, email, describe, createdAt)
	ccstore.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryBare(contractInfo))

	isGenesis, _ = ctx.Value(types.SystemContract).(bool)
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
func (k *Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, invoker sdk.AccAddress, args utils.WasmInput, gasWanted uint64) (sdk.Result, error) {
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
	wasmRuntime := new(wasmRuntime)
	res, err := wasmRuntime.Call(code, args.Sink, args.Method, runtimeCfg)
	if err != nil {
		return sdk.Result{}, err
	}
	ctx.GasMeter().ConsumeGas(sdk.Gas(runtimeCfg.GasUsed),"wasm cost")
	return sdk.Result{
		Data:   []byte(fmt.Sprintf("%s", string(res))),
	}, nil
}

func (k *Keeper) Migrate(ctx sdk.Context, codeHash []byte, invoker sdk.AccAddress, oldContract sdk.AccAddress, args utils.WasmInput, name, version, author, email, describe string, gasWanted uint64) (sdk.AccAddress, error) {
	canMigrate := utils.WasmInput{
		Method: CAN_MIGRATE,
		Sink:   nil,
	}
	newCtx := ctx
	res, err := k.Query(newCtx, oldContract, invoker, canMigrate)
	if err != nil || res.Result != CAN_MIGRATE_RESULT {
		return sdk.AccAddress{}, errors.New("Cannot migrate")
	}

	newContract, err := k.Instantiate(ctx, codeHash, invoker, args, name, version, author, email, describe, types.EmptyAddress, gasWanted)

	if err != nil {
		return sdk.AccAddress{}, err
	}

	//todo:fix iterator
	oldKey := types.GetContractStorePrefixKey(oldContract)
	oldStore := NewStore(ctx.KVStore(k.storeKey), oldKey)
	iter := oldStore.parent.RemoteIterator(oldKey, sdk.PrefixEndBytes(oldKey))

	defer iter.Close()
	newStoreKey := types.GetContractStorePrefixKey(newContract)
	newStore := NewStore(ctx.KVStore(k.storeKey), newStoreKey)

	for iter.Valid() {
		newStore.Set(iter.Key(), iter.Value())
		iter.Next()
	}

	return newContract, nil
}

// queryContract
func (k Keeper) Query(ctx sdk.Context, contractAddress, invoker sdk.AccAddress, args utils.WasmInput) (types.ContractState, error) {
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
	wasmRuntime := new(wasmRuntime)
	res, err := wasmRuntime.Call(code, args.Sink, args.Method, runtimeCfg)
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
		return types.CodeInfo{}, sdkerrors.Wrap(sdkerrors.ErrInternal, " get contract address failed")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)
	codeHash, _ := hex.DecodeString(contract.CodeInfo.CodeHash)
	bz := ccstore.Get(types.GetCodeKey(codeHash))
	if bz == nil {
		return types.CodeInfo{}, sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("invalid code_hash: %v", codeHash))
	}
	wasmerBz := ccstore.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return types.CodeInfo{}, sdkerrors.Wrap(sdkerrors.ErrInternal, "unexpected wasmer info")
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

func (k *Keeper) create(ctx sdk.Context, invokerAddr sdk.AccAddress, wasmCode []byte) (codeHash []byte, err error) {
	wasmCode, err = UnCompress(wasmCode)
	if err != nil {
		return nil, err
	}
	//checks if the file contents are of wasm binary
	ok := IsValidaWasmFile(wasmCode)
	if ok != nil {
		return nil, ok
	}
	ccstore := ctx.KVStore(k.storeKey)
	codeHash = MakeCodeHash(wasmCode)

	// addgas
	gasedCode, err := tryAddgas(wasmCode)
	if err != nil {
		return nil, err
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
				return nil, err
			}
			localFileHash := MakeCodeHash(localCode)
			//the content if different, delete local file and save remote file.
			if !bytes.Equal(localFileHash, codeHash) {
				err = os.Remove(filePath)
				if err != nil {
					return nil, err
				}
				err = ioutil.WriteFile(filePath, gasedCode, types.ModePerm)
				if err != nil {
					return nil, err
				}
			}
			return codeHash, nil
		}else {
			err = ioutil.WriteFile(filePath, gasedCode, types.ModePerm)
			if err != nil {
				return nil, err
			}
			return codeHash, nil
		}
	}
	newWasmer, err := k.wasmer.Create(k.homeDir, fmt.Sprintf("%x", codeHash))
	if err != nil {
		return nil, err
	}
	bz := k.cdc.MustMarshalJSON(newWasmer)
	if bz == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "cdc marshal failed")
	}
	ccstore.Set(types.GetWasmerKey(), bz)
	ccstore.Set(codeHash, wasmCode)
	codeInfo := types.NewCodeInfo(strings.ToUpper(hex.EncodeToString(codeHash)), invokerAddr)
	ccstore.Set(types.GetCodeKey(codeHash), k.cdc.MustMarshalBinaryBare(codeInfo))
	hash := fmt.Sprintf("%x", codeHash)
	filePath := path.Join(k.homeDir, WASMDIR, k.wasmer.FilePathMap[hash])
	err = ioutil.WriteFile(filePath, gasedCode, types.ModePerm)
	if err != nil {
		return nil, err
	}
	return codeHash, nil
}

func (k Keeper) generateContractAddress(codeHash []byte, creatorAddr sdk.AccAddress, payload utils.WasmInput, nonce uint64) sdk.AccAddress {
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
		// Create an Engine
		engine := wasmer.NewEngine()
		store := wasmer.NewStore(engine)
		_, err := wasmer.NewModule(store, code)
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


func (k *Keeper) RecordSection(ctx sdk.Context, height int64, bloom ethtypes.Bloom) {
	index := (height-1)/evmtypes.SectionSize

	gen, found := k.GetSectionBloom(ctx, index)
	if !found {
		gen, _ = NewGenerator(evmtypes.SectionSize)
	}
	gen.AddBloom(uint((height-1) % evmtypes.SectionSize), bloom)

	k.SetSectionBloom(ctx, index, gen)
}

func (k *Keeper) GetSectionBloom(ctx sdk.Context, index int64) (*Generator, bool) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixSection)

	//has := store.Has(evmtypes.BloomKey(index))
	//if !has {
	//	return nil, false
	//}

	bz := store.Get(evmtypes.BloomKey(index))
	if len(bz) == 0 {
		return nil, false
	}
	var section Generator
	err := json.Unmarshal(bz, &section)
	if err != nil {
		return nil, false
	}
	return &section, true
}

func (k *Keeper) GetSectionBlooms(ctx sdk.Context, start, end uint64) (map[uint64]*Generator, error) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixSection)

	iterator := store.Iterator(evmtypes.BloomKey(int64(start)), evmtypes.BloomKey(int64(end+1)))
	defer iterator.Close()

	ges := map[uint64]*Generator{}

	for ; iterator.Valid(); iterator.Next() {
		if iterator.Error() != nil {
			return nil, iterator.Error()
		}
		key := evmtypes.BloomKeyFromByte(iterator.Key()[len(evmtypes.KeyPrefixSection):])

		var ge Generator
		err := json.Unmarshal(iterator.Value(), &ge)
		if err != nil {
			return nil, err
		}
		ges[key] = &ge
	}
	return ges, nil
}

func (k *Keeper) SetSectionBloom(ctx sdk.Context, index int64, gen *Generator) {
	store := NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixSection)
	by, _ := json.Marshal(gen)
	store.Set(evmtypes.BloomKey(index), by)
}

var (
	// errSectionOutOfBounds is returned if the user tried to add more bloom filters
	// to the batch than available space, or if tries to retrieve above the capacity.
	errSectionOutOfBounds = evmtypes.ErrSectionOutOfBounds

	// errBloomBitOutOfBounds is returned if the user tried to retrieve specified
	// bit bloom above the capacity.
	errBloomBitOutOfBounds = evmtypes.ErrBloomBitOutOfBounds
)

// Generator takes a number of bloom filters and generates the rotated bloom bits
// to be used for batched filtering.
type Generator struct {
	Blooms   [ethtypes.BloomBitLength][]byte `json:"blooms"`// Rotated blooms for per-bit matching
	Sections uint  `json:"sections"`                       // Number of sections to batch together
	//NextSec  uint   `json:"next_sec"`                     // Next section to set when adding a bloom
}

// NewGenerator creates a rotated bloom generator that can iteratively fill a
// batched bloom filter's bits.
func NewGenerator(sections uint) (*Generator, error) {
	if sections%8 != 0 {
		return nil, evmtypes.ErrBloomFilterSectionNum
	}
	b := &Generator{Sections: sections}
	for i := 0; i < ethtypes.BloomBitLength; i++ {
		b.Blooms[i] = make([]byte, sections/8)
	}
	return b, nil
}

// AddBloom takes a single bloom filter and sets the corresponding bit column
// in memory accordingly.
func (b *Generator) AddBloom(index uint, bloom ethtypes.Bloom) error {
	// Make sure we're not adding more bloom filters than our capacity
	if index >= b.Sections {
		return errSectionOutOfBounds
	}
	//if b.NextSec != index {
	//	return errors.New("bloom filter with unexpected index")
	//}
	// Rotate the bloom and insert into our collection
	byteIndex := index / 8
	bitMask := byte(1) << byte(7-index%8)

	for i := 0; i < ethtypes.BloomBitLength; i++ {
		bloomByteIndex := ethtypes.BloomByteLength - 1 - i/8
		bloomBitMask := byte(1) << byte(i%8)

		if (bloom[bloomByteIndex] & bloomBitMask) != 0 {
			b.Blooms[i][byteIndex] |= bitMask
		}
	}
	//b.NextSec++

	return nil
}

// Bitset returns the bit vector belonging to the given bit index after all
// blooms have been added.
func (b *Generator) Bitset(idx uint) ([]byte, error) {
	//if b.NextSec != b.Sections {
	//	return nil, errors.New("bloom not fully generated yet")
	//}
	if idx >= ethtypes.BloomBitLength {
		return nil, errBloomBitOutOfBounds
	}
	return b.Blooms[idx], nil
}

