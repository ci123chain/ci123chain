package types


const (
	EventsMsgPreStaking = "pre_staking"
	EventMsgStaking     = "staking"
	EventUndelegate     = "undelegate"
	EventTypeRedelegate           = "redelegate"

	AttributeKeySrcValidator      = "source_validator"
	AttributeKeyDstValidator      = "destination_validator"
	AttributeKeyCompletionTime    = "completion_time"
	AttributeValueCategory        = ModuleName
)
