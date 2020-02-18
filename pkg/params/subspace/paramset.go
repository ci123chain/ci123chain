package subspace

/*
type ParamSetPair struct {
	Key		[]byte
	Value 	interface{}
}

func NewParamSetPair(key []byte, value interface{}) ParamSetPair {
	return ParamSetPair{key, value}
}
*/
type (
	ValueValidatorFn func(value interface{}) error

	// ParamSetPair is used for associating paramsubspace key and field of param
	// structs.
	ParamSetPair struct {
		Key         []byte
		Value       interface{}
		ValidatorFn ValueValidatorFn
	}
)

// NewParamSetPair creates a new ParamSetPair instance.
func NewParamSetPair(key []byte, value interface{}, vfn ValueValidatorFn) ParamSetPair {
	return ParamSetPair{key, value, vfn}
}
type ParamSetPairs []ParamSetPair

type ParamSet interface {
	ParamSetPairs() ParamSetPairs
}