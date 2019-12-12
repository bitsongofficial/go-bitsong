package types

// Artist module event types
const (
	EventTypeCreateArtist    = "create_artist"
	EventTypeSetArtistStatus = "set_artist_status"
	EventTypeDepositArtist   = "deposit_artist"

	// Artist attributes
	AttributeValueCategory   = ModuleName
	AttributeKeyArtistID     = "id"
	AttributeKeyArtistName   = "name"
	AttributeKeyArtistOwner  = "owner"
	AttributeKeyArtistStatus = "status"
)
