package types

import (
	"fmt"
	"testing"
)

func TestPrefixKey(t *testing.T) {
	pk := NewPrefixedKey([]byte("store_name"), []byte("real_key"))
	rk := GetRealKey(pk)
	fmt.Println(string(pk))
	fmt.Println(string(rk))
}
