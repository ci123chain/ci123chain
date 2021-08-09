package keeper
import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
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


type VMRes struct {
	err []byte // error  response tip
	res []byte // success
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
	//defer instance.Close()
	methodEncode := "x" + hex.EncodeToString([]byte(method))
	invokeCall, err := instance.Exports.GetRawFunction(methodEncode)
	invoke := invokeCall.Native()
	if err != nil {
		return nil, evmtypes.ErrContractMethodInvalid.Wrap("get method name " + method)
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

	_, err2 := invoke()
	if err2 != nil {
		panic(err2)
	}
	return res, err
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