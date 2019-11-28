package types

// Artist module event types
const (
	EventTypeCreateArtist = "create_artist"

	// Artist attributes
	AttributeValueCategory  = ModuleName
	AttributeKeyArtistID    = ModuleName + "_id"
	AttributeKeyArtistName  = ModuleName + "_name"
	AttributeKeyArtistOwner = ModuleName + "_owner"
)
