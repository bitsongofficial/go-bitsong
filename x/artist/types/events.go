package types

// Artist module event types
const (
	EventTypeCreateArtist   = "create_artist"
	EventTypeSetArtistImage = "set_artist_image"

	// Artist attributes
	AttributeValueCategory  = ModuleName
	AttributeKeyArtistID    = ModuleName + "_id"
	AttributeKeyArtistName  = ModuleName + "_name"
	AttributeKeyArtistImage = ModuleName + "_image"
	AttributeKeyArtistOwner = ModuleName + "_owner"
)
