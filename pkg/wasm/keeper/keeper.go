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
	keeper2 "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	dbm "github.com/tendermint/tm-db"
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
		wasmer:        wasmer,
		AccountKeeper: accountKeeper,
		StakingKeeper: stakingKeeper,
		cdb:		   cdb,
	}
	SetAccountKeeper(accountKeeper)
	SetWasmKeeper(&wk)
	SetStakingKeeper(stakingKeeper)
	return wk
}

//
func (k Keeper) Instantiate(ctx sdk.Context, wasmCode []byte, invoker sdk.AccAddress, nonce uint64, args json.RawMessage, name, version, author, email, describe string, isGenesis bool, genesisContractAddress sdk.AccAddress) (sdk.AccAddress, error) {
	codeHash, isExist, err := k.create(ctx, invoker, wasmCode)
	if err != nil {
		return sdk.AccAddress{}, err
	}

	SetGasUsed()
	SetCtx(&ctx)
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
	var contractAddress sdk.AccAddress
	if isGenesis {
		contractAddress = genesisContractAddress
	}else {
		contractAddress = k.generateContractAddress(invoker, nonce)
	}
	existingAcct := k.AccountKeeper.GetAccount(ctx, contractAddress)
	if existingAcct != nil {
		return sdk.AccAddress{}, sdk.ErrInternal("Contract account exists")
	}
	SetPreCaller(invoker)
	SetInvoker(invoker)
	SetCreator(invoker)
	SetSelfAddr(contractAddress)
	var contractAccount exported.Account
	contractAccount = k.AccountKeeper.NewAccountWithAddress(ctx, contractAddress)
	k.AccountKeeper.SetAccount(ctx, contractAccount)

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCodeKey(codeHash))
	if bz == nil {
		return sdk.AccAddress{}, sdk.ErrInternal("codeHash not found")
	}
	wasmerBz := store.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return sdk.AccAddress{}, sdk.ErrInternal("empty wasmer")
		}
		k.wasmer = wasmer
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &codeInfo)

	wc, err := k.wasmer.GetWasmCode(codeHash)
	if err != nil {
		wc = store.Get(codeHash)

		fileName := k.wasmer.FilePathMap[strings.ToLower(codeInfo.CodeHash)]
		err = ioutil.WriteFile(k.wasmer.HomeDir + "/" + fileName, wc, types.ModePerm)
		if err != nil {
			return sdk.AccAddress{}, err
		}
	}
	code = wc
	//create store

	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	SetStore(prefixStore)
	if len(args) != 0 {
		input, err := handleArgs(args)
		if err != nil {
			return sdk.AccAddress{}, err
		}
		_, err = k.wasmer.Call(code, input, INIT)
		if err != nil {
			return sdk.AccAddress{}, err
		}
	}

	//save the contract info.
	createdAt := types.NewCreatedAt(ctx)
	contractInfo := types.NewContractInfo(codeHash, invoker, args, name, version, author, email, describe, createdAt)
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryBare(contractInfo))

	//save contractAddress into account
	if !isGenesis {
		contractAddrStr := contractAddress.String()
		accountAddr := k.AccountKeeper.GetAccount(ctx, invoker).GetAddress()

		var contractList []string
		contractListBytes := store.Get(types.GetAccountContractListKey(accountAddr))
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
		store.Set(types.GetAccountContractListKey(accountAddr), contractListBytes)
	}
	ctx.GasMeter().ConsumeGas(sdk.Gas(GasUsed),"wasm cost")

	//store code in local.
	if !isExist {
		hash := fmt.Sprintf("%x", codeHash)
		err = ioutil.WriteFile(k.wasmer.HomeDir + "/" + k.wasmer.FilePathMap[hash], wasmCode, types.ModePerm)
		if err != nil {
			return sdk.AccAddress{}, err
		}
	}
	return contractAddress, nil
}

