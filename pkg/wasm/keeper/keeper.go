package keeper

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	dbm "github.com/tendermint/tm-db"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const UINT_MAX uint64 = ^uint64(0)
type Keeper struct {
	storeKey    		sdk.StoreKey
	cdc         		*codec.Codec
	wasmer     	 		Wasmer
	AccountKeeper 		account.AccountKeeper
	cdb					dbm.DB
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, homeDir string, wasmConfig types.WasmConfig,  accountKeeper account.AccountKeeper, cdb dbm.DB) Keeper {
	wasmer, err := NewWasmer(homeDir, wasmConfig)
	if err != nil {
		panic(err)
	}

	wk := Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		wasmer:        wasmer,
		AccountKeeper: accountKeeper,
		cdb:		   cdb,
	}
	SetAccountKeeper(accountKeeper)
	SetWasmKeeper(&wk)
	return wk
}

//　Install uploads and compiles a WASM contract, returning a short identifier for the contract
func (k Keeper) Install(ctx sdk.Context, invokerAddr sdk.AccAddress, wasmCode []byte) (codeHash []byte, err error) {
	wasmCode, err = uncompress(wasmCode)
	if err != nil {
		return nil, err
	}
	//checks if the file contents are of wasm binary
	ok := types.IsValidaWasmFile(wasmCode)
	if ok != nil {
		return nil, ok
	}
	// addgas
	wasmCode, err = tryAddgas(wasmCode)
	if err != nil {
		return nil, err
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
	codeHash = MakeCodeHash(wasmCode)
	//check if it has been saved in couchDB.
	codeByte := store.Get(codeHash)
	if codeByte != nil {
		hash := fmt.Sprintf("%x", codeHash)
		filePath := path.Join(k.wasmer.HomeDir, k.wasmer.FilePathMap[hash])
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
			}
			err = ioutil.WriteFile(filePath, wasmCode, types.ModePerm)
			if err != nil {
				return nil, err
			}
			return codeHash, nil
		}else {
			err = ioutil.WriteFile(filePath, wasmCode, types.ModePerm)
			if err != nil {
				return nil, err
			}
			return codeHash,nil
		}
	}
	newWasmer, codeHash, err := k.wasmer.Create(wasmCode)
	if err != nil {
		return codeHash, err
	}
	bz := k.cdc.MustMarshalJSON(newWasmer)
	if bz == nil {
		return nil, sdk.ErrInternal("marshal json failed")
	}
	store.Set(codeHash, wasmCode)
	codeInfo := types.NewCodeInfo(strings.ToUpper(hex.EncodeToString(codeHash)), invokerAddr)
	store.Set(types.GetCodeKey(codeHash), k.cdc.MustMarshalBinaryBare(codeInfo))
	store.Set(types.GetWasmerKey(), bz)

	//store code in local.
	hash := fmt.Sprintf("%x", codeHash)
	err = ioutil.WriteFile(newWasmer.HomeDir + "/" + newWasmer.FilePathMap[hash], wasmCode, types.ModePerm)
	if err != nil {
		return nil, err
	}

	return codeHash, nil
}

//　Uninstall remove a WASM contract, returning an error or nil
func (k Keeper) Uninstall(ctx sdk.Context, invokerAddr sdk.AccAddress, codeHash []byte) error {
	store := ctx.KVStore(k.storeKey)
	codeByte := store.Get(codeHash)
	if codeByte == nil {
		return errors.New("uninstall failed: error codeHash")
	}

	store.Delete(types.GetCodeKey(codeHash))

	var wasmer Wasmer
	wasmerBz := store.Get(types.GetWasmerKey())
	if wasmerBz != nil {
		k.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
			if wasmer.LastFileID == 0 {
				panic(errors.New("unexpected wasmer info"))
			}
			k.wasmer = wasmer
			err := k.wasmer.DeleteCode(codeHash)
			if err != nil {
				panic(err)
			}
			bz := k.cdc.MustMarshalJSON(k.wasmer)
			store.Set(types.GetWasmerKey(), bz)
		} else {
			panic(errors.New("no wasmer"))
	}
	return nil
}

//
func (k Keeper) Instantiate(ctx sdk.Context, codeHash []byte, invoker sdk.AccAddress, args json.RawMessage, name, version, author, email, describe string) (sdk.AccAddress, error) {
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
	contractAddress := k.generateContractAddress(codeHash)
	existingAcct := k.AccountKeeper.GetAccount(ctx, contractAddress)
	if existingAcct != nil {
		return sdk.AccAddress{}, sdk.ErrInternal("Contract account exists")
	}
	SetInvoker(invoker)
	SetCreator(contractAddress)
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
	if len(args) != 0 {
		input, err := handleArgs(args)
		if err != nil {
			return sdk.AccAddress{}, err
		}
		_, err = k.wasmer.Call(code, input)
		if err != nil {
			return sdk.AccAddress{}, err
		}
	}
	//save the contract info.
	createdAt := types.NewCreatedAt(ctx)
	contractInfo := types.NewContractInfo(codeHash, invoker, args, name, version, author, email, describe, createdAt)
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryBare(contractInfo))
	//save contractAddress into account
	Account := k.AccountKeeper.GetAccount(ctx, invoker)
	Account.AddContract(contractAddress)
	k.AccountKeeper.SetAccount(ctx, Account)
	ctx.GasMeter().ConsumeGas(sdk.Gas(GasUsed),"wasm cost")
	return contractAddress, nil
}

//
func (k Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, invoker sdk.AccAddress, args json.RawMessage) (sdk.Result, error) {
	SetGasUsed()
	SetCreator(contractAddress)
	SetInvoker(invoker)
	SetCtx(&ctx)
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
	input, err := handleArgs(args)
	if err != nil {
		return sdk.Result{}, err
	}
	res, err := k.wasmer.Call(code, input)
	if err != nil {
		return sdk.Result{}, err
	}
	ctx.GasMeter().ConsumeGas(sdk.Gas(GasUsed),"wasm cost")
	return sdk.Result{
		Data:   []byte(fmt.Sprintf("%s", string(res))),
	}, nil
}

// query?
func (k Keeper) Query(ctx sdk.Context, contractAddress sdk.AccAddress, args json.RawMessage) (types.ContractState, error) {
	SetCreator(contractAddress)
	SetInvoker(sdk.AccAddress{})
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
	input, err := handleArgs(args)
	if err != nil {
		return types.ContractState{}, err
	}
	res, err := k.wasmer.Call(code, input)
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

func (k Keeper) generateContractAddress(codeHash []byte) sdk.AccAddress {
	//fmt.Println(sdk.ToAccAddress(codeHash))
	return sdk.ToAccAddress(codeHash)
}

func handleArgs(args json.RawMessage) ([]byte, error){
	var param Param
	inputByte, _ := args.MarshalJSON()
	err := json.Unmarshal(inputByte, &param)
	if err != nil {
		return nil, err
	}

	inputArgs := []interface{}{param.Method}
	for i := 0; i < len(param.Args); i++ {
		inputArgs = append(inputArgs, param.Args[i])
	}

	input := Serialize(inputArgs)
	return input, nil
}

