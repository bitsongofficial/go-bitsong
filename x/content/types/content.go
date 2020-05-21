package types

import (
	"fmt"
	"strings"
	"time"
)

type Content struct {
	Uri           string        `json:"uri" yaml:"uri"`
	Hash          []byte        `json:"hash" yaml:"hash"` // /ipfs/QM.......
	Tags          []string      `json:"tags" yaml:"tags"`
	RightsHolders RightsHolders `json:"rights_holders" yaml:"rights_holders"`
	Counters      Counters      `json:"counters" yaml:"counters"`
	CreatedAt     time.Time     `json:"created_at" yaml:"created_at"`
}

func NewContent(uri string, hash []byte, rhs RightsHolders) Content {
	return Content{
		Uri:           uri,
		Hash:          hash,
		Tags:          []string{},
		RightsHolders: rhs,
		Counters:      NewCounters(),
	}
}

func (c Content) String() string {
	return fmt.Sprintf(`Name: %s
Uri: %s
Metadata: %s
ContentUri: %s
Stream Price: %s
Download Price: %s
CreatedAt: %s
Rights Hoders: %s
Total Streams: %d
Total Downloads: %d`,
		c.Name, c.Uri, c.Metadata, c.ContentUri, c.StreamPrice, c.DownloadPrice, c.CreatedAt, c.RightsHolders, c.TotalStreams, c.TotalDownloads,
	)
}

func (c Content) Equals(content Content) bool {
	return c.Name == content.Name && c.Uri == content.Uri && c.Metadata == content.Metadata &&
		c.ContentUri == content.ContentUri && c.StreamPrice == content.StreamPrice && c.DownloadPrice == content.DownloadPrice &&
		c.RightsHolders.Equals(content.RightsHolders)
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

	if len(c.Metadata) > MaxMetadataLength {
		return fmt.Errorf("metadata cannot be longer than %d characters", MaxMetadataLength)
	}

	if len(strings.TrimSpace(c.ContentUri)) == 0 {
		return fmt.Errorf("content-uri cannot be empty")
	}

	if len(c.ContentUri) > MaxUriLength {
		return fmt.Errorf("content-uri cannot be longer than %d characters", MaxUriLength)
	}

	if err := c.RightsHolders.Validate(); err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	return nil
}
