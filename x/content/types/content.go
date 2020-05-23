package types

import (
	"fmt"
	"strings"
	"time"
)

type Content struct {
	Uri  string   `json:"uri" yaml:"uri"`
	Hash string   `json:"hash" yaml:"hash"` // /ipfs/QM.......
	Tags []string `json:"tags" yaml:"tags"`
	Dao  Dao      `json:"dao" yaml:"dao"`
	//Actions   Actions   `json:"actions" yaml:"actions"`
	Actions   []string  `json:"actions" yaml:"actions"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

func NewContent(uri, hash string, dao Dao) *Content {
	return &Content{
		Uri:  uri,
		Hash: hash,
		Tags: []string{},
		Dao:  dao,
		//Actions: NewActions(),
		Actions: []string{},
	}
}

func (c Content) String() string {
	return fmt.Sprintf(`Uri: %s
Hash: %s
Dao: %v
Tags: %s
Actions: %s
CreatedAt: %s`,
		c.Uri, c.Hash, c.Dao, c.Tags, c.Actions, c.CreatedAt,
	)
}

func (c Content) Equals(content Content) bool {
	return c.Uri == content.Uri && c.Hash == content.Hash
}

func (c Content) Validate() error {
	if len(strings.TrimSpace(c.Uri)) == 0 {
		return fmt.Errorf("uri cannot be empty")
	}

	if len(c.Uri) > MaxUriLength {
		return fmt.Errorf("uri cannot be longer than %d characters", MaxUriLength)
	}

	if len(c.Hash) > MaxHashLength {
		return fmt.Errorf("hash cannot be longer than %d characters", MaxHashLength)
	}

	if err := c.Dao.Validate(); err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	return nil
}
