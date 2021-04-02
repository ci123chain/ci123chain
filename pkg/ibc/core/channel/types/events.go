package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
)

// IBC channel events vars
var (
	EventTypeChannelOpenInit     = MsgChannelOpenInit{}.MsgType()
	EventTypeChannelOpenTry      = MsgChannelOpenTry{}.MsgType()
	EventTypeChannelOpenAck      = MsgChannelOpenAck{}.MsgType()
	EventTypeChannelOpenConfirm  = MsgChannelOpenConfirm{}.MsgType()
	//EventTypeChannelCloseInit    = MsgChannelCloseInit{}.MsgType()
	//EventTypeChannelCloseConfirm = MsgChannelCloseConfirm{}.MsgType()

	AttributeValueCategory = fmt.Sprintf("%s_%s", host.ModuleName, SubModuleName)
)


// IBC channel events
const (
	AttributeKeyConnectionID       = "connection_id"
	AttributeKeyPortID             = "port_id"
	AttributeKeyChannelID          = "channel_id"
	AttributeCounterpartyPortID    = "counterparty_port_id"
	AttributeCounterpartyChannelID = "counterparty_channel_id"

	EventTypeSendPacket        = "send_packet"
	EventTypeRecvPacket        = "recv_packet"
	EventTypeWriteAck          = "write_acknowledgement"
	EventTypeAcknowledgePacket = "acknowledge_packet"
	EventTypeTimeoutPacket     = "timeout_packet"

	AttributeKeyData             = "packet_data"
	AttributeKeyAck              = "packet_ack"
	AttributeKeyTimeoutHeight    = "packet_timeout_height"
	AttributeKeyTimeoutTimestamp = "packet_timeout_timestamp"
	AttributeKeySequence         = "packet_sequence"
	AttributeKeySrcPort          = "packet_src_port"
	AttributeKeySrcChannel       = "packet_src_channel"
	AttributeKeyDstPort          = "packet_dst_port"
	AttributeKeyDstChannel       = "packet_dst_channel"
	AttributeKeyChannelOrdering  = "packet_channel_ordering"
	AttributeKeyConnection       = "packet_connection"
)