//
func (k Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, invoker sdk.AccAddress, args json.RawMessage) (sdk.Result, error) {
	SetGasUsed()
	SetSelfAddr(contractAddress)
	SetInvoker(invoker)
	SetPreCaller(invoker)
	SetCtx(&ctx)

	contract := k.GetContractInfo(ctx, contractAddress)
	SetCreator(contract.CodeInfo.Creator)
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
	wc, err := k.wasmer.GetWasmCode(codeHash)
	ccstore := ctx.KVStore(k.storeKey)
	if err != nil {
		wc = ccstore.Get(codeHash)

		fileName := k.wasmer.FilePathMap[strings.ToLower(codeInfo.CodeHash)]
		err = ioutil.WriteFile(k.wasmer.HomeDir + "/" + fileName, wc, types.ModePerm)
		if err != nil {
			return sdk.Result{}, err
		}
	}
	code = wc
	//get store
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	SetStore(prefixStore)
	input, err := handleArgs(args)
	if err != nil {
		return sdk.Result{}, err
	}
	res, err := k.wasmer.Call(code, input, INVOKE)
	if err != nil {
		return sdk.Result{}, err
	}
	ctx.GasMeter().ConsumeGas(sdk.Gas(GasUsed),"wasm cost")
	return sdk.Result{
		Data:   []byte(fmt.Sprintf("%s", string(res))),
	}, nil
}

func (k Keeper) Migrate(ctx sdk.Context, wasmCode []byte, invoker sdk.AccAddress, oldContract sdk.AccAddress, nonce uint64, args json.RawMessage, name, version, author, email, describe string) (sdk.AccAddress, error) {
	newContract, err := k.Instantiate(ctx, wasmCode, invoker, nonce, args, name, version, author, email, describe, false, types.EmptyAddress)
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

// query?
func (k Keeper) Query(ctx sdk.Context, contractAddress, invokerAddress sdk.AccAddress, args json.RawMessage) (types.ContractState, error) {
	SetCreator(k.GetCreator(ctx, contractAddress))
	SetPreCaller(invokerAddress)
	SetInvoker(invokerAddress)
	SetSelfAddr(contractAddress)
	SetCtx(&ctx)
	SetGasUsed()
	SetGasWanted(UINT_MAX)
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
	wc, err := k.wasmer.GetWasmCode(codeHash)
	if err != nil {
		wc = store.Get(codeHash)

		fileName := k.wasmer.FilePathMap[strings.ToLower(codeInfo.CodeHash)]
		err = ioutil.WriteFile(k.wasmer.HomeDir + "/" + fileName, wc, types.ModePerm)
		if err != nil {
			return types.ContractState{}, err
		}
	}
	code = wc

	//get store
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	SetStore(prefixStore)
	input, err := handleArgs(args)
	if err != nil {
		return types.ContractState{}, err
	}
	res, err := k.wasmer.Call(code, input, INVOKE)
	if err != nil {
		return types.ContractState{}, err
	}

	contractState := types.ContractState{Result: string(res)}

	return contractState, nil
}


func (k *Keeper) contractInstance(ctx sdk.Context, contractAddress sdk.AccAddress) (types.CodeInfo, error) {

	var wasmer Wasmer
	store := ctx.KVStore(k.storeKey)
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return types.CodeInfo{}, sdk.ErrInternal(" get contract address failed")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)
	codeHash, _ := hex.DecodeString(contract.CodeInfo.CodeHash)
	bz := store.Get(types.GetCodeKey(codeHash))
	if bz == nil {
		return types.CodeInfo{}, sdk.ErrInternal("get code key failed")
	}
	wasmerBz := store.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return types.CodeInfo{}, sdk.ErrInternal("unexpected wasmer info")
		}
		k.wasmer = wasmer
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

func (k Keeper) create(ctx sdk.Context, invokerAddr sdk.AccAddress, wasmCode []byte) (codeHash []byte, isExist bool, err error) {
	wasmCode, err = UnCompress(wasmCode)
	if err != nil {
		return nil, false, err
	}
	//checks if the file contents are of wasm binary
	ok := types.IsValidaWasmFile(wasmCode)
	if ok != nil {
		return nil, false, ok
	}
	// addgas
	wasmCode, err = tryAddgas(wasmCode)
	if err != nil {
		return nil, false, err
	}

	ccstore := ctx.KVStore(k.storeKey)
	var wasmer Wasmer
	wasmerBz := ccstore.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return nil, false, sdk.ErrInternal("empty wasmer")
		}
		k.wasmer = wasmer
	}
	codeHash = MakeCodeHash(wasmCode)
	//check if it has been saved in couchDB.
	codeByte := ccstore.Get(codeHash)
	if codeByte != nil {
		hash := fmt.Sprintf("%x", codeHash)
		filePath := path.Join(k.wasmer.HomeDir, k.wasmer.FilePathMap[hash])
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
	newWasmer, codeHash, err := k.wasmer.Create(wasmCode)
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

func (k Keeper) generateContractAddress(creatorAddr sdk.AccAddress, nonce uint64) sdk.AccAddress {
	contract, _ := rlp.EncodeToBytes([]interface{}{creatorAddr, nonce})
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