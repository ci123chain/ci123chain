package subspace

type ParamSetPair struct {
	Key		[]byte
	Value 	interface{}
}

func NewParamSetPair(key []byte, value interface{}) ParamSetPair {
	return ParamSetPair{key, value}
}

type ParamSetPairs []ParamSetPair

type ParamSet interface {
	ParamSetPairs() ParamSetPairs
}