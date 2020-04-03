package keeper

// #include <stdlib.h>
//
// extern int read_db(void *context, int key, int value);
// extern void write_db(void *context, int key, int value);
// extern void delete_db(void *context, int key);
import "C"
import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/wasm/types"
	"github.com/wasmerio/go-ext-wasm/wasmer"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"unsafe"
)

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


func (w *Wasmer) Create(code []byte) (Wasmer,[]byte, error) {

	codeHash := MakeCodeHash(code)
	hash := fmt.Sprintf("%x", codeHash)
	if Str := w.FilePathMap[hash]; Str != "" {
		return Wasmer{}, nil, errors.New("the contract file has existed")
	}

	id, fileName := makeFilePath(w.LastFileID)
	err := ioutil.WriteFile(w.HomeDir + "/" + fileName, code, types.ModePerm)
	if err != nil {
		return Wasmer{}, nil, err
	}
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
	/*//直接引用go-ext-wasm的instance.
	instance, err := wasmer.NewInstance(code)
	//instance, err := wasmer.NewInstanceWithImports(code, imports)
	if err != nil {
		return "", err
	}*/
	init, exist := instance.Exports["init"]
	if !exist {
		fmt.Println(exist)
		return "", errors.New("no expected function")
	}

	res, err := init(args)
	if err != nil {
		return "", err
	}
	/*function:= instance.Exports[funcName]

	//TODO
	result, _ := function(1, 2)*/
	//result, _ := function(args)
	Result := res.String()
	return Result, nil

}

func (w *Wasmer) Execute(code []byte, funcName string, args json.RawMessage) (string, error) {
	instance , err := getInstance(code)
	if err != nil {
		return "", err
	}
	/*instance, err := wasmer.NewInstance(code)
	//instance, err := wasmer.NewInstanceWithImports(code, imports)
	if err != nil {
		return "", err
	}*/
	function:= instance.Exports[funcName]

	//result, _ := function(3, 4)
	result, _ := function(args)
	Result := result.String()
	return Result, nil
}

func (w *Wasmer) Query(code []byte, funcName string, args json.RawMessage) (string, error) {
	instance, err := getInstance(code)
	if err != nil {
		return "", err
	}
	/*instance, err := wasmer.NewInstance(code)
	//instance, err := wasmer.NewInstanceWithImports(code, imports)
	if err != nil {
		return "", err
	}*/
	function := instance.Exports[funcName]

	//TODO
	result, _ := function(3, 4)
	//result, _ := function(args)
	Result := result.String()
	return Result, nil
}

func getInstance(code []byte) (*wasmer.Instance, error) {
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

	module, err := wasmer.Compile(code)
	if err != nil {
		panic(err)
	}
	defer module.Close()

	instance, err := module.InstantiateWithImports(imports)
	if err != nil {
		panic(err)
	}
	defer instance.Close()
	allocate, exist := instance.Exports["allocate"]
	if !exist {
		fmt.Println(exist)
		return &wasmer.Instance{}, errors.New("no expected function")
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

func (w *Wasmer) GetWasmCode(id []byte) ([]byte, error) {
	hash := fmt.Sprintf("%x", id)
	filePath := w.FilePathMap[hash]
	code, err := ioutil.ReadFile(w.HomeDir + "/" + filePath)
	if err != nil {
		return nil, err
	}
	//the file may be not exist.
	return code, nil
}