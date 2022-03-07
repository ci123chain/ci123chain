package cryptosuite

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestVerify(t *testing.T)  {
	ec := NewEccK1()

	msg := "this is a dog"

	sig := "0x5820209ae507fbbc8636114b4a209c87a5267802761a996f5f86bac8462835e903998d9163ee2e91f88c6a440a82f4f2c7c77e8459ebf851ef9ea1a1bf7204691b"
	sigbz := hexutil.MustDecode(sig)

	priv := "2cd9367f83f81341d17c6f126ea771e7c588c7fa902aca79284ec60ddc52de3a"
	privbz, _ := hex.DecodeString(priv)
	priKeyIns, err := crypto.ToECDSA(privbz)
	address := crypto.PubkeyToAddress(priKeyIns.PublicKey)
	res, err := ec.Verify(address.Bytes(), Hash([]byte(msg)), sigbz)
	if err != nil{
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestVerify2(t *testing.T)  {
	ec := NewEccK1()

	msg := "this is a dog"

	sig := "0x5820209ae507fbbc8636114b4a209c87a5267802761a996f5f86bac8462835e903998d9163ee2e91f88c6a440a82f4f2c7c77e8459ebf851ef9ea1a1bf72046900"
	sigbz := hexutil.MustDecode(sig)

	priv := "2cd9367f83f81341d17c6f126ea771e7c588c7fa902aca79284ec60ddc52de3a"
	privbz, _ := hex.DecodeString(priv)
	priKeyIns, err := crypto.ToECDSA(privbz)
	address := crypto.PubkeyToAddress(priKeyIns.PublicKey)
	res, err := ec.Verify(address.Bytes(), Hash([]byte(msg)), sigbz)
	if err != nil{
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestSign(t *testing.T)  {
	msg := "this is a dog"
	priv := "2cd9367f83f81341d17c6f126ea771e7c588c7fa902aca79284ec60ddc52de3a"
	ec := NewEccK1()
	digst := Hash([]byte(msg))
	fmt.Println("digst:", hex.EncodeToString(digst))
	privbz, _ := hex.DecodeString(priv)
	res, err := ec.Sign(privbz, digst)
	if err != nil{
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(res))
}

