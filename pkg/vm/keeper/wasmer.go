package keeper

// #include <stdlib.h>
//
// extern int read_db(void*, int, int, int, int, int);
// extern void write_db(void*, int, int, int, int);
// extern void delete_db(void*, int, int);
// extern int new_db_iter(void*, int, int);
// extern int db_iter_next(void*, int);
// extern int db_iter_key(void*, int);
// extern int db_iter_value(void*, int);
//
// extern int send(void*, int, long long);
// extern void get_pre_caller(void*, int);
// extern void get_creator(void*, int);
// extern void get_invoker(void*, int);
// extern void self_address(void*, int);
// extern void get_block_header(void*, int);
//
// extern int get_input_length(void*, int);
// extern void get_input(void*, int, int, int);
// extern void notify_contract(void*, int, int);
// extern void return_contract(void*, int, int);
// extern int call_contract(void*, int, int, int);
// extern void new_contract(void*, int, int, int, int, int);
// extern void destroy_contract(void*);
// extern void panic_contract(void*, int, int);
// extern void get_validator_power(void*, int, int, int);
// extern void total_power(void*, int);
// extern void get_balance(void*, int, int);
// extern void addgas(void*, int);
//
// extern void debug_print(void*, int, int);
import "C"
import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"github.com/wasmerio/go-ext-wasm/wasmer"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"unsafe"
)

var inputData = map[int32][]byte{}

const (
	InputDataTypeParam          = 0
	InputDataTypeContractResult = 1
)

type runtimeConfig struct {
	Store Store
	GasUsed int64
	GasWanted uint64
	PreCaller sdk.AccAddress
	Creator sdk.AccAddress
	Invoker sdk.AccAddress
	SelfAddress sdk.AccAddress
	Keeper *Keeper
	Context *sdk.Context
}


type wasmRuntime struct {

}

//set the store that be used by rust contract.
func(cfg *runtimeConfig) SetStore(kvStore Store) {
	cfg.Store = kvStore
}

func(cfg *runtimeConfig) SetGasUsed(){
	cfg.GasUsed = 0
}

func(cfg *runtimeConfig) SetGasWanted(gaswanted uint64){
	cfg.GasWanted = gaswanted
}

func(cfg *runtimeConfig) SetPreCaller(addr sdk.AccAddress) {
	cfg.PreCaller = addr
}

func(cfg *runtimeConfig) SetCreator(addr sdk.AccAddress) {
	cfg.Creator = addr
}

func(cfg *runtimeConfig) SetInvoker(addr sdk.AccAddress) {
	cfg.Invoker = addr
}

func(cfg *runtimeConfig) SetSelfAddr(addr sdk.AccAddress) {
	cfg.SelfAddress = addr
}

func(cfg *runtimeConfig) SetWasmKeeper(wk *Keeper) {
	cfg.Keeper = wk
}

func(cfg *runtimeConfig) SetCtx(con *sdk.Context) {
	cfg.Context = con
}

//export read_db
func read_db(context unsafe.Pointer, keyPtr, keySize, valuePtr, valueSize, offset int32) int32 {
	return readDB(context, keyPtr, keySize, valuePtr, valueSize, offset)
}

//export write_db
func write_db(context unsafe.Pointer, keyPtr, keySize, valuePtr, valueSize int32) {
	writeDB(context, keyPtr, keySize, valuePtr, valueSize)
}

//export delete_db
func delete_db(context unsafe.Pointer, keyPtr, keySize int32) {
	deleteDB(context, keyPtr, keySize)
}

//export new_db_iter
func new_db_iter(context unsafe.Pointer, prefixPtr, prefixSize int32) int32 {
	return newDBIter(context, prefixPtr, prefixSize)
}

//export db_iter_next
func db_iter_next(context unsafe.Pointer, token int32) int32 {
	return dbIterNext(context, token)
}

//export db_iter_key
func db_iter_key(context unsafe.Pointer, token int32) int32 {
	return dbIterKey(context, token)
}

//export db_iter_value
func db_iter_value(context unsafe.Pointer, token int32) int32 {
	return dbIterValue(context, token)
}

//export send
func send(context unsafe.Pointer, to int32, amount int64) int32 {
	return performSend(context, to, amount)
}

//export get_pre_caller
func get_pre_caller(context unsafe.Pointer, callerPtr int32) {
	getPreCaller(context, callerPtr)
}

//export get_creator
func get_creator(context unsafe.Pointer, creatorPtr int32) {
	getCreator(context, creatorPtr)
}

//export get_invoker
func get_invoker(context unsafe.Pointer, invokerPtr int32) {
	getInvoker(context, invokerPtr)
}

//export self_address
func self_address(context unsafe.Pointer, contractPtr int32) {
	selfAddress(context, contractPtr)
}

//export get_block_header
func get_block_header(context unsafe.Pointer, valuePtr int32) {
	getBlockHeader(context, valuePtr)
}

//export get_input_length
func get_input_length(context unsafe.Pointer, token int32) int32 {
	return getInputLength(context, token)
}

//export get_input
func get_input(context unsafe.Pointer, token, ptr, size int32) {
	getInput(context, token, ptr, size)
}

//export notify_contract
func notify_contract(context unsafe.Pointer, ptr, size int32) {
	notifyContract(context, ptr, size)
}

//export return_contract
func return_contract(context unsafe.Pointer, ptr, size int32) {
	returnContract(context, ptr, size)
}

//export call_contract
func call_contract(context unsafe.Pointer, addrPtr, paramPtr, paramSize int32) int32 {
	return callContract(context, addrPtr, paramPtr, paramSize)
}

//export new_contract
func new_contract(context unsafe.Pointer, codeHashPtr, codeHashSize, argsPtr, argsSize, newContractPtr int32) {
	newContract(context, codeHashPtr, codeHashSize, argsPtr, argsSize, newContractPtr)
}

//export destroy_contract
func destroy_contract(context unsafe.Pointer) {
	destroyContract(context)
}

//export panic_contract
func panic_contract(context unsafe.Pointer, dataPtr, dataSize int32) {
	panicContract(context, dataPtr, dataSize)
}
//export debug_print
func debug_print(context unsafe.Pointer, dataPtr, dataSize int32) {
	debugPrint(context, dataPtr, dataSize)
}
//export get_validator_power
func get_validator_power(context unsafe.Pointer, dataPtr, dataSize, valuePtr int32) {
	getValidatorPower(context, dataPtr, dataSize, valuePtr)
}

//export total_power
func total_power(context unsafe.Pointer, valuePtr int32) {
	totalPower(context, valuePtr)
}

//export get_balance
func get_balance(context unsafe.Pointer, addrPtr, balancePtr int32) {
	getBalance(context, addrPtr, balancePtr)
}

type VMRes struct {
	err []byte // error  response tip
	res []byte // success
}

//export addgas
func addgas(context unsafe.Pointer, gas int32) {
	instanceContext := wasmer.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	runtimeCfg.GasUsed += int64(gas)
	if(uint64(runtimeCfg.GasUsed) > runtimeCfg.GasWanted) {
		panic(sdk.ErrorOutOfGas{Descriptor: "run vm"})
	}
	return
}

type Wasmer struct {
	FilePathMap  map[string]string  `json:"file_path_map"`
	SortMaps     SortMaps 			`json:"sort_maps"`
	LastFileID   int				`json:"last_file_id"`
}

func NewWasmer(homeDir string, _ types.WasmConfig) (*Wasmer, error){
	dir := filepath.Join(homeDir, types.FolderName)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	filePathMap := make(map[string]string)
	LastFileID := 0

	return &Wasmer{
		FilePathMap: filePathMap,
		SortMaps: nil,
		LastFileID:  LastFileID,
	}, nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func (w *Wasmer) Create(homeDir, codeHash string) (Wasmer, error) {
	if fName := w.FilePathMap[codeHash]; fName != "" {
		if FileExist(path.Join(homeDir, WASMDIR, fName)) {
			//file exist, remove file and delete
			err := os.Remove(path.Join(homeDir, WASMDIR, fName))
			if err != nil {
				return Wasmer{}, err
			}
			delete(w.FilePathMap, codeHash)
		}else {
			//file not exist, delete
			delete(w.FilePathMap, codeHash)
		}
	}
	id, fileName := makeFilePath(w.LastFileID)
/*
	err := ioutil.WriteFile(w.HomeDir + "/" + fileName, code, types.ModePerm)
	if err != nil {
		return Wasmer{}, nil, err
	}
*/
	w.FilePathMap[codeHash] = fileName
	newWasmer := Wasmer{
		FilePathMap: nil,
		SortMaps:    sortMapByValue(w.FilePathMap),
		LastFileID:  id,
	}
	return newWasmer, nil
}

func (w wasmRuntime) Call(code []byte, input []byte, method string, cfg *runtimeConfig) (res []byte, err error) {
	instance , err := getInstance(code, cfg)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	method = "x" + hex.EncodeToString([]byte(method))
	call, exist := instance.Exports[method]
	if !exist {
		return nil, errors.New("no expected function")
	}

	inputData[InputDataTypeParam] = input
	defer func() {/**/
		if r := recover(); r != nil{
			switch x := r.(type) {
			case string:
				err = errors.New(x)
				res = nil
			case error:
				err = x
				res = nil
			case sdk.ErrorOutOfGas:
				panic(x)
			case VMRes:
				if x.err == nil {
					res = x.res
					err = nil
				} else {
					res = nil
					err = errors.New(string(x.err))
				}
			default:
				err = errors.New("")
				res = nil
			}
		}
	}()

	_, err2 := call()
	if err2 != nil {
		panic(err2)
	}
	return res, err
}

func getInstance(code []byte, cfg *runtimeConfig) (*wasmer.Instance, error) {
	imports, err := wasmer.NewImports().Namespace("env").Append("send", send, C.send)
	if err != nil {
		panic(err)
	}

	_, _ = imports.Append("read_db", read_db, C.read_db)
	_, _ = imports.Append("write_db", write_db, C.write_db)
	_, _ = imports.Append("delete_db", delete_db, C.delete_db)
	_, _ = imports.Append("new_db_iter", new_db_iter, C.new_db_iter)
	_, _ = imports.Append("db_iter_next", db_iter_next, C.db_iter_next)
	_, _ = imports.Append("db_iter_key", db_iter_key, C.db_iter_key)
	_, _ = imports.Append("db_iter_value", db_iter_value, C.db_iter_value)

	_, _ = imports.Append("get_pre_caller", get_pre_caller, C.get_pre_caller)
	_, _ = imports.Append("self_address", selfAddress, C.self_address)
	_, _ = imports.Append("get_creator", get_creator, C.get_creator)
	_, _ = imports.Append("get_invoker", get_invoker, C.get_invoker)
	_, _ = imports.Append("get_block_header", get_block_header, C.get_block_header)

	_, _ = imports.Append("get_input_length", get_input_length, C.get_input_length)
	_, _ = imports.Append("get_input", get_input, C.get_input)
	_, _ = imports.Append("return_contract", return_contract, C.return_contract)
	_, _ = imports.Append("notify_contract", notify_contract, C.notify_contract)
	_, _ = imports.Append("call_contract", call_contract, C.call_contract)
	_, _ = imports.Append("new_contract", new_contract, C.new_contract)
	_, _ = imports.Append("destroy_contract", destroy_contract, C.destroy_contract)
	_, _ = imports.Append("panic_contract", panic_contract, C.panic_contract)

	_, _ = imports.Append("addgas", addgas, C.addgas)
	_, _ = imports.Append("debug_print", debug_print, C.debug_print)
	_, _ = imports.Append("get_validator_power", get_validator_power, C.get_validator_power)
	_, _ = imports.Append("total_power", total_power, C.total_power)
	_, _ = imports.Append("get_balance", get_balance, C.get_balance)

	module, err := wasmer.Compile(code)
	if err != nil {
		panic(err)
	}
	defer module.Close()

	instance, err := module.InstantiateWithImports(imports)
	if err != nil {
		panic(err)
	}
	instance.SetContextData(cfg)
	return &instance, nil
}


func makeFilePath(id int) (int, string) {
	id ++
	filePath := strconv.Itoa(id) + types.SuffixName
	return id, filePath
}

func MakeCodeHash(code []byte) []byte {
	//get hash
	Md5Inst := md5.New()
	Md5Inst.Write(code)
	Result := Md5Inst.Sum([]byte(""))
	return Result
}

func (w *Wasmer) GetWasmCode(homeDir string, hash []byte) ([]byte, error) {
	Hash := fmt.Sprintf("%x", hash)
	filePath := w.FilePathMap[Hash]
	code, err := ioutil.ReadFile(homeDir + WASMDIR + filePath)
	if err != nil {
		return nil, err
	}
	//the file may be not exist.
	return code, nil
}

func (w *Wasmer) DeleteCode(homeDir string, hash []byte) error {
	Hash := fmt.Sprintf("%x", hash)
	filePath := w.FilePathMap[Hash]
	_, err := os.Lstat(homeDir + WASMDIR + filePath)
	if err == nil {
		err := os.Remove(homeDir + WASMDIR + filePath)
		if err != nil {
			return err
		}
	}
	delete(w.FilePathMap, Hash)

	//splits := strings.Split(filePath, ".")
	//idStr := splits[0]
	//id, _ := strconv.ParseInt(idStr,10,64)
	//files, err := ioutil.ReadDir(w.HomeDir)
	//for index, f := range files {
	//	if int64(index) >= id  {
	//		_ = os.Rename(w.HomeDir + "/"  + f.Name(), w.HomeDir + "/" + fmt.Sprintf("%d.wasm", index))
	//	}
	//}
	//w.LastFileID--

	return nil
}

type SortMap struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type SortMaps []SortMap

func (s SortMaps) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortMaps) Len() int           { return len(s) }
func (s SortMaps) Less(i, j int) bool { return s[i].Value < s[j].Value }

func sortMapByValue(m map[string]string) SortMaps {
	s := make(SortMaps, len(m))
	i := 0
	for k, v := range m {
		s[i] = SortMap{k, v}
		i++
	}
	sort.Sort(s)
	return s
}

func mapFromSortMaps(s SortMaps) map[string]string {
	m := make(map[string]string, len(s))
	for _, v := range s {
		m[v.Key] = v.Value
	}

	return m
}