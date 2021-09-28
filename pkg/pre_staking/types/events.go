package types


const (
	EventsMsgPreStaking = "pre_staking"
	EventMsgStaking     = "staking"
	EventUndelegate     = "undelegate"
	EventTypeRedelegate           = "redelegate"
	EventTypeDeploy     = "deploy"
	EventTypeCreateValidator  = "create_validator"

	AttributeKeySrcValidator      = "source_validator"
	AttributeKeyDstValidator      = "destination_validator"
	AttributeKeyCompletionTime    = "completion_time"
	AttributeValueCategory        = ModuleName


	AttributeKeyVaultID  =  "VaultID"
	AttributeKeyContract = "contract"
)
