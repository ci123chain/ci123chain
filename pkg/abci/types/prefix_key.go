package types

import "strings"

func NewPrefixedKey(prefix, realKey []byte) (prefixedKey []byte){
	prefixKey := []byte("s/k:" + string(prefix) + "/")
	prefixedKey = append(prefixKey, realKey...)
	return
}

func GetRealKey(prefixedKey []byte) []byte {
	//iterator.Key()
	key := string(prefixedKey)
	keys := strings.SplitAfterN(key, "/",3)
	return []byte(keys[2])
}