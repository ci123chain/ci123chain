package types

// IBC transfer events
const (
	EventTypePacket       = "fungible_token_packet"
	EventTypeDenomTrace   = "denomination_trace"
	AttributeKeyTraceHash      = "trace_hash"
	AttributeKeyReceiver       = "receiver"
	AttributeKeyDenom          = "denom"
	AttributeKeyAmount         = "amount"
	AttributeKeyAckSuccess     = "success"
	AttributeKeyAck            = "acknowledgement"
	AttributeKeyAckError       = "error"

)