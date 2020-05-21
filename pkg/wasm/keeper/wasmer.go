package keeper

// #include <stdlib.h>
//
// extern int read_db(void *context, int key, int value);
// extern void write_db(void *context, int key, int value);
// extern void delete_db(void *context, int key);
// extern int send(void *context, int toPtr, int amountPtr);
// extern void get_creator(void *context, int creatorPtr);
// extern void get_invoker(void *context, int invokerPtr);
// extern void get_time(void *context, int timePtr);
// extern void addgas(void *context, int gas);
import "C"
import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wasmerio/go-ext-wasm/wasmer"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"unsafe"
)

var GasUsed int64
var GasWanted uint64
func SetGasUsed(){
	GasUsed = 0
}

func SetGasWanted(gaswanted uint64){
	GasWanted = gaswanted
}

//export addgas
func addgas(context unsafe.Pointer, gas int32) {
	GasUsed += int64(gas)
	if(uint64(GasUsed) > GasWanted) {
		panic("out of gas in location: vm")
	}
	return
}

var blockHeader abci.Header
func SetBlockHeader(header abci.Header) {
	blockHeader = header
}

var creator string
func SetCreator(addr string) {
	creator = addr
}

var invoker string
func SetInvoker(addr string) {
	invoker = addr
}

var accountKeeper account.AccountKeeper
func SetAccountKeeper(ac account.AccountKeeper) {
	accountKeeper = ac
}

var ctx *sdk.Context
func SetCtx(con *sdk.Context) {
	ctx = con
}

//export send
func send(context unsafe.Pointer, toPtr int32, amountPtr int32) int32{
	return perform_send(context, toPtr, amountPtr)
}

//export get_creator
func get_creator(context unsafe.Pointer, creatorPtr int32) {
	getCreator(context, creatorPtr)
}

//export get_invoker
func get_invoker(context unsafe.Pointer, invokerPtr int32) {
	getInvoker(context, invokerPtr)
}

//export get_time
func get_time(context unsafe.Pointer, timePtr int32) {
	getTime(context, timePtr)
}


//export read_db
func read_db(context unsafe.Pointer, key, value int32) int32 {
	return readDB(context, key, value)
}

//export write_db
func write_db(context unsafe.Pointer, key, value int32) {
	writeDB(context, key, value)
}

//export delete_db
func delete_db(context unsafe.Pointer, key int32) {
	deleteDB(context, key)
}

type Wasmer struct {
	HomeDir      string              `json:"home_dir"`
	FilePathMap  map[string]string  `json:"file_path_map"`
	LastFileID   int				`json:"last_file_id"`
}


func NewWasmer(homeDir string, _ types.WasmConfig) (Wasmer, error){
	dir := filepath.Join(homeDir, types.FolderName)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return Wasmer{}, err
	}
	filePathMap := make(map[string]string)
	LastFileID := 0

	return Wasmer{
		HomeDir:     dir,
		FilePathMap: filePathMap,
		LastFileID:  LastFileID,
	}, nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func (w *Wasmer) Create(code []byte) (Wasmer,[]byte, error) {

	codeHash := MakeCodeHash(code)
	hash := fmt.Sprintf("%x", codeHash)
	if fName := w.FilePathMap[hash]; fName != "" {
		if FileExist(path.Join(w.HomeDir, fName)) {
			//file exist, remove file and delete
			err := os.Remove(path.Join(w.HomeDir, fName))
			if err != nil {
				return Wasmer{}, nil, err
			}
			delete(w.FilePathMap, hash)
		}else {
			//file not exist, delete
			delete(w.FilePathMap, hash)
		}
	}
	id, fileName := makeFilePath(w.LastFileID)
/*
	err := ioutil.WriteFile(w.HomeDir + "/" + fileName, code, types.ModePerm)
	if err != nil {
		return Wasmer{}, nil, err
	}
*/
	w.FilePathMap[hash] = fileName
	newWasmer := Wasmer{
		HomeDir:     w.HomeDir,
		FilePathMap: w.FilePathMap,
		LastFileID:  id,
	}
	return newWasmer,codeHash, nil
}

func (w *Wasmer) Instantiate(code []byte, funcName string, args json.RawMessage) (string, error) {
	instance , err := getInstance(code)
	if err != nil {
		return "", err
	}
	defer instance.Close()
	/*//直接引用go-ext-wasm的instance.
	instance, err := wasmer.NewInstance(code)
	//instance, err := wasmer.NewInstanceWithImports(code, imports)
	if err != nil {
		return "", err
	}*/
	init, exist := instance.Exports[funcName]
	if !exist {
		fmt.Println(exist)
		return "", errors.New("no expected function")
	}

	res, err := wasmCall(*instance, init, args)
	if err != nil {
		return "", err
	}
	if res.Err() {
		errStr := fmt.Sprintf("err: [%s]\n", res.ParseError())
		return "", errors.New(errStr)
	}else {
		res.Ok()
		resStr := fmt.Sprintf("ok: [%s]\n", string(res.Parse().Data))
		return resStr, nil
	}
}

