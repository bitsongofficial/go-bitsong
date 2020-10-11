package types

type BaseArtist struct {
	Name string `json:"name"`
	ID   ID     `json:"id"`
	URLs URLs   `json:"urls"`
}

type Artist struct {
	BaseArtist
	Genres []string `json:"genres"`
	//Followers
	//Images
}
