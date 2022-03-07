package cryptosuite

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestSaveKey(T *testing.T) {

	priKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
	}
	priByte := crypto.FromECDSA(priKey)
	priHex := hex.EncodeToString(priByte)
	//b := make([]byte, 16)
	//for i := 0; i < 10; i++ {
	//	_, err = rand.Read(b)
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	break
	//}
	b := []byte("123456")
	syKey := "1A2b34p567f8R90w"
	savedKey, err := EncryptKey(priByte, []byte(syKey), b)
	if err != nil {
		fmt.Println(err)
	}
	decryptKeyByte, err := DecryptKey(savedKey, []byte(syKey), b)
	if err != nil {
		fmt.Println(err)
	}
	decryptPriKey, err := crypto.ToECDSA(decryptKeyByte)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(decryptPriKey)
	decryptHex := hex.EncodeToString(decryptKeyByte)
	fmt.Println(priHex == decryptHex)

	//priKey , err := ParsePemPriKey(eccR1PriKey)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//b,err := json.Marshal(priKey)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(b))
	//
	//s := &ecdsa.PrivateKey{}
	//err  = json.Unmarshal(b, &s)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(s)

}
