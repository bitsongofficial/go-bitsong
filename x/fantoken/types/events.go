// nolint
package types

const (
	EventTypeIssueFanToken         = "issue_fan_token"
	EventTypeEditFanToken          = "edit_fan_token"
	EventTypeMintFanToken          = "mint_fan_token"
	EventTypeBurnFanToken          = "burn_fan_token"
	EventTypeTransferFanTokenOwner = "transfer_fan_token_owner"

	AttributeValueCategory = ModuleName

	AttributeKeyCreator   = "creator"
	AttributeKeySymbol    = "symbol"
	AttributeKeyDenom     = "denom"
	AttributeKeyAmount    = "amount"
	AttributeKeyOwner     = "owner"
	AttributeKeyDstOwner  = "dst_owner"
	AttributeKeyRecipient = "recipient"
)
