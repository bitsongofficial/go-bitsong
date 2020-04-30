package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type Content struct {
	Name       string         `json:"name" yaml:"name"`
	Uri        string         `json:"uri" yaml:"uri"`
	MetaUri    string         `json:"meta_uri" yaml:"meta_uri"`
	ContentUri string         `json:"content_uri" yaml:"content_uri"`
	Denom      string         `json:"denom" yaml:"denom"`
	CreatedAt  string         `json:"created_at" yaml:"created_at"`
	Creator    sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewContent(name, uri, metaUri, contentUri, denom string, creator sdk.AccAddress) Content {
	return Content{
		Name:       name,
		Uri:        uri,
		MetaUri:    metaUri,
		ContentUri: contentUri,
		Denom:      denom,
		Creator:    creator,
	}
}

func (c Content) String() string {
	return fmt.Sprintf(`Name: %s
Uri: %s
MetaUri: %s
ContentUri: %s
Denom: %s
CreatedAt: %s
Creator: %s`,
		c.Name, c.Uri, c.MetaUri, c.ContentUri, c.Denom, c.CreatedAt, c.Creator,
	)
}

func (c Content) Equals(content Content) bool {
	return c.Name == content.Name && c.Uri == content.Uri && c.MetaUri == content.MetaUri &&
		c.ContentUri == content.ContentUri && c.Denom == content.Denom
}

func (c Content) Validate() error {
	if len(strings.TrimSpace(c.Name)) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	if len(c.Name) > MaxNameLength {
		return fmt.Errorf("name cannot be longer than %d characters", MaxUriLength)
	}

	if len(strings.TrimSpace(c.Uri)) == 0 {
		return fmt.Errorf("uri cannot be empty")
	}

	if len(c.Uri) > MaxUriLength {
		return fmt.Errorf("uri cannot be longer than %d characters", MaxUriLength)
	}

	if len(strings.TrimSpace(c.MetaUri)) == 0 {
		return fmt.Errorf("meta-uri cannot be empty")
	}

	if len(c.MetaUri) > MaxUriLength {
		return fmt.Errorf("meta-uri cannot be longer than %d characters", MaxUriLength)
	}

	if len(strings.TrimSpace(c.ContentUri)) == 0 {
		return fmt.Errorf("content-uri cannot be empty")
	}

	if len(c.ContentUri) > MaxUriLength {
		return fmt.Errorf("content-uri cannot be longer than %d characters", MaxUriLength)
	}

	if err := sdk.ValidateDenom(c.Denom); err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	if c.Creator == nil {
		return fmt.Errorf("invalid creator: %s", c.Creator)
	}

	return nil
}
