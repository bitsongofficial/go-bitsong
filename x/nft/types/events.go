package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	EventTypeTransferNft = "transfer_nft"
	EventTypeNftReceived = "nft_received"

	AttributeKeySender   = "sender"
	AttributeKeyReceiver = "receiver"

	AttributeKeyCollection = "collection"
	AttributeKeyTokenId    = "token_id"
)

func NewNftReceivedEvent(receiver sdk.AccAddress, collection string, tokenId string) sdk.Event {
	return sdk.NewEvent(
		EventTypeNftReceived,
		sdk.NewAttribute(AttributeKeyReceiver, receiver.String()),
		sdk.NewAttribute(AttributeKeyCollection, collection),
		sdk.NewAttribute(AttributeKeyTokenId, tokenId),
	)
}
