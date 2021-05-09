// nolint
package types

const (
	EventTypeIssueFanToken          = "issue_fan_token"
	EventTypeUpdateFanTokenMintable = "update_fan_token_mintable"
	EventTypeMintFanToken           = "mint_fan_token"
	EventTypeBurnFanToken           = "burn_fan_token"
	EventTypeTransferFanTokenOwner  = "transfer_fan_token_owner"

	AttributeValueCategory = ModuleName

	AttributeKeyCreator   = "creator"
	AttributeKeyDenom     = "denom"
	AttributeKeyAmount    = "amount"
	AttributeKeyOwner     = "owner"
	AttributeKeyDstOwner  = "dst_owner"
	AttributeKeyRecipient = "recipient"
)
