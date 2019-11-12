package types

import (
	"encoding/json"
	"github.com/pkg/errors"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"time"
)

const (
	StateKey = ".state="
	UniqueKey = ".uniqueid="
	StateReady = "ready"
	StateProcessing = "processing"
	StateDone = "done"

	TimeoutProcessing = 10
)

func ValidateState(state string) error {
	if state == StateReady ||
		state == StateProcessing ||
		state == StateDone {
		return nil
	}
	return errors.New("unknown state type")
}

type IBCMsg struct {
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



func (aa *IBCMsg) MarshalJSON() ([]byte, error) {
	type Alias IBCMsg
	return json.Marshal(&struct {
		UniqueID 	string		`json:"unique_id"`
		ObserverID 	string		`json:"observer_id"`
		*Alias
	}{
		UniqueID: string(aa.UniqueID),
		ObserverID: string(aa.ObserverID),
		Alias: (*Alias)(aa),
	})
}

func (aa *IBCMsg) UnmarshalJSON(data []byte) error {
	type Alias IBCMsg
	aux := &struct {
		UniqueID 	string		`json:"unique_id"`
		ObserverID 	string		`json:"observer_id"`
		*Alias
	}{
		Alias: (*Alias)(aa),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	aa.UniqueID = []byte(aux.UniqueID)
	aa.ObserverID = []byte(aux.ObserverID)
	return nil
}


func (msg IBCMsg) CanProcess() bool {
	if msg.State == StateReady {
		return true
	}
	if msg.State == StateProcessing {
		if time.Now().Unix() - msg.ApplyTime.Unix() > TimeoutProcessing {
			return true
		}
	}
	return false
}

type SignedIBCMsg struct {
	Signature 	[]byte 	`json:"signature"`
	IBCMsgBytes []byte	`json:"ibc_msg_bytes"`
}

func (sim SignedIBCMsg) Sign(priv []byte) (SignedIBCMsg, error) {
	signBytes := sim.GetSignBytes()

	sid := cryptosuit.NewFabSignIdentity()
	signature, err := sid.Sign(signBytes, priv)
	if err != nil {
		return sim, err
	}
	if len(signature) < 1 {
		return sim, errors.New("signature error")
	}
	sim.Signature = signature
	return sim, nil
}

func (msg *SignedIBCMsg)Bytes() []byte {
	bytes, err := IbcCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}


func (sim *SignedIBCMsg) GetSignBytes() []byte {
	tsim := *sim
	tsim.Signature = nil
	signBytes, err := IbcCdc.MarshalJSON(tsim)
	if err != nil {
		panic(err)
	}
	return signBytes
}