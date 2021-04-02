package types

// IBC transfer events
const (
	EventTypeTimeout      = "timeout"

	EventTypePacket       = "fungible_token_packet"
	EventTypeDenomTrace   = "denomination_trace"
	EventTypeTransfer     = "ibc_transfer"


	AttributeKeyTraceHash      = "trace_hash"
	AttributeKeyReceiver       = "receiver"
	AttributeKeyDenom          = "denom"
	AttributeKeyAmount         = "amount"
	AttributeKeyAckSuccess     = "success"
	AttributeKeyAck            = "acknowledgement"
	AttributeKeyAckError       = "error"
	AttributeKeyRefundReceiver = "refund_receiver"
	AttributeKeyRefundDenom    = "refund_denom"
	AttributeKeyRefundAmount   = "refund_amount"
)