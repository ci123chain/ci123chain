package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	"github.com/tanhuiya/ci123chain/pkg/account/exported"
	"github.com/tanhuiya/ci123chain/pkg/wasm/types"
	"io/ioutil"
	"strings"
)

type Keeper struct {
	storeKey    sdk.StoreKey
	cdc         *codec.Codec
	wasmer      Wasmer
	AccountKeeper       account.AccountKeeper
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, homeDir string, wasmConfig types.WasmConfig,  accountKeeper account.AccountKeeper) Keeper {
	wasmer, err := NewWasmer(homeDir, wasmConfig)
	if err != nil {
		panic(err)
	}
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		wasmer:        wasmer,
		AccountKeeper: accountKeeper,
	}
}

//　Create uploads and compiles a WASM contract, returning a short identifier for the contract
func (k Keeper) Create(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte) (codeHash []byte, err error) {

	wasmCode, err = uncompress(wasmCode)
	if err != nil {
		return nil, err
	}
	//checks if the file contents are of wasm binary
	ok := types.IsValidaWasmFile(wasmCode)
	if ok != nil {
		return nil, ok
	}
	store := ctx.KVStore(k.storeKey)
	var wasmer Wasmer
	wasmerBz := store.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			return nil, sdk.ErrInternal("empty wasmer")
		}
		k.wasmer = wasmer
	}
	newWasmer, codeHash, err := k.wasmer.Create(wasmCode)
	if err != nil {
		return nil, err
	}
	bz := k.cdc.MustMarshalJSON(newWasmer)
	if bz == nil {
		return nil, sdk.ErrInternal("marshal json failed")
	}
	codeInfo := types.NewCodeInfo(strings.ToUpper(hex.EncodeToString(codeHash)), creator)
	store.Set(types.GetCodeKey(codeHash), k.cdc.MustMarshalBinaryBare(codeInfo))
	store.Set(types.GetWasmerKey(), bz)
	store.Set(codeHash, wasmCode)
	return codeHash, nil
}

//
func (k Keeper) Instantiate(ctx sdk.Context, codeHash []byte, creator sdk.AccAddress, args json.RawMessage, label string) (sdk.AccAddress, error) {
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
	contractAddress := k.generateContractAddress(codeHash)
	existingAcct := k.AccountKeeper.GetAccount(ctx, contractAddress)
	if existingAcct != nil {
		return sdk.AccAddress{}, sdk.ErrInternal("account exists")
	}
	var contractAccount exported.Account
	/*if !deposit.IsZero() {
		sdkerr := k.AccountKeeper.Transfer(ctx, creator, contractAddress, deposit)
		if sdkerr != nil {
			return sdk.AccAddress{}, sdk.ErrInternal("transfer failed")
		}
	}else {
		contractAccount = k.AccountKeeper.NewAccountWithAddress(ctx, contractAddress)
		k.AccountKeeper.SetAccount(ctx, contractAccount)
	}*/
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

		fileName := k.wasmer.FilePathMap[fmt.Sprintf("%x",codeInfo.CodeHash)]
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

	_, err = k.wasmer.Instantiate(code,types.InitFunctionName, args)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	//save the contract info.
	createdAt := types.NewCreatedAt(ctx)
	contractInfo := types.NewContractInfo(codeHash, creator, args, label, createdAt)
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryBare(contractInfo))
	//save contractAddress into account
	Account := k.AccountKeeper.GetAccount(ctx, creator)
	Account.AddContract(contractAddress)
	k.AccountKeeper.SetAccount(ctx, Account)
	return contractAddress, nil
}

//
func (k Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, args json.RawMessage) (sdk.Result, error) {

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
	store := ctx.KVStore(k.storeKey)
	var code []byte
	codeHash, _ := hex.DecodeString(codeInfo.CodeHash)
	wc, err := k.wasmer.GetWasmCode(codeHash)
	if err != nil {
		wc = store.Get(codeHash)

		fileName := k.wasmer.FilePathMap[fmt.Sprintf("%x",codeInfo.CodeHash)]
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

	res, err := k.wasmer.Execute(code, types.HandleFunctionName, args)
	if err != nil {
		return sdk.Result{}, err
	}
	return sdk.Result{
		Data:   []byte(fmt.Sprintf("%s", res)),
	}, nil
}

// query?
func (k Keeper) Query(ctx sdk.Context, contractAddress sdk.AccAddress, msg json.RawMessage) (types.ContractState, error) {

	var params types.CallContractParam
	if msg != nil {
		err := json.Unmarshal(msg, &params)
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

		fileName := k.wasmer.FilePathMap[fmt.Sprintf("%x",codeInfo.CodeHash)]
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

	res, err := k.wasmer.Query(code, types.QueryFunctionName, msg)
	if err != nil {
		return types.ContractState{}, err
	}
	contractState := types.ContractState{Result:res}
	return contractState, nil
}


func (k *Keeper) contractInstance(ctx sdk.Context, contractAddress sdk.AccAddress) (types.CodeInfo, error) {

	var wasmer Wasmer
	store := ctx.KVStore(k.storeKey)
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return types.CodeInfo{}, sdk.ErrInternal("get contract address failed")
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

func (k Keeper) generateContractAddress(codeHash []byte) sdk.AccAddress {
	fmt.Println(sdk.ToAccAddress(codeHash))
	return sdk.ToAccAddress(codeHash)
}