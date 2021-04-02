package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	capabilitykeeper "github.com/ci123chain/ci123chain/pkg/capability/keeper"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	porttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/port/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the IBC channel keeper
type Keeper struct {
	// implements gRPC QueryServer interface
	//types.QueryServer

	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	clientKeeper     types.ClientKeeper
	connectionKeeper types.ConnectionKeeper
	portKeeper       types.PortKeeper
	scopedKeeper     capabilitykeeper.ScopedKeeper
}


// NewKeeper creates a new IBC channel Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey,
	clientKeeper types.ClientKeeper, connectionKeeper types.ConnectionKeeper,
	portKeeper types.PortKeeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		clientKeeper:     clientKeeper,
		connectionKeeper: connectionKeeper,
		portKeeper:       portKeeper,
		scopedKeeper:     scopedKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"/"+types.SubModuleName)
}

// GetChannel returns a channel with a particular identifier binded to a specific port
func (k Keeper) GetChannel(ctx sdk.Context, portID, channelID string) (types.Channel, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.ChannelKey(portID, channelID))
	if bz == nil {
		return types.Channel{}, false
	}

	var channel types.Channel
	k.cdc.MustUnmarshalBinaryBare(bz, &channel)
	return channel, true
}

// SetChannel sets a channel to the store
func (k Keeper) SetChannel(ctx sdk.Context, portID, channelID string, channel types.Channel) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&channel)
	store.Set(host.ChannelKey(portID, channelID), bz)
}

// GenerateChannelIdentifier returns the next channel identifier.
func (k Keeper) GenerateChannelIdentifier(ctx sdk.Context) string {
	nextChannelSeq := k.GetNextChannelSequence(ctx)
	channelID := types.FormatChannelIdentifier(nextChannelSeq)

	nextChannelSeq++
	k.SetNextChannelSequence(ctx, nextChannelSeq)
	return channelID
}


// GetNextChannelSequence gets the next channel sequence from the store.
func (k Keeper) GetNextChannelSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyNextChannelSequence))
	if bz == nil {
		panic("next channel sequence is nil")
	}

	return sdk.BigEndianToUint64(bz)
}

// SetNextChannelSequence sets the next channel sequence to the store.
func (k Keeper) SetNextChannelSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set([]byte(types.KeyNextChannelSequence), bz)
}



// GetNextSequenceSend gets a channel's next send sequence from the store
func (k Keeper) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.NextSequenceSendKey(portID, channelID))
	if bz == nil {
		return 0, false
	}

	return sdk.BigEndianToUint64(bz), true
}

// SetNextSequenceSend sets a channel's next send sequence to the store
func (k Keeper) SetNextSequenceSend(ctx sdk.Context, portID, channelID string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set(host.NextSequenceSendKey(portID, channelID), bz)
}



// GetNextSequenceRecv gets a channel's next receive sequence from the store
func (k Keeper) GetNextSequenceRecv(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.NextSequenceRecvKey(portID, channelID))
	if bz == nil {
		return 0, false
	}

	return sdk.BigEndianToUint64(bz), true
}

// SetNextSequenceRecv sets a channel's next receive sequence to the store
func (k Keeper) SetNextSequenceRecv(ctx sdk.Context, portID, channelID string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set(host.NextSequenceRecvKey(portID, channelID), bz)
}


// GetNextSequenceAck gets a channel's next ack sequence from the store
func (k Keeper) GetNextSequenceAck(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.NextSequenceAckKey(portID, channelID))
	if bz == nil {
		return 0, false
	}

	return sdk.BigEndianToUint64(bz), true
}

// SetNextSequenceAck sets a channel's next ack sequence to the store
func (k Keeper) SetNextSequenceAck(ctx sdk.Context, portID, channelID string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set(host.NextSequenceAckKey(portID, channelID), bz)
}

// LookupModuleByChannel will return the IBCModule along with the capability associated with a given channel defined by its portID and channelID
func (k Keeper) LookupModuleByChannel(ctx sdk.Context, portID, channelID string) (string, *capabilitytypes.Capability, error) {
	modules, cap, err := k.scopedKeeper.LookupModules(ctx, host.ChannelCapabilityPath(portID, channelID))
	if err != nil {
		return "", nil, err
	}

	return porttypes.GetModuleOwner(modules), cap, nil
}

// SetPacketCommitment sets the packet commitment hash to the store
func (k Keeper) SetPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64, commitmentHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketCommitmentKey(portID, channelID, sequence), commitmentHash)
}

// GetPacketCommitment gets the packet commitment hash from the store
func (k Keeper) GetPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketCommitmentKey(portID, channelID, sequence))
	return bz
}


func (k Keeper) deletePacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(host.PacketCommitmentKey(portID, channelID, sequence))
}


// GetPacketReceipt gets a packet receipt from the store
func (k Keeper) GetPacketReceipt(ctx sdk.Context, portID, channelID string, sequence uint64) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketReceiptKey(portID, channelID, sequence))
	if bz == nil {
		return "", false
	}

	return string(bz), true
}

// SetPacketReceipt sets an empty packet receipt to the store
func (k Keeper) SetPacketReceipt(ctx sdk.Context, portID, channelID string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketReceiptKey(portID, channelID, sequence), []byte{byte(1)})
}

// HasPacketAcknowledgement check if the packet ack hash is already on the store
func (k Keeper) HasPacketAcknowledgement(ctx sdk.Context, portID, channelID string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(host.PacketAcknowledgementKey(portID, channelID, sequence))
}

// GetPacketAcknowledgement gets the packet ack hash from the store
func (k Keeper) GetPacketAcknowledgement(ctx sdk.Context, portID, channelID string, sequence uint64) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketAcknowledgementKey(portID, channelID, sequence))
	if bz == nil {
		return nil, false
	}
	return bz, true
}


// SetPacketAcknowledgement sets the packet ack hash to the store
func (k Keeper) SetPacketAcknowledgement(ctx sdk.Context, portID, channelID string, sequence uint64, ackHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketAcknowledgementKey(portID, channelID, sequence), ackHash)
}
