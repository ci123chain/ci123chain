package collactor

import (
	"fmt"
	"github.com/avast/retry-go"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	transfertypes "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	chantypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	conntypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	tmclient "github.com/ci123chain/ci123chain/pkg/ibc/light-clients/07-tendermint/types"
	"github.com/pkg/errors"
)

// CreateClient creates an sdk.Msg to update the client on src with consensus state from dst
func (c *Chain) CreateClient(
//nolint:interfacer
	clientState *tmclient.ClientState,
	dstHeader *tmclient.Header) sdk.Msg {

	if err := dstHeader.ValidateBasic(); err != nil {
		panic(err)
	}

	msg, err := clienttypes.NewMsgCreateClient(
		clientState,
		dstHeader.ConsensusState(),
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)
	if err != nil {
		panic(err)
	}

	if err := msg.ValidateBasic(); err != nil {
		panic(err)
	}
	return msg
}


// ConnInit creates a MsgConnectionOpenInit
func (c *Chain) ConnInit(counterparty *Chain) ([]sdk.Msg, error) {
	updateMsg, err := c.UpdateClient(counterparty)
	if err != nil {
		return nil, err
	}

	var version *conntypes.Version
	msg := conntypes.NewMsgConnectionOpenInit(
		c.PathEnd.ClientID,
		counterparty.PathEnd.ClientID,
		defaultChainPrefix,
		version,
		defaultDelayPeriod,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

	return []sdk.Msg{updateMsg, msg}, nil

}

// ConnTry creates a MsgConnectionOpenTry
func (c *Chain) ConnTry(
	counterparty *Chain,
) ([]sdk.Msg, error) {
	updateMsg, err := c.UpdateClient(counterparty)
	if err != nil {
		return nil, err
	}

	// NOTE: the proof height uses - 1 due to tendermint's delayed execution model
	useHeight := counterparty.MustGetLatestLightHeight() - 1
	clientState, clientStateProof, consensusStateProof, connStateProof,
	proofHeight, err := counterparty.GenerateConnHandshakeProof(useHeight)
	if err != nil {
		return nil, err
	}
	// TODO: Get DelayPeriod from counterparty connection rather than using default value
	msg := conntypes.NewMsgConnectionOpenTry(
		c.PathEnd.ConnectionID,
		c.PathEnd.ClientID,
		counterparty.PathEnd.ConnectionID,
		counterparty.PathEnd.ClientID,
		clientState,
		defaultChainPrefix,
		conntypes.ExportedVersionsToProto(conntypes.GetCompatibleVersions()),
		defaultDelayPeriod,
		connStateProof,
		clientStateProof,
		consensusStateProof,
		proofHeight,
		clientState.GetLatestHeight().(clienttypes.Height),
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return []sdk.Msg{updateMsg, msg}, nil
}


// ConnAck creates a MsgConnectionOpenAck
func (c *Chain) ConnAck(
	counterparty *Chain,
) ([]sdk.Msg, error) {
	updateMsg, err := c.UpdateClient(counterparty)
	if err != nil {
		return nil, err
	}

	// NOTE: the proof height uses - 1 due to tendermint's delayed execution model
	clientState, clientStateProof, consensusStateProof, connStateProof,
	proofHeight, err := counterparty.GenerateConnHandshakeProof(counterparty.MustGetLatestLightHeight() - 1)
	if err != nil {
		return nil, err
	}

	msg := conntypes.NewMsgConnectionOpenAck(
		c.PathEnd.ConnectionID,
		counterparty.PathEnd.ConnectionID,
		clientState,
		connStateProof,
		clientStateProof,
		consensusStateProof,
		proofHeight,
		clientState.GetLatestHeight().(clienttypes.Height),
		conntypes.DefaultIBCVersion,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return []sdk.Msg{updateMsg, msg}, nil
}


// ConnConfirm creates a MsgConnectionOpenConfirm
func (c *Chain) ConnConfirm(counterparty *Chain) ([]sdk.Msg, error) {
	updateMsg, err := c.UpdateClient(counterparty)
	if err != nil {
		return nil, err
	}

	counterpartyConnState, err := counterparty.QueryConnection(int64(counterparty.MustGetLatestLightHeight()) - 1)
	if err != nil {
		return nil, err
	}
	if counterpartyConnState.Connection.State != conntypes.OPEN {
		return nil, errors.New(fmt.Sprintf("counterparty connection state error, expected Open(3), Got %d", counterpartyConnState.Connection.State))
	}
	msg := conntypes.NewMsgConnectionOpenConfirm(
		c.PathEnd.ConnectionID,
		counterpartyConnState.Proof,
		counterpartyConnState.ProofHeight,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return []sdk.Msg{updateMsg, msg}, nil
}


// ChanInit creates a MsgChannelOpenInit
func (c *Chain) ChanInit(counterparty *Chain) ([]sdk.Msg, error) {
	updateMsg, err := c.UpdateClient(counterparty)
	if err != nil {
		return nil, err
	}

	msg := chantypes.NewMsgChannelOpenInit(
		c.PathEnd.PortID,
		c.PathEnd.Version,
		c.PathEnd.GetOrder(),
		[]string{c.PathEnd.ConnectionID},
		counterparty.PathEnd.PortID,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return []sdk.Msg{updateMsg, msg}, nil
}

// ChanTry creates a MsgChannelOpenTry
func (c *Chain) ChanTry(
	counterparty *Chain,
) ([]sdk.Msg, error) {
	updateMsg, err := c.UpdateClient(counterparty)
	if err != nil {
		return nil, err
	}

	// NOTE: the proof height uses - 1 due to tendermint's delayed execution model
	useHeight := counterparty.MustGetLatestLightHeight() - 1
	counterpartyChannelRes, err := counterparty.QueryChannel(int64(useHeight))
	if err != nil {
		return nil, err
	}

	msg := chantypes.NewMsgChannelOpenTry(
		c.PathEnd.PortID,
		c.PathEnd.ChannelID,
		c.PathEnd.Version,
		counterpartyChannelRes.Channel.Ordering,
		[]string{c.PathEnd.ConnectionID},
		counterparty.PathEnd.PortID,
		counterparty.PathEnd.ChannelID,
		counterpartyChannelRes.Channel.Version,
		counterpartyChannelRes.Proof,
		counterpartyChannelRes.ProofHeight,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return []sdk.Msg{updateMsg, msg}, nil
}


// ChanAck creates a MsgChannelOpenAck
func (c *Chain) ChanAck(
	counterparty *Chain,
) ([]sdk.Msg, error) {
	updateMsg, err := c.UpdateClient(counterparty)
	if err != nil {
		return nil, err
	}

	// NOTE: the proof height uses - 1 due to tendermint's delayed execution model
	counterpartyChannelRes, err := counterparty.QueryChannel(int64(counterparty.MustGetLatestLightHeight()) - 1)
	if err != nil {
		return nil, err
	}

	msg := chantypes.NewMsgChannelOpenAck(
		c.PathEnd.PortID,
		c.PathEnd.ChannelID,
		counterparty.PathEnd.ChannelID,
		counterpartyChannelRes.Channel.Version,
		counterpartyChannelRes.Proof,
		counterpartyChannelRes.ProofHeight,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return []sdk.Msg{updateMsg, msg}, nil
}

// ChanConfirm creates a MsgChannelOpenConfirm
func (c *Chain) ChanConfirm(counterparty *Chain) ([]sdk.Msg, error) {
	updateMsg, err := c.UpdateClient(counterparty)
	if err != nil {
		return nil, err
	}

	// NOTE: the proof height uses - 1 due to tendermint's delayed execution model
	counterpartyChanState, err := counterparty.QueryChannel(int64(counterparty.MustGetLatestLightHeight()) - 1)
	if err != nil {
		return nil, err
	}
	if counterpartyChanState.Channel.State != chantypes.OPEN {
		return nil, errors.New(fmt.Sprintf("counterparty channel state error, expected Open(3), Got %d", counterpartyChanState.Channel.State))
	}

	msg := chantypes.NewMsgChannelOpenConfirm(
		c.PathEnd.PortID,
		c.PathEnd.ChannelID,
		counterpartyChanState.Proof,
		counterpartyChanState.ProofHeight,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return []sdk.Msg{updateMsg, msg}, nil
}


// UpdateClient creates an sdk.Msg to update the client on src with data pulled from dst
// at the request height..
func (c *Chain) UpdateClient(dst *Chain) (sdk.Msg, error) {
	header, err := dst.GetIBCUpdateHeader(c)
	if err != nil {
		return nil, err
	}

	if err := header.ValidateBasic(); err != nil {
		return nil, err
	}
	msg, err := clienttypes.NewMsgUpdateClient(
		c.PathEnd.ClientID,
		header,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)
	if err != nil {
		return nil, err
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

func (c *Chain) UpdateClient2(dst *Chain) (sdk.Msg, int64, error) {
	header, err := dst.GetIBCUpdateHeader(c)
	if err != nil {
		return nil, -1, err
	}

	if err := header.ValidateBasic(); err != nil {
		return nil, -1,err
	}
	msg, err := clienttypes.NewMsgUpdateClient(
		c.PathEnd.ClientID,
		header,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)
	if err != nil {
		return nil, -1, err
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, -1, err
	}
	return msg, header.Header.Height,nil
}



// MsgRelayRecvPacket constructs the MsgRecvPacket which is to be sent to the receiving chain.
// The counterparty represents the sending chain where the packet commitment would be stored.
func (c *Chain) MsgRelayRecvPacket(counterparty *Chain, packet *relayMsgRecvPacket) (msgs []sdk.Msg, err error) {
	var comRes *chantypes.QueryPacketCommitmentResponse
	if err = retry.Do(func() (err error) {
		// NOTE: the proof height uses - 1 due to tendermint's delayed execution model
		comRes, err = counterparty.QueryPacketCommitment(int64(counterparty.MustGetLatestLightHeight()) - 1, packet.seq)
		if err != nil {
			return err
		}

		if comRes.Proof == nil || comRes.Commitment == nil {
			return fmt.Errorf("recv packet commitment query seq(%d) is nil", packet.seq)
		}

		return nil
	}, rtyAtt, rtyDel, rtyErr, retry.OnRetry(func(n uint, _ error) {
		// clear messages
		msgs = []sdk.Msg{}

		// OnRetry we want to update the light clients and then debug log
		updateMsg, err := c.UpdateClient(counterparty)
		if err != nil {
			return
		}

		msgs = append(msgs, updateMsg)

		if counterparty.debug {
			counterparty.Log(fmt.Sprintf("- [%s]@{%d} - try(%d/%d) query packet commitment: %s",
				counterparty.ChainID, counterparty.MustGetLatestLightHeight()-1, n+1, rtyAttNum, err))
		}

	})); err != nil {
		counterparty.Error(err)
		return
	}

	if comRes == nil {
		return nil, fmt.Errorf("receive packet [%s]seq{%d} has no associated proofs", c.ChainID, packet.seq)
	}

	msg := chantypes.NewMsgRecvPacket(
		chantypes.NewPacket(
			packet.packetData,
			packet.seq,
			counterparty.PathEnd.PortID,
			counterparty.PathEnd.ChannelID,
			c.PathEnd.PortID,
			c.PathEnd.ChannelID,
			packet.timeout,
			packet.timeoutStamp,
		),
		comRes.Proof,
		comRes.ProofHeight,
		c.MustGetAddressString(),
	)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return append(msgs, msg), nil
}


// MsgRelayTimeout constructs the MsgTimeout which is to be sent to the sending chain.
// The counterparty represents the receiving chain where the receipts would have been
// stored.
func (c *Chain) MsgRelayTimeout(counterparty *Chain, packet *relayMsgTimeout) (msgs []sdk.Msg, err error) {
	var recvRes *chantypes.QueryPacketReceiptResponse
	if err = retry.Do(func() (err error) {
		// NOTE: Timeouts currently only work with ORDERED channels for nwo
		// NOTE: the proof height uses - 1 due to tendermint's delayed execution model
		useDstHeight := counterparty.MustGetLatestLightHeight() - 1
		recvRes, err = counterparty.QueryPacketReceipt(int64(useDstHeight), packet.seq)

		if err != nil {
			return err
		}

		if recvRes.Proof == nil {
			return fmt.Errorf("timeout packet receipt proof seq(%d) is nil", packet.seq)
		}

		return nil
	}, rtyAtt, rtyDel, rtyErr, retry.OnRetry(func(n uint, _ error) {
		// clear messages
		msgs = []sdk.Msg{}

		// OnRetry we want to update the light clients and then debug log
		updateMsg, err := c.UpdateClient(counterparty)
		if err != nil {
			return
		}

		msgs = append(msgs, updateMsg)

		if counterparty.debug {
			counterparty.Log(fmt.Sprintf("- [%s]@{%d} - try(%d/%d) query packet receipt: %s",
				counterparty.ChainID, counterparty.MustGetLatestLightHeight()-1, n+1, rtyAttNum, err))
		}

	})); err != nil {
		counterparty.Error(err)
		return
	}

	if recvRes == nil {
		return nil, fmt.Errorf("timeout packet [%s]seq{%d} has no associated proofs", c.ChainID, packet.seq)
	}

	msg := chantypes.NewMsgTimeout(
		chantypes.NewPacket(
			packet.packetData,
			packet.seq,
			c.PathEnd.PortID,
			c.PathEnd.ChannelID,
			counterparty.PathEnd.PortID,
			counterparty.PathEnd.ChannelID,
			packet.timeout,
			packet.timeoutStamp,
		),
		packet.seq,
		recvRes.Proof,
		recvRes.ProofHeight,
		c.MustGetAddressString(),
	)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return append(msgs, msg), nil
}



// MsgTransfer creates a new transfer message
func (c *Chain) MsgTransfer(dst *PathEnd, amount sdk.Coin, dstAddr string,
	timeoutHeight, timeoutTimestamp uint64) sdk.Msg {
	version := clienttypes.ParseChainID(dst.ChainID)
	return transfertypes.NewMsgTransfer(
		c.PathEnd.PortID,
		c.PathEnd.ChannelID,
		amount,
		c.MustGetAddressString(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
		dstAddr,
		clienttypes.NewHeight(version, timeoutHeight),
		timeoutTimestamp,
	)
}
