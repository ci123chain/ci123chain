package collactor

import (
	"fmt"
	"github.com/gogo/protobuf/proto"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

// RelayMsgs contains the msgs that need to be sent to both a src and dst chain
// after a given relay round. MaxTxSize and MaxMsgLength are ignored if they are
// set to zero.
type RelayMsgs struct {
	Src          []sdk.Msg `json:"src"`
	Dst          []sdk.Msg `json:"dst"`
	MaxTxSize    uint64    `json:"max_tx_size"`    // maximum permitted size of the msgs in a bundled relay transaction
	MaxMsgLength uint64    `json:"max_msg_length"` // maximum amount of messages in a bundled relay transaction

	Last      bool `json:"last"`
	Succeeded bool `json:"success"`
}


// Ready returns true if there are messages to relay
func (r *RelayMsgs) Ready() bool {
	if r == nil {
		return false
	}

	if len(r.Src) == 0 && len(r.Dst) == 0 {
		return false
	}
	return true
}


// Send sends the messages with appropriate output
// TODO: Parallelize? Maybe?
func (r *RelayMsgs) Send(src, dst *Chain) {
	r.SendWithController(src, dst, true)
}

// Success returns the success var
func (r *RelayMsgs) Success() bool {
	return r.Succeeded
}



func (r *RelayMsgs) SendWithController(src, dst *Chain, useController bool) {
	//if useController && SendToController != nil {
	//	action := &DeliverMsgsAction{
	//		Src:       MarshalChain(src),
	//		Dst:       MarshalChain(dst),
	//		Last:      r.Last,
	//		Succeeded: r.Succeeded,
	//		Type:      "RELAYER_SEND",
	//	}
	//
	//	action.SrcMsgs = EncodeMsgs(src, r.Src)
	//	action.DstMsgs = EncodeMsgs(dst, r.Dst)
	//
	//	// Get the messages that are actually sent.
	//	cont, err := ControllerUpcall(&action)
	//	if !cont {
	//		if err != nil {
	//			fmt.Println("Error calling controller", err)
	//			r.Succeeded = false
	//		} else {
	//			r.Succeeded = true
	//		}
	//		return
	//	}
	//}

	//nolint:prealloc // can not be pre allocated
	var (
		msgLen, txSize uint64
		msgs           []sdk.Msg
	)

	r.Succeeded = true

	// submit batches of relay transactions
	for _, msg := range r.Src {
		pbMsg := msg.(sdk.PbMsg)
		bz, err := proto.Marshal(pbMsg)
		if err != nil {
			panic(err)
		}

		msgLen++
		txSize += uint64(len(bz))

		if r.IsMaxTx(msgLen, txSize) {
			// Submit the transactions to src chain and update its status
			_, success, err := src.SendMsgs(msgs)
			if err != nil && src.debug {
				fmt.Println(err)
			}
			r.Succeeded = r.Succeeded && success

			// clear the current batch and reset variables
			msgLen, txSize = 1, uint64(len(bz))
			msgs = []sdk.Msg{}
		}
		msgs = append(msgs, msg)
	}

	// submit leftover msgs
	if len(msgs) > 0 {
		_, success, err := src.SendMsgs(msgs)
		if err != nil && src.debug {
			fmt.Println(err)
		}

		r.Succeeded = success
	}

	// reset variables
	msgLen, txSize = 0, 0
	msgs = []sdk.Msg{}

	for _, msg := range r.Dst {
		pbMsg := msg.(sdk.PbMsg)
		bz, err := proto.Marshal(pbMsg)
		if err != nil {
			panic(err)
		}

		msgLen++
		txSize += uint64(len(bz))

		if r.IsMaxTx(msgLen, txSize) {
			// Submit the transaction to dst chain and update its status
			_, success, err := dst.SendMsgs(msgs)
			if err != nil && dst.debug {
				fmt.Println(err)
			}

			r.Succeeded = r.Succeeded && success

			// clear the current batch and reset variables
			msgLen, txSize = 1, uint64(len(bz))
			msgs = []sdk.Msg{}
		}
		msgs = append(msgs, msg)
	}

	// submit leftover msgs
	if len(msgs) > 0 {
		_, success, err := dst.SendMsgs(msgs)
		if err != nil && dst.debug {
			fmt.Println(err)
		}

		r.Succeeded = success
	}
}

func (r *RelayMsgs) IsMaxTx(msgLen, txSize uint64) bool {
	return (r.MaxMsgLength != 0 && msgLen > r.MaxMsgLength) ||
		(r.MaxTxSize != 0 && txSize > r.MaxTxSize)
}