package validator

import (
	"fmt"
	"testing"
)

func TestNewValidatorKey(t *testing.T) {
	validatorKey, pubKey, err := NewValidatorKey()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(validatorKey)
	fmt.Println(pubKey)
}