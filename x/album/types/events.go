package types

// Album module event types
const (
	EventTypeCreateAlbum = "create_album"

	// Album attributes
	AttributeValueCategory                = ModuleName
	AttributeKeyAlbumID                   = ModuleName + "_id"
	AttributeKeyAlbumTitle                = ModuleName + "_track"
	AttributeKeyAlbumType                 = ModuleName + "_type"
	AttributeKeyAlbumReleaseDate          = ModuleName + "_release_date"
	AttributeKeyAlbumReleaseDatePrecision = ModuleName + "_release_date_precision"
	AttributeKeyAlbumOwner                = ModuleName + "_owner"
)
