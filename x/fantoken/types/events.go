// nolint
package types

const (
	EventTypeIssueFanToken         = "issue_fantoken"
	EventTypeEditFanToken          = "edit_fantoken"
	EventTypeMintFanToken          = "mint_fantoken"
	EventTypeBurnFanToken          = "burn_fantoken"
	EventTypeTransferFanTokenOwner = "transfer_fantoken_owner"

	AttributeValueCategory = ModuleName

	AttributeKeyDenom     = "denom"
	AttributeKeyAmount    = "amount"
	AttributeKeyOwner     = "owner"
	AttributeKeyDstOwner  = "dst_owner"
	AttributeKeyRecipient = "recipient"
)
