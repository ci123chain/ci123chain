package channel

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
)

// InitGenesis initializes the ibc channel submodule's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	for _, channel := range gs.Channels {
		ch := types.NewChannel(channel.State, channel.Ordering, channel.Counterparty, channel.ConnectionHops, channel.Version)
		k.SetChannel(ctx, channel.PortId, channel.ChannelId, ch)
	}
	for _, ack := range gs.Acknowledgements {
		k.SetPacketAcknowledgement(ctx, ack.PortId, ack.ChannelId, ack.Sequence, ack.Data)
	}
	for _, commitment := range gs.Commitments {
		k.SetPacketCommitment(ctx, commitment.PortId, commitment.ChannelId, commitment.Sequence, commitment.Data)
	}
	for _, receipt := range gs.Receipts {
		k.SetPacketReceipt(ctx, receipt.PortId, receipt.ChannelId, receipt.Sequence)
	}
	for _, ss := range gs.SendSequences {
		k.SetNextSequenceSend(ctx, ss.PortId, ss.ChannelId, ss.Sequence)
	}
	for _, rs := range gs.RecvSequences {
		k.SetNextSequenceRecv(ctx, rs.PortId, rs.ChannelId, rs.Sequence)
	}
	for _, as := range gs.AckSequences {
		k.SetNextSequenceAck(ctx, as.PortId, as.ChannelId, as.Sequence)
	}
	k.SetNextChannelSequence(ctx, gs.NextChannelSequence)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {

	var ics = make([]types.IdentifiedChannel, 0)
	k.GetChannels(ctx, func(channel types.Channel, portID, channelID string) (stop bool) {
		ic := types.NewIdentifiedChannel(portID, channelID, channel)
		ics = append(ics, ic)
		return false
	})

	var als = make([]types.PacketState, 0)
	k.GetAllPacketAAcknowledgements(ctx, func(data []byte, portID, channelID string, seq uint64) (stop bool) {
		al := types.PacketState{
			PortId:    portID,
			ChannelId: channelID,
			Sequence:  seq,
			Data:      data,
		}
		als = append(als,al)
		return false
	})

	var css = make([]types.PacketState, 0)
	k.GetAllPacketCommitments(ctx, func(data []byte, portID, channelID string, Seq uint64) (stop bool) {
		cs := types.PacketState{
			PortId:    portID,
			ChannelId: channelID,
			Sequence:  Seq,
			Data:      data,
		}
		css = append(css, cs)
		return false
	})

	var rss = make([]types.PacketState, 0)
	k.GetAllPacketReceipts(ctx, func(data []byte, portID, channelID string, Seq uint64) (stop bool) {
		r := types.PacketState{
			PortId:    portID,
			ChannelId: channelID,
			Sequence:  Seq,
			Data:      data,
		}
		rss = append(rss, r)
		return false
	})

	var ss = make([]types.PacketSequence, 0)
	k.GetAllNextSequenceSend(ctx, func(seq uint64, portID, channelID string) (stop bool) {
		s := types.PacketSequence{
			PortId:    portID,
			ChannelId: channelID,
			Sequence:  seq,
		}
		ss = append(ss, s)
		return false
	})

	var recvs = make([]types.PacketSequence, 0)
	k.GetAllNextSequenceRecv(ctx, func(seq uint64, portID, channelID string) (stop bool) {
		recv := types.PacketSequence{
			PortId:    portID,
			ChannelId: channelID,
			Sequence:  seq,
		}
		recvs = append(recvs, recv)
		return false
	})

	var acks = make([]types.PacketSequence, 0)
	k.GetAllNextSequenceAck(ctx, func(seq uint64, portID, channelID string) (stop bool) {
		ack := types.PacketSequence{
			PortId:    "",
			ChannelId: channelID,
			Sequence:  seq,
		}
		acks = append(acks, ack)
		return false
	})

	gs := types.GenesisState{
		Channels:            ics,
		Acknowledgements:    als,
		Commitments:         css,
		Receipts:            rss,
		SendSequences:       ss,
		RecvSequences:       recvs,
		AckSequences:        acks,
		NextChannelSequence: k.GetNextChannelSequence(ctx),
	}
	return gs
}
