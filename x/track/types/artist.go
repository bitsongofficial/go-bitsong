package types

/*type Artist struct {
	Cid          string    `json:"cid" yaml:"cid"`   // cid of the track
	Name         string    `json:"name" yaml:"name"` // artist or band name e.g.: The Beatles
	Images       []Image   `json:"images" yaml:"images"`
	Genres       []string  `json:"genres" yaml:"genres"`               // a list of the genres the artist is associated with. For example: "Rock", "Pop"
	ExternalUrls Externals `json:"external_urls" yaml:"external_urls"` // known external URLs for this artist eg. key: spotify|youtube|soundcloud -> value...
	Followers    uint64    `json:"followers" yaml:"followers"`         // total followers
	//Popularity
	//Uri string `json:"uri" yaml:"uri"` // the bitsong uri for the artist e.g.: bitsong:artist:zmsdksd394394
}

func NewArtist(name string, images []Image, genres []string, extUrls Externals) *Artist {
	pref := cid.Prefix{
		Version:  1,
		Codec:    cid.DagCBOR,
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}

	cid, err := pref.Sum([]byte(name)) // TODO: add more data
	if err != nil {
		return nil
	}

	return &Artist{
		Cid:          cid.String(),
		Name:         name,
		Images:       images,
		Genres:       genres,
		ExternalUrls: extUrls,
		Followers:    0,
	}
}

func (a *Artist) String() string {
	// TODO
	return fmt.Sprintf("Name: %s", a.Name)
}

func (a *Artist) Equals(artist Artist) bool {
	// TODO
	return true
}

func (a *Artist) Validate() error {
	// TODO

	if len(strings.TrimSpace(a.Name)) == 0 {
		return fmt.Errorf("artist name cannot be empty")
	}

	return nil
}
*/
