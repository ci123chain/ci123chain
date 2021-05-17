package cosmos_gravity

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNonce(t *testing.T) {
	var a uint64
	b := []byte("IjAi")
	err := json.Unmarshal(b, &a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}