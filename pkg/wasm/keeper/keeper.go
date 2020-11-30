package keeper

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/logger"
	keeper2 "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	dbm "github.com/tendermint/tm-db"
	"github.com/wasmerio/go-ext-wasm/wasmer"
	"io/ioutil"
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
	cdb					dbm.DB
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, homeDir string, wasmConfig types.WasmConfig,  accountKeeper account.AccountKeeper, stakingKeeper keeper2.StakingKeeper, cdb dbm.DB) Keeper {
	wasmer, err := NewWasmer(homeDir, wasmConfig)
	if err != nil {
		panic(err)
	}

	wk := Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		wasmer:        *wasmer,
		homeDir:       homeDir,
		AccountKeeper: accountKeeper,
		StakingKeeper: stakingKeeper,
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

func (k *Keeper) Instantiate(ctx sdk.Context, codeHash []byte, invoker sdk.AccAddress, args json.RawMessage, name, version, author, email, describe string, genesisContractAddress sdk.AccAddress, gasWanted uint64) (sdk.AccAddress, error) {
	runtimeCfg := &runtimeConfig{
		GasUsed: 0,
		PreCaller:   invoker,
		Invoker:     invoker,
		Creator:     invoker,
		SelfAddress: sdk.AccAddress{},
		Keeper:      k,
		Context:     &ctx,
	}

	// 如果是官方合约，不限制gas数量
	isGenesis, ok := ctx.Value(types.SystemContract).(bool)
	if ok && isGenesis {
		runtimeCfg.SetGasWanted(UINT_MAX)
	} else {
		runtimeCfg.SetGasWanted(gasWanted)
	}

	var codeInfo types.CodeInfo
	var wasmer Wasmer
	var code []byte
	var params types.CallContractParam
	if args != nil {
		err := json.Unmarshal(args, &params)
		if err != nil {
			return sdk.AccAddress{}, sdk.ErrInternal("invalid instantiate message")
		}
	}
	ccstore := ctx.KVStore(k.storeKey)
	bz := ccstore.Get(types.GetCodeKey(codeHash))
	if bz == nil {
		return sdk.AccAddress{}, sdk.ErrInternal("codeHash not found")
	}
	var contractAddress sdk.AccAddress
	if isGenesis {
		contractAddress = genesisContractAddress
	}else {
		contractAddress = k.generateContractAddress(codeHash, invoker, args)
	}
	existingAcct := k.AccountKeeper.GetAccount(ctx, contractAddress)
	if existingAcct != nil {
		//return sdk.AccAddress{}, sdk.ErrInternal("Contract account exists")
		return contractAddress, nil
	}
	runtimeCfg.SetSelfAddr(contractAddress)
	var contractAccount exported.Account
	contractAccount = k.AccountKeeper.NewAccountWithAddress(ctx, contractAddress)
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
	if len(args) != 0 {
		input, err := handleArgs(args)
		if err != nil {
			return sdk.AccAddress{}, err
		}
		_, err = k.wasmer.Call(code, input, INIT, runtimeCfg)
		if err != nil {
			return sdk.AccAddress{}, err
		}
	}

	//save the contract info.
	createdAt := types.NewCreatedAt(ctx)
	contractInfo := types.NewContractInfo(codeInfo, args, name, version, author, email, describe, createdAt)
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
func (k *Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, invoker sdk.AccAddress, args json.RawMessage, gasWanted uint64) (sdk.Result, error) {
	logger.GetLogger().Info(("!!!!!!!!!!!!!Enter Execute!!!!!!!!!!!!! "))
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
	var params types.CallContractParam
	if args != nil {
		err := json.Unmarshal(args, &params)
		if err != nil {
			return sdk.Result{}, sdk.ErrInternal("invalid handle message")
		}
	}
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
	input, err := handleArgs(args)
	if err != nil {
		return sdk.Result{}, err
	}
	res, err := k.wasmer.Call(code, input, INVOKE, runtimeCfg)
	if err != nil {
		return sdk.Result{}, err
	}
	ctx.GasMeter().ConsumeGas(sdk.Gas(runtimeCfg.GasUsed), "wasm cost")
	return sdk.Result{
		Data:   []byte(fmt.Sprintf("%s", string(res))),
	}, nil
}

func (k *Keeper) Migrate(ctx sdk.Context, codeHash []byte, invoker sdk.AccAddress, oldContract sdk.AccAddress, args json.RawMessage, name, version, author, email, describe string, gasWanted uint64) (sdk.AccAddress, error) {
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
func (k Keeper) Query(ctx sdk.Context, contractAddress, invoker sdk.AccAddress, args json.RawMessage) (types.ContractState, error) {
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

	var params types.CallContractParam
	if args != nil {
		err := json.Unmarshal(args, &params)
		if err != nil {
			return types.ContractState{}, sdk.ErrInternal("invalid query message")
		}
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
	input, err := handleArgs(args)
	if err != nil {
		return types.ContractState{}, err
	}
	res, err := k.wasmer.Call(code, input, INVOKE, runtimeCfg)
	if err != nil {
		return types.ContractState{}, err
	}

	contractState := types.ContractState{Result: string(res)}

	return contractState, nil
}


func (k *Keeper) contractInstance(ctx sdk.Context, contractAddress sdk.AccAddress) (types.CodeInfo, error) {
	logger.GetLogger().Info(("!!!!!!!!!!!!!Enter contractInstance!!!!!!!!!!!!! "))
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

func (k Keeper) generateContractAddress(codeHash []byte, creatorAddr sdk.AccAddress, payload json.RawMessage) sdk.AccAddress {
	contract, _ := rlp.EncodeToBytes([]interface{}{codeHash, creatorAddr, payload})
	return sdk.ToAccAddress(crypto.Keccak256Hash(contract).Bytes()[12:])
}

func handleArgs(args json.RawMessage) ([]byte, error){
	var param types.CallContractParam
	inputByte, _ := args.MarshalJSON()
	err := json.Unmarshal(inputByte, &param)
	if err != nil {
		return nil, err
	}

	var inputArgs []interface{}
	for i := 0; i < len(param.Args); i++ {
		inputArgs = append(inputArgs, param.Args[i])
	}

	input := Serialize(inputArgs)
	return input, nil
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
