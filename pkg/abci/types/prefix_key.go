package types

func NewPrefixedKey(prefix, key []byte) (realKey []byte){
	prefixKey := []byte("s/k:" + string(prefix) + "/")
	realKey = append(prefixKey, key...)
	return
}