func (w *Wasmer) Execute(code []byte, funcName string, args json.RawMessage) (string, error) {
	instance , err := getInstance(code)
	if err != nil {
		return "", err
	}
	defer instance.Close()
	/*instance, err := wasmer.NewInstance(code)
	//instance, err := wasmer.NewInstanceWithImports(code, imports)
	if err != nil {
		return "", err
	}*/
	handle, exist := instance.Exports[funcName]
	if !exist {
		return "", errors.New("no expected function")
	}
	res, err := wasmCall(*instance, handle, args)
	if err != nil {
		return "", errors.New("handle failed")
	}
	if res.Err() {
		errStr := fmt.Sprintf("err: [%s]\n", res.ParseError())
		return "", errors.New(errStr)
	}else {
		res.Ok()
		resStr := fmt.Sprintf("ok: [%s]\n", string(res.Parse().Data))
		return resStr, nil
	}
}

func (w *Wasmer) Query(code []byte, funcName string, args json.RawMessage) (string, error) {
	instance, err := getInstance(code)
	if err != nil {
		return "", err
	}
	defer instance.Close()

	query, exist := instance.Exports[funcName]
	if !exist {
		fmt.Println(exist)
		return "", errors.New("no expected function")
	}

	res, err := wasmCall(*instance, query, args)
	if err != nil {
		return "", errors.New("query failed")
	}
	if res.Err() {
		errStr := fmt.Sprintf("err: [%s]\n", res.ParseError())
		return "", errors.New(errStr)
	}else {
		res.Ok()
		resStr := fmt.Sprintf("ok: [%s]\n", string(res.Parse().Data))
		return resStr, nil
	}
}

func getInstance(code []byte) (*wasmer.Instance, error) {
	//store
	imports, err := wasmer.NewImports().Namespace("env").Append("read_db", read_db, C.read_db)
	if err != nil {
		return &wasmer.Instance{}, err
	}
	imports, err = imports.Namespace("env").Append("write_db", write_db, C.write_db)
	if err != nil {
		return &wasmer.Instance{}, err
	}
	imports, err = imports.Namespace("env").Append("delete_db", delete_db, C.delete_db)
	if err != nil {
		return &wasmer.Instance{}, err
	}

	//api
	imports, err = imports.Namespace("env").Append("send", send, C.send)
	if err != nil {
		return &wasmer.Instance{}, err
	}
	imports, err = imports.Namespace("env").Append("get_creator", get_creator, C.get_creator)
	if err != nil {
		return &wasmer.Instance{}, err
	}
	imports, err = imports.Namespace("env").Append("get_invoker", get_invoker, C.get_invoker)
	if err != nil {
		return &wasmer.Instance{}, err
	}
	imports, err = imports.Namespace("env").Append("get_time", get_time, C.get_time)
	if err != nil {
		return &wasmer.Instance{}, err
	}

	//gas
	imports, err = imports.Namespace("env").Append("addgas", addgas, C.addgas)
	if err != nil {
		return &wasmer.Instance{}, err
	}

	module, err := wasmer.Compile(code)
	if err != nil {
		panic(err)
	}
	defer module.Close()

	instance, err := module.InstantiateWithImports(imports)
	if err != nil {
		panic(err)
	}
	//defer instance.Close()
	allocate, exist := instance.Exports["allocate"]
	if !exist {
		fmt.Println(exist)
		return &wasmer.Instance{}, errors.New("no allocate")
	}
	middleIns.fun["allocate"] = allocate
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

func (w *Wasmer) GetWasmCode(hash []byte) ([]byte, error) {
	Hash := fmt.Sprintf("%x", hash)
	filePath := w.FilePathMap[Hash]
	code, err := ioutil.ReadFile(w.HomeDir + "/" + filePath)
	if err != nil {
		return nil, err
	}
	//the file may be not exist.
	return code, nil
}

func readCString(memory []byte) string {
	var res []byte
	for i := range memory {
		if memory[i] == 0 {
			break
		}
		res = append(res, memory[i])
	}
	return string(res)
}

func wasmCall(instance wasmer.Instance, fun func(...interface{}) (wasmer.Value, error), msg json.RawMessage) (ContractResult, error) {
	allocate, exist := middleIns.fun["allocate"]
	if !exist {
		panic("allocate not found")
	}

	var data []byte
	data = msg
	data = append(data, 0) // c str, + \0

	offset, err := allocate(len(data))
	if err != nil {
		return ContractResult{}, err
	}
	copy(instance.Memory.Data()[offset.ToI32():offset.ToI32()+int32(len(data))], data)

	res, err := fun(offset)
	if err != nil {
		return ContractResult{}, err
	}
	 str := readCString(instance.Memory.Data()[res.ToI32():])
	 resultMap := make(map[string]interface{})
	 if err := json.Unmarshal([]byte(str), &resultMap); err != nil {
	 	return ContractResult{}, err
	 }
	 return ContractResult{ resultMap, nil}, nil
}
