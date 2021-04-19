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
	transferutils "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/utils"
	transfertypes "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	"github.com/pkg/errors"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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

// QueryPacketCommitment returns the packet commitment proof at a given height
func (c *Chain) QueryPacketCommitment(
	height int64, seq uint64) (comRes *chantypes.QueryPacketCommitmentResponse, err error) {
	return chanutils.QueryPacketCommitment(c.CLIContext(height), c.PathEnd.PortID, c.PathEnd.ChannelID, seq, true)
}

// QueryPacketAcknowledgement returns the packet ack proof at a given height
func (c *Chain) QueryPacketAcknowledgement(height int64,
	seq uint64) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error) {
	return chanutils.QueryPacketAcknowledgement(c.CLIContext(height), c.PathEnd.PortID, c.PathEnd.ChannelID, seq, true)
}

// QueryPacketCommitments returns an array of packet commitments
func (c *Chain) QueryPacketCommitments(
	offset, limit, height uint64) (comRes *chantypes.QueryPacketCommitmentsResponse, err error) {

	return chanutils.QueryPacketCommitments(c.CLIContext(int64(int(height))), c.PathEnd.PortID, c.PathEnd.ChannelID, offset, limit)

	//qc := chantypes.NewQueryClient(c.CLIContext(int64(height)))
	//return qc.PacketCommitments(context.Background(), &chantypes.QueryPacketCommitmentsRequest{
	//	PortId:    c.PathEnd.PortID,
	//	ChannelId: c.PathEnd.ChannelID,
	//	Pagination: &querytypes.PageRequest{
	//		Offset:     offset,
	//		Limit:      limit,
	//		CountTotal: true,
	//	},
	//})
}

// QueryUnreceivedPackets returns a list of unrelayed packet commitments
func (c *Chain) QueryUnreceivedPackets(height uint64, seqs []uint64) ([]uint64, error) {
	res, err :=  chanutils.QueryUnreceivedPackets(c.CLIContext(int64(height)), c.PathEnd.PortID, c.PathEnd.ChannelID, seqs)
	if err != nil {
		return nil, err
	}
	return res.Sequences, nil
	//qc := chantypes.NewQueryClient(c.CLIContext(int64(height)))
	//res, err := qc.UnreceivedPackets(context.Background(), &chantypes.QueryUnreceivedPacketsRequest{
	//	PortId:                    c.PathEnd.PortID,
	//	ChannelId:                 c.PathEnd.ChannelID,
	//	PacketCommitmentSequences: seqs,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return res.Sequences, nil
}


// QueryDenomTraces returns all the denom traces from a given chain
func (c *Chain) QueryDenomTraces(offset, limit uint64, height int64) (*transfertypes.QueryDenomTracesResponse, error) {
	return transferutils.QueryDenomTraces(c.CLIContext(height), offset, limit)

	//return transfertypes.NewQueryClient(c.CLIContext(height)).DenomTraces(context.Background(),
	//	&transfertypes.QueryDenomTracesRequest{
	//		Pagination: &querytypes.PageRequest{
	//			Key:        []byte(""),
	//			Offset:     offset,
	//			Limit:      limit,
	//			CountTotal: true,
	//		},
	//	})
}


// QueryPacketReceipt returns the packet receipt proof at a given height
func (c *Chain) QueryPacketReceipt(height int64, seq uint64) (recRes *chantypes.QueryPacketReceiptResponse, err error) {
	return chanutils.QueryPacketReceipt(c.CLIContext(height), c.PathEnd.PortID, c.PathEnd.ChannelID, seq, true)
}



// QueryTxs returns an array of transactions given a tag
func (c *Chain) QueryTxs(height uint64, page, limit int, events []string) ([]*ctypes.ResultTx, error) {
	if len(events) == 0 {
		return nil, errors.New("must declare at least one event to search")
	}

	if page <= 0 {
		return nil, errors.New("page must greater than 0")
	}

	if limit <= 0 {
		return nil, errors.New("limit must greater than 0")
	}

	res, err := c.Client.TxSearch(context.Background(), strings.Join(events, " AND "), true, &page, &limit, "")
	if err != nil {
		return nil, err
	}
	return res.Txs, nil
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
	res, err := chanutils.QueryChannels(c.CLIContext(0), offset, limit)
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
