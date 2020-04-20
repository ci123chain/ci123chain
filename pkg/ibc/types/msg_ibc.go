package types

import (
	"encoding/json"
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"time"
)

const (
	StateKey = ".state="
	UniqueKey = ".uniqueid="
	TimestampKey = ".timestamp="
	StateReady = "ready"
	StateProcessing = "processing"
	StateDone = "done"
	StateCancel = "cancel"

	TimeoutProcessing = 30
)

func ValidateState(state string) sdk.Error {
	if state == StateReady ||
		state == StateProcessing ||
		state == StateDone {
		return nil
	}
	return ErrState(DefaultCodespace, errors.New("state not ready processing or done"))
}

// 链上保存的跨链交易信息
type IBCInfo struct {
	// 银行地址
	BankAddress sdk.AccAddress	`json:"bank_address"`
	// 跨链交易ID
	UniqueID 	[]byte			`json:"unique_id"`
	ObserverID 	[]byte			`json:"observer_id"`
	ApplyTime 	time.Time		`json:"apply_time"`
	State 	 	string 			`json:"state"`
	FromAddress sdk.AccAddress 	`json:"from_address"`
	ToAddress 	sdk.AccAddress 	`json:"to_address"`
	Amount 	    sdk.Coin   		`json:"amount"`
}


//
//func (aa *IBCInfo) MarshalJSON() ([]byte, error) {
//	type Alias IBCInfo
//	return json.Marshal(&struct {
//		UniqueID 	string		`json:"unique_id"`
//		ObserverID 	string		`json:"observer_id"`
//		*Alias
//	}{
//		UniqueID: string(aa.UniqueID),
//		ObserverID: string(aa.ObserverID),
//		Alias: (*Alias)(aa),
//	})
//}
//
//func (aa *IBCInfo) UnmarshalJSON(data []byte) error {
//	type Alias IBCInfo
//	aux := &struct {
//		UniqueID 	string		`json:"unique_id"`
//		ObserverID 	string		`json:"observer_id"`
//		*Alias
//	}{
//		Alias: (*Alias)(aa),
//	}
//	if err := json.Unmarshal(data, &aux); err != nil {
//		return err
//	}
//	fmt.Println(aux.UniqueID)
//	aa.UniqueID = []byte(aux.UniqueID)
//	aa.ObserverID = []byte(aux.ObserverID)
//	return nil
//}


func (msg IBCInfo) CanProcess(now time.Time) bool {
	if msg.State == StateReady {
		return true
	}
	if msg.State == StateProcessing {
		if now.Unix() - msg.ApplyTime.Unix() > TimeoutProcessing {
			return true
		}
	}
	return false
}

// 调用 Apply 返回的签名信息
type ApplyReceipt struct {
	Signature 	[]byte 	`json:"signature"`
	IBCMsgBytes []byte	`json:"ibc_msg_bytes"`
}

func (sim ApplyReceipt) Sign(priv []byte) (ApplyReceipt, error) {
	signBytes := sim.GetSignBytes()

	sid := cryptosuit.NewFabSignIdentity()
	signature, err := sid.Sign(signBytes, priv)
	if err != nil {
		return sim, transaction.ErrSignature(DefaultCodespace, err)
	}
	if len(signature) < 1 {
		return sim, transaction.ErrSignature(DefaultCodespace, errors.New("len signature less than 1"))
	}
	sim.Signature = signature
	return sim, nil
}

func (msg *ApplyReceipt)Bytes() []byte {
	bytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}


func (sim *ApplyReceipt) GetSignBytes() []byte {
	tsim := *sim
	tsim.Signature = nil
	signBytes := tsim.Bytes()
	return signBytes
}


// ci 给 fabric 的转账回执
type BankReceipt struct {
	UniqueID 	string	`json:"unique_id"`
	ObserverID 	string	`json:"observer_id"`
	Signature 	[]byte	`json:"signature"`
}

func (br *BankReceipt) GetSignBytes() []byte {
	tsim := *br
	tsim.Signature = nil
	signBytes := tsim.Bytes()
	return signBytes
}

func (br *BankReceipt) Bytes() []byte {
	bytes, err := json.Marshal(*br)
	if err != nil {
		panic(err)
	}
	return bytes
}

func NewBankReceipt(uniqueID, observerID string) *BankReceipt {
	return &BankReceipt{
		UniqueID: uniqueID,
		ObserverID: observerID,
	}
}

func (br *BankReceipt) Sign(priv []byte) (*BankReceipt, error) {
	signBytes := br.GetSignBytes()
	sid := cryptosuit.NewFabSignIdentity()
	signature, err := sid.Sign(signBytes, priv)
	if err != nil {
		return br, transaction.ErrSignature(DefaultCodespace, err)
	}
	if len(signature) < 1 {
		return br, transaction.ErrSignature(DefaultCodespace, errors.New("len signature less than 1"))
	}
	br.Signature = signature
	return br, nil
}