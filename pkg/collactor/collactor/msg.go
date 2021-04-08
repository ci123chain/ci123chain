package collactor

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	chantypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	conntypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	tmclient "github.com/ci123chain/ci123chain/pkg/ibc/light-clients/07-tendermint/types"
)

// CreateClient creates an sdk.Msg to update the client on src with consensus state from dst
func (c *Chain) CreateClient(
//nolint:interfacer
	clientState *tmclient.ClientState,
	dstHeader *tmclient.Header) sdk.Msg {

	if err := dstHeader.ValidateBasic(); err != nil {
		panic(err)
	}

	msg := clienttypes.NewMsgCreateClient(
		clientState,
		dstHeader.ConsensusState(),
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

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
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
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
	clientState, clientStateProof, consensusStateProof, connStateProof,
	proofHeight, err := counterparty.GenerateConnHandshakeProof(counterparty.MustGetLatestLightHeight() - 1)
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
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
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
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

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

	msg := conntypes.NewMsgConnectionOpenConfirm(
		c.PathEnd.ConnectionID,
		counterpartyConnState.Proof,
		counterpartyConnState.ProofHeight,
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

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
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

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
	counterpartyChannelRes, err := counterparty.QueryChannel(int64(counterparty.MustGetLatestLightHeight()) - 1)
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
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)
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
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

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

	msg := chantypes.NewMsgChannelOpenConfirm(
		c.PathEnd.PortID,
		c.PathEnd.ChannelID,
		counterpartyChanState.Proof,
		counterpartyChanState.ProofHeight,
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)

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
	msg := clienttypes.NewMsgUpdateClient(
		c.PathEnd.ClientID,
		header,
		c.MustGetAddress(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
	)
	return msg, nil
}
