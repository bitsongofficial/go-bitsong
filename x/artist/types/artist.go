package types

import (
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ArtistI interface {
	//IsActive() bool  // is artist status active?
	GetName() string // get artist name
}

// Implements Artist interface
var _ ArtistI = Artist{}

// nolint - for ArtistI
func (a Artist) GetName() string { return a.Meta.Name }

type Artist struct {
	Meta     Meta           `json:"meta" yaml:"meta"`     // Artist meta
	Images   []Image        `json:"images" yaml:"images"` // Artist images
	ArtistID uint64         `json:"id" yaml:"id"`         // ID of the Artist
	Status   ArtistStatus   `json:"status" yaml:"status"` // Status of the Artist {Nil, Verified, Rejected, Failed}
	Owner    sdk.AccAddress `json:"owner" yaml:"owner"`   // Owner of the Artist`
}

type (
	// ArtistStatus is a type alias that represents an artist status as a byte
	ArtistStatus byte
)

//nolint
const (
	StatusNil      ArtistStatus = 0x00
	StatusVerified ArtistStatus = 0x01
	StatusRejected ArtistStatus = 0x02
	StatusFailed   ArtistStatus = 0x03
)

func NewArtist(id uint64, meta Meta, images []Image, owner sdk.AccAddress) Artist {
	return Artist{
		ArtistID: id,
		Meta:     meta,
		Images:   images,
		Status:   StatusNil,
		Owner:    owner,
	}
}

// nolint
func (a Artist) String() string {
	return fmt.Sprintf(`Artist %d:
  Name:    %s
  Status:  %s
  Owner:   %s`,
		a.ArtistID, a.GetName(), a.Status.String(), a.Owner.String(),
	)
}

// Artists is an array of artist
type Artists []Artist

// nolint
func (a Artists) String() string {
	out := "ID - (Status) Name\n"
	for _, art := range a {
		out += fmt.Sprintf("%d - (%s) %s\n",
			art.ArtistID, art.Status, art.GetName())
	}
	return strings.TrimSpace(out)
}

// ArtistStatusFromString turns a string into a ArtistStatus
func ArtistStatusFromString(str string) (ArtistStatus, error) {
	switch str {
	case "Verified":
		return StatusVerified, nil

	case "Rejected":
		return StatusRejected, nil

	case "Failed":
		return StatusFailed, nil

	case "":
		return StatusNil, nil

	default:
		return ArtistStatus(0xff), fmt.Errorf("'%s' is not a valid artist status", str)
	}
}

// ValidArtistStatus returns true if the artist status is valid and false
// otherwise.
func ValidArtistStatus(status ArtistStatus) bool {
	if status == StatusNil ||
		status == StatusVerified ||
		status == StatusRejected ||
		status == StatusFailed {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status ArtistStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ArtistStatus) Unmarshal(data []byte) error {
	*status = ArtistStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status ArtistStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *ArtistStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ArtistStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// String implements the Stringer interface.
func (status ArtistStatus) String() string {
	switch status {
	case StatusVerified:
		return "Verified"

	case StatusRejected:
		return "Rejected"

	case StatusFailed:
		return "Failed"

	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (status ArtistStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}

// MetaFromProposalType returns a Content object based on the proposal type.
func MetaFromArtist(name string) Meta {
	return NewMeta(name)
}

func ImagesFromArtist(images []Image) []Image {
	var images2 []Image

	for i := 0; i < len(images); i++ {
		images2 = append(images2, NewImage(images[i].CID, images[i].Width, images[i].Height))
	}

	return images2
}
