package keeper

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
)

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
	//直接引用go-ext-wasm的instance.
	instance, err := wasmer.NewInstance(code)
	//instance, err := wasmer.NewInstanceWithImports(code, imports)
	if err != nil {
		return "", err
	}
	function:= instance.Exports[funcName]

	//TODO
	result, _ := function(1, 2)
	//result, _ := function(args)
	Result := result.String()
	return Result, nil

}

func (w *Wasmer) Execute(code []byte, funcName string, args json.RawMessage) (string, error) {
	instance, err := wasmer.NewInstance(code)
	//instance, err := wasmer.NewInstanceWithImports(code, imports)
	if err != nil {
		return "", err
	}
	function:= instance.Exports[funcName]

	//TODO
	result, _ := function(3, 4)
	//result, _ := function(args)
	Result := result.String()
	return Result, nil
}


func (w *Wasmer) Query(code []byte, funcName string, args json.RawMessage) (string, error) {
	instance, err := wasmer.NewInstance(code)
	//instance, err := wasmer.NewInstanceWithImports(code, imports)
	if err != nil {
		return "", err
	}
	function := instance.Exports[funcName]

	//TODO
	result, _ := function(3, 4)
	//result, _ := function(args)
	Result := result.String()
	return Result, nil
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