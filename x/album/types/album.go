package types

import (
	btsg "github.com/bitsongofficial/go-bitsong/types"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "2006-01-02"

type BaseAlbum struct {
	Name                 string    `json:"name"`
	Artists              []btsg.ID `json:"artists"`
	AlbumGroup           string    `json:"album_group"`
	AlbumType            string    `json:"album_type"`
	ID                   btsg.ID   `json:"id"`
	Markets              []string  `json:"markets"`
	URLs                 btsg.URLs `json:"urls"`
	ReleaseDate          string    `json:"release_date"`
	ReleaseDatePrecision string    `json:"release_date_precision"` // year, month or day
}

func (ba *BaseAlbum) GetReleaseDate() time.Time {
	if ba.ReleaseDatePrecision == "day" {
		result, _ := time.Parse(DateLayout, ba.ReleaseDate)
		return result
	}
	if ba.ReleaseDatePrecision == "month" {
		ym := strings.Split(ba.ReleaseDate, "-")
		year, _ := strconv.Atoi(ym[0])
		month, _ := strconv.Atoi(ym[1])
		return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}
	year, _ := strconv.Atoi(ba.ReleaseDate)
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}

type Album struct {
	BaseAlbum
	Copyrights []Copyright `json:"copyrights"`
	Genres     []string    `json:"genres"`
	Tracks     []btsg.ID   `json:"tracks"`
	EIDs       btsg.EIDs   `json:"eids"`
}

type AlbumType int

const (
	AlbumTypeAlbum AlbumType = 1 << iota
	AlbumTypeSingle
	AlbumTypeCompilation
)
