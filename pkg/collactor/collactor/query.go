package collactor

import (
	"context"
	"fmt"
	accountutils "github.com/ci123chain/ci123chain/pkg/account/utils"
	chantypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	chanutils "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/utils"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	clientutils "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/utils"
	committypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	conntypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	connutils "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/utils"
	ibcexported "github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	stakingutils "github.com/ci123chain/ci123chain/pkg/staking/client/utils"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

// QueryLatestHeight queries the chain for the latest height and returns it
func (c *Chain) QueryLatestHeight() (int64, error) {
	res, err := c.Client.Status(context.Background())
	if err != nil {
		return -1, err
	} else if res.SyncInfo.CatchingUp {
		return -1, fmt.Errorf("node at %s running chain %s not caught up", c.RPCAddr, c.ChainID)
	}

	return res.SyncInfo.LatestBlockHeight, nil
}

// QueryClientState retrevies the latest consensus state for a client in state at a given height
func (c *Chain) QueryClientState(height int64) (*clienttypes.QueryClientStateResponse, error) {
	return clientutils.QueryClientStateABCI(c.CLIContext(height), c.PathEnd.ClientID)
}

// QueryClients queries all the clients!
func (c *Chain) QueryClients(offset, limit uint64) (*clienttypes.QueryClientStatesResponse, error) {
	return clientutils.QueryClientStatesABCI(c.CLIContext(0), offset, limit)
}


// QueryConnection returns the remote end of a given connection
func (c *Chain) QueryConnection(height int64) (*conntypes.QueryConnectionResponse, error) {
	res, err := connutils.QueryConnection(c.CLIContext(height), c.PathEnd.ConnectionID, true)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return emptyConnRes, nil
	} else if err != nil {
		return nil, err
	}
	return res, nil
}


// ////////////////////////////
//  ICS 03 -> CONNECTIONS   //
// ////////////////////////////

// QueryConnections gets any connections on a chain
func (c *Chain) QueryConnections(
	offset, limit uint64) (*conntypes.QueryConnectionsResponse, error) {
	return connutils.QueryConnectionsABCI(c.CLIContext(0), offset, limit)
}


// QueryClientConsensusState retrevies the latest consensus state for a client in state at a given height
func (c *Chain) QueryClientConsensusState(
	height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	return clientutils.QueryConsensusStateABCI(
		c.CLIContext(height),
		c.PathEnd.ClientID,
		dstClientConsHeight,
	)
}


// QueryChannelPair returns a pair of channel responses
func QueryChannelPair(src, dst *Chain, srcH, dstH int64) (srcChan, dstChan *chantypes.QueryChannelResponse, err error) {
	var eg = new(errgroup.Group)
	eg.Go(func() error {
		srcChan, err = src.QueryChannel(srcH)
		return err
	})
	eg.Go(func() error {
		dstChan, err = dst.QueryChannel(dstH)
		return err
	})
	err = eg.Wait()
	return
}

// QueryChannel returns the channel associated with a channelID
func (c *Chain) QueryChannel(height int64) (chanRes *chantypes.QueryChannelResponse, err error) {
	res, err := chanutils.QueryChannel(c.CLIContext(height), c.PathEnd.PortID, c.PathEnd.ChannelID, true)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return emptyChannelRes, nil
	} else if err != nil {
		return nil, err
	}
	return res, nil
}


// QueryChannels returns all the channels that are registered on a chain
func (c *Chain) QueryChannels(offset, limit uint64) (*chantypes.QueryChannelsResponse, error) {
	res, err := chanutils.QueryChannelsABCI(c.CLIContext(0), offset, limit)
	return res, err
}



// QueryUnbondingPeriod returns the unbonding period of the chain
func (c *Chain) QueryUnbondingPeriod() (time.Duration, error) {
	params, err := stakingutils.QueryParms(c.CLIContext(0))
	if err != nil {
		return 0, err
	}
	return params.UnbondingTime, nil
}

func (c *Chain) QueryNonce() (uint64, error) {
	nonce, err := accountutils.QueryNonce(c.CLIContext(0), c.MustGetAddress())
	return nonce, err
}



var emptyConnRes = conntypes.NewQueryConnectionResponse(
	conntypes.NewConnectionEnd(
		conntypes.UNINITIALIZED,
		"client",
		conntypes.NewCounterparty(
			"client",
			"connection",
			committypes.NewMerklePrefix([]byte{}),
		),
		[]*conntypes.Version{},
		0,
	),
	[]byte{},
	clienttypes.NewHeight(0, 0),
)

var emptyChannelRes = chantypes.NewQueryChannelResponse(
	chantypes.NewChannel(
		chantypes.UNINITIALIZED,
		chantypes.UNORDERED,
		chantypes.NewCounterparty(
			"port",
			"channel",
		),
		[]string{},
		"version",
	),
	[]byte{},
	clienttypes.NewHeight(0, 0),
)

// QueryConnectionPair returns a pair of connection responses
func QueryConnectionPair(
	src, dst *Chain,
	srcH, dstH int64) (srcConn, dstConn *conntypes.QueryConnectionResponse, err error) {
	var eg = new(errgroup.Group)
	eg.Go(func() error {
		srcConn, err = src.QueryConnection(srcH)
		return err
	})
	eg.Go(func() error {
		dstConn, err = dst.QueryConnection(dstH)
		return err
	})
	err = eg.Wait()
	return
}
