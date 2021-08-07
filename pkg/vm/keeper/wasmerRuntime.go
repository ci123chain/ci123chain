package keeper
import "C"
import (
	"fmt"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type WasmerContext struct {
	*wasmer.Instance
	cfg *runtimeConfig
}

func (wc *WasmerContext) getMemory() []byte {
	memoryIns, err := wc.Instance.Exports.GetMemory("memory")
	if err != nil {
		panic("get memory should not error")
	}
	memory := memoryIns.Data()
	return memory
}


func getInstance(wasmBytes []byte, cfg *runtimeConfig) (*WasmerContext, error) {
	wCtx := &WasmerContext{
		cfg: cfg,
	}

	// Create an Engine
	engine := wasmer.NewEngine()
	// Create a Store
	store := wasmer.NewStore(engine)
	module, err := wasmer.NewModule(store, wasmBytes)

	if err != nil {
		fmt.Println("Failed to compile module:", err)
	}
	importObject := wasmer.NewImportObject()

	newFunc := func(f func([]wasmer.Value) ([]wasmer.Value, error), in []*wasmer.ValueType, out []*wasmer.ValueType) *wasmer.Function {
		function := wasmer.NewFunction(
			store,
			wasmer.NewFunctionType(in, out),
			f,
		)
		return function
	}

	importObject.Register(
		"env",
		map[string]wasmer.IntoExtern{
			"send": newFunc(wCtx.performSend, wasmer.NewValueTypes(wasmer.I32, wasmer.I64), wasmer.NewValueTypes(wasmer.I32)),
			"read_db": newFunc(wCtx.readDB, wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32), wasmer.NewValueTypes(wasmer.I32)),
			"write_db": newFunc(wCtx.writeDB, wasmer.NewValueTypes(wasmer.I32,wasmer.I32,wasmer.I32,wasmer.I32), wasmer.NewValueTypes()),
			"delete_db": newFunc(wCtx.deleteDB, wasmer.NewValueTypes(wasmer.I32,wasmer.I32), wasmer.NewValueTypes()),
			"new_db_iter": newFunc(wCtx.newDBIter, wasmer.NewValueTypes(wasmer.I32,wasmer.I32), wasmer.NewValueTypes(wasmer.I32)),
			"db_iter_next": newFunc(wCtx.dbIterNext, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes(wasmer.I32)),
			"db_iter_key": newFunc(wCtx.dbIterKey, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes(wasmer.I32)),
			"db_iter_value": newFunc(wCtx.dbIterValue, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes(wasmer.I32)),

			"get_creator": newFunc(wCtx.getCreator, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
			"get_invoker": newFunc(wCtx.getInvoker, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
			"self_address": newFunc(wCtx.selfAddress, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
			"get_pre_caller": newFunc(wCtx.getPreCaller, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
			"get_block_header": newFunc(wCtx.getBlockHeader, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),

			"get_input_length": newFunc(wCtx.getInputLength, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes(wasmer.I32)),
			"get_input": newFunc(wCtx.getInput, wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32), wasmer.NewValueTypes()),
			"return_contract": newFunc(wCtx.returnContract, wasmer.NewValueTypes(wasmer.I32,wasmer.I32), wasmer.NewValueTypes()),
			"notify_contract": newFunc(wCtx.notifyContract, wasmer.NewValueTypes(wasmer.I32,wasmer.I32), wasmer.NewValueTypes()),
			"call_contract": newFunc(wCtx.callContract, wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32), wasmer.NewValueTypes(wasmer.I32)),
			"new_contract": newFunc(wCtx.newContract, wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32), wasmer.NewValueTypes()),
			"destroy_contract": newFunc(wCtx.destroyContract, wasmer.NewValueTypes(), wasmer.NewValueTypes()),
			"panic_contract": newFunc(wCtx.panicContract, wasmer.NewValueTypes(wasmer.I32,wasmer.I32), wasmer.NewValueTypes()),

			"get_validator_power": newFunc(wCtx.getValidatorPower, wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32), wasmer.NewValueTypes()),
			"total_power": newFunc(wCtx.totalPower, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
			"get_balance": newFunc(wCtx.getBalance, wasmer.NewValueTypes(wasmer.I32, wasmer.I32), wasmer.NewValueTypes()),
			"addgas": newFunc(wCtx.addGas, wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
			"debug_print": newFunc(wCtx.debugPrint, wasmer.NewValueTypes(wasmer.I32, wasmer.I32), wasmer.NewValueTypes()),
		},
	)


	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		return nil, err
	}
	wCtx.Instance = instance
	return wCtx, nil
}