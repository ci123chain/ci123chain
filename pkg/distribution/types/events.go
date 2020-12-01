package types


const (
	EventTypeModifyWithdrawAddress = "modify_withdraw_address"
	EventTypeWithdrawCommission = "withdraw_commission"
	EventTypeWithdrawRewards    = "withdraw_rewards"
	EventTypeFundCommunityPool  = "fund_community_pool"

	AttributeKeyWithdrawAddress = "withdraw_address"
	AttributeValueCategory = ModuleName
	AttributeKeyValidator       = "validator"
)