// nolint
package types

const (
	EventTypeIssueFanToken         = "issue_fantoken"
	EventTypeEditFanToken          = "edit_fantoken"
	EventTypeMintFanToken          = "mint_fantoken"
	EventTypeBurnFanToken          = "burn_fantoken"
	EventTypeTransferFanTokenOwner = "transfer_fantoken_owner"

	AttributeValueCategory = ModuleName

	AttributeKeyCreator   = "creator"
	AttributeKeySymbol    = "symbol"
	AttributeKeyDenom     = "denom"
	AttributeKeyAmount    = "amount"
	AttributeKeyOwner     = "owner"
	AttributeKeyDstOwner  = "dst_owner"
	AttributeKeyRecipient = "recipient"
)
