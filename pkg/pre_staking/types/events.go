package types


const (
	EventsMsgPreStaking = "pre_staking"
	EventMsgStaking     = "staking"
	EventUndelegate     = "undelegate"
	EventTypeRedelegate           = "redelegate"
	EventTypeDeploy     = "deploy"

	AttributeKeySrcValidator      = "source_validator"
	AttributeKeyDstValidator      = "destination_validator"
	AttributeKeyCompletionTime    = "completion_time"
	AttributeValueCategory        = ModuleName


	AttributeKeyVaultID  =  "VaultID"
	AttributeKeyContract = "contract"
)
