package cli

const (
	FlagOwner    = "owner"
	FlagName     = "name"
	FlagImages   = "images"
	FlagStatus   = "status"
	flagNumLimit = "limit"
	FlagArtist   = "artist"
)

// ArtistFlags defines the core required fields of a artist. It is used to
// verify that these values are not provided in conjunction with a JSON artist
// file.
var ArtistFlags = []string{
	FlagName,
}
