package types

type AlbumType uint8

const (
	AlbumAlbum AlbumType = iota + 1
	AlbumSingle
	AlbumCompilation
)

var AlbumTypeMap = map[string]AlbumType{
	"album":       AlbumAlbum,
	"single":      AlbumSingle,
	"compilation": AlbumCompilation,
}

type Album struct {
	AlbumType            AlbumType            `json:"album_type" yaml:"album_type"`
	Artists              []Artist             `json:"artists" yaml:"artists"` // the artists who performed the track
	Copyrights           []string             `json:"copyrights" yaml:"copyrights"`
	ExternalIds          Externals            `json:"external_ids" yaml:"external_ids"`   // Known external IDs for the track. eg. key: isrc|ean|upc -> value...
	ExternalUrls         Externals            `json:"external_urls" yaml:"external_urls"` // known external URLs for this artist eg. key: spotify|youtube|soundcloud -> value...
	Genres               []string             `json:"genres" yaml:"genres"`               // a list of the genres the artist is associated with. For example: "Rock", "Pop"
	Images               []Image              `json:"images" yaml:"images"`
	Label                string               `json:"label" yaml:"label"`
	Licensor             string               `json:"licensor" yaml:"licensor"`
	Title                string               `json:"title" yaml:"title"`
	ReleaseDate          string               `json:"release_date" yaml:"release_date"`                     // The date the album was first released, for example "1981-12-15". Depending on the precision, it might be shown as "1981" or "1981-12".
	ReleaseDatePrecision ReleaseDatePrecision `json:"release_date_precision" yaml:"release_date_precision"` // The precision with which release_date value is known: "year" , "month" , or "day".
	Tracks               []Track              `json:"tracks" yaml:"tracks"`
	Year                 uint                 // year of the original release (as supplied to bitsong by the issuer, where data is not available, year of the digital release)
	Explicit             bool                 `json:"explicit" yaml:"explicit"` // parental advisory, explicit content tag, as supplied to bitsong by issuer
	// Popularity
	Duration uint `json:"duration" yaml:"duration"` // the length of the track in milliseconds
	// download
	// subscriptionStreaming
	// Uri string `json:"uri" yaml:"uri"` // the bitsong uri for the artist e.g.: bitsong:artist:zmsdksd394394
}

type ReleaseDatePrecision uint8

const (
	Year ReleaseDatePrecision = iota + 1
	Month
	Day
)

var ReleaseDatePrecisionMap = map[string]ReleaseDatePrecision{
	"year":  Year,
	"month": Month,
	"day":   Day,
}
