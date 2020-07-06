package keeper

// #include <stdlib.h>
//
// extern int read_db(void*, int, int, int, int, int);
// extern void write_db(void*, int, int, int, int);
// extern void delete_db(void*, int, int);
//
// extern int send(void*, int, long long);
// extern void get_creator(void*, int);
// extern void get_invoker(void*, int);
// extern long long get_time(void*);
//
// extern int get_input_length(void*, int);
// extern void get_input(void*, int, int, int);
// extern void notify_contract(void*, int, int);
// extern void return_contract(void*, int, int);
// extern int call_contract(void*, int, int, int);
// extern void destroy_contract(void*);
// extern int migrate_contract(void*, int, int, int, int, int, int, int, int, int, int, int, int, int);
// extern void panic_contract(void*, int, int);
//
// extern void addgas(void*, int);
import "C"
import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	"github.com/wasmerio/go-ext-wasm/wasmer"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"unicode/utf8"
	"unsafe"
)

var inputData = map[int32][]byte{}

const (
	InputDataTypeParam          = 0
	InputDataTypeContractResult = 1
)

type Param struct {
	Method string 	`json:"method"`
	Args   []string	`json:"args"`
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

//export send
func send(context unsafe.Pointer, to int32, amount int64) int32 {
	return performSend(context, to, amount)
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
func get_time(context unsafe.Pointer) int64 {
	return getTime(context)
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

//export migrate_contract
func migrate_contract(context unsafe.Pointer, codePtr, codeSize, namePtr, nameSize, verPtr, verSize,
	authorPtr, authorSize, emailPtr, emailSize, descPtr, descSize, newAddrPtr int32) int32 {
	return migrateContract(context, codePtr, codeSize, namePtr, nameSize, verPtr, verSize,
		authorPtr, authorSize, emailPtr, emailSize, descPtr, descSize, newAddrPtr)
}

//export destroy_contract
func destroy_contract(context unsafe.Pointer) {
	destroyContract(context)
}

//export panic_contract
func panic_contract(context unsafe.Pointer, dataPtr, dataSize int32) {
	panicContract(context, dataPtr, dataSize)
}

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
		panic(sdk.ErrorOutOfGas{Descriptor: "out of gas in location: vm"})
	}
	return
}

var creator sdk.AccAddress
func SetCreator(addr sdk.AccAddress) {
	creator = addr
}

var invoker sdk.AccAddress
func SetInvoker(addr sdk.AccAddress) {
	invoker = addr
}

var keeper *Keeper
func SetWasmKeeper(wk *Keeper) {
	keeper = wk
}

var invokeResult string
func ResetResult() {
	invokeResult = ""
}

var callResult []byte
func ResetCallResult() {
	callResult = []byte{}
}

var accountKeeper account.AccountKeeper
func SetAccountKeeper(ac account.AccountKeeper) {
	accountKeeper = ac
}

var ctx *sdk.Context
func SetCtx(con *sdk.Context) {
	ctx = con
}

type Wasmer struct {
	HomeDir      string             `json:"home_dir"`
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

func (w *Wasmer) Call(code []byte, args json.RawMessage) error {
	instance , err := getInstance(code)
	if err != nil {
		return err
	}
	defer instance.Close()

	invoke, exist := instance.Exports["invoke"]
	if !exist {
		fmt.Println(exist)
		return errors.New("no expected function")
	}
	var param Param
	inputByte, _ := args.MarshalJSON()
	fmt.Println(args)
	err = json.Unmarshal(inputByte, &param)
	if err != nil {
		return err
	}

	input := []interface{}{param.Method}
	for i := 0; i < len(param.Args); i++ {
		input = append(input, param.Args[i])
	}

	inputData[InputDataTypeParam] = serialize(input)

	fmt.Println(inputData)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	_, err = invoke()
	if err != nil {
		panic(err)
	}
	return nil
}

func (w *Wasmer) IndirectCall(code []byte, input []byte) error {
	instance , err := getInstance(code)
	if err != nil {
		return err
	}
	defer instance.Close()

	invoke, exist := instance.Exports["invoke"]
	if !exist {
		fmt.Println(exist)
		return errors.New("no expected function")
	}

	inputData[InputDataTypeParam] = input
	fmt.Println(inputData)
	_, err = invoke()
	if err != nil {
		return err
	}
	return nil
}

func getInstance(code []byte) (*wasmer.Instance, error) {
	imports, err := wasmer.NewImports().Namespace("env").Append("send", send, C.send)
	if err != nil {
		panic(err)
	}

	_, _ = imports.Append("read_db", read_db, C.read_db)
	_, _ = imports.Append("write_db", write_db, C.write_db)
	_, _ = imports.Append("delete_db", delete_db, C.delete_db)

	_, _ = imports.Append("get_creator", get_creator, C.get_creator)
	_, _ = imports.Append("get_invoker", get_invoker, C.get_invoker)
	_, _ = imports.Append("get_time", get_time, C.get_time)

	_, _ = imports.Append("get_input_length", get_input_length, C.get_input_length)
	_, _ = imports.Append("get_input", get_input, C.get_input)
	_, _ = imports.Append("return_contract", return_contract, C.return_contract)
	_, _ = imports.Append("notify_contract", notify_contract, C.notify_contract)
	_, _ = imports.Append("call_contract", call_contract, C.call_contract)
	_, _ = imports.Append("destroy_contract", destroy_contract, C.destroy_contract)
	_, _ = imports.Append("migrate_contract", migrate_contract, C.migrate_contract)
	_, _ = imports.Append("panic_contract", panic_contract, C.panic_contract)

	_, _ = imports.Append("addgas", addgas, C.addgas)
	module, err := wasmer.Compile(code)
	if err != nil {
		panic(err)
	}
	defer module.Close()

	instance, err := module.InstantiateWithImports(imports)
	if err != nil {
		panic(err)
	}
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

func (w *Wasmer) DeleteCode(hash []byte) error {
	Hash := fmt.Sprintf("%x", hash)
	filePath := w.FilePathMap[Hash]
	_, err := os.Lstat(w.HomeDir + "/" + filePath)
	if err == nil {
		err := os.Remove(w.HomeDir + "/" + filePath)
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

func serialize(raw []interface{}) (res []byte) {
	sink := NewSink(res)

	for i := range raw {
		switch r := raw[i].(type) {
		case string:
			//字符串必须是合法的utf8字符串
			if !utf8.ValidString(r) {
				panic("invalid utf8 string")
			}
			sink.WriteString(r)

		case uint64:
			sink.WriteU64(r)

		case uint32:
			sink.WriteU32(r)

		case []byte:
			sink.WriteBytes(r)

		case Address:
			sink.WriteAddress(r)

		default:
			panic("unexpected type")
		}
	}

	return sink.Bytes()
}

