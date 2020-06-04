package types

type Artist struct {
	ExternalUrls Externals `json:"external_urls" yaml:"external_urls"` // known external URLs for this artist eg. key: spotify|youtube|soundcloud -> value...
	Followers    uint64    `json:"followers" yaml:"followers"`         // total followers
	Genres       []string  `json:"genres" yaml:"genres"`               // a list of the genres the artist is associated with. For example: "Rock", "Pop"
	Images       []Image   `json:"images" yaml:"images"`
	Name         string    `json:"name" yaml:"name"`
	//Popularity
	//Uri string `json:"uri" yaml:"uri"` // the bitsong uri for the artist e.g.: bitsong:artist:zmsdksd394394
	//ID
}
