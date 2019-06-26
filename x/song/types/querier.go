package types

//import "strings"

type QueryResSearch struct {
	Title string `json:"title"`
}

func (r QueryResSearch) String() string {
	return r.Title
}