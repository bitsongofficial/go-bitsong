package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"time"
)

type Content struct {
	Name           string        `json:"name" yaml:"name"`
	Uri            string        `json:"uri" yaml:"uri"`
	Metadata       string        `json:"metadata" yaml:"metadata"`       // JSON.stringify()
	ContentUri     string        `json:"content_uri" yaml:"content_uri"` // /ipfs/QM.......
	StreamPrice    sdk.Coin      `json:"stream_price" yaml:"stream_price"`
	DownloadPrice  sdk.Coin      `json:"download_price" yaml:"download_price"`
	RightsHolders  RightsHolders `json:"rights_holders" yaml:"rights_holders"`
	TotalStreams   uint64        `json:"total_streams" yaml:"total_streams"`
	TotalDownloads uint64        `json:"total_downloads" yaml:"total_downloads"`
	CreatedAt      time.Time     `json:"created_at" yaml:"created_at"`

	// Aggiungere tags []string{}
}

func NewContent(name, uri, metadata, contentUri string, streamPrice, downloadPrice sdk.Coin, rhs RightsHolders) Content {
	return Content{
		Name:           name,
		Uri:            uri,
		Metadata:       metadata,
		ContentUri:     contentUri,
		StreamPrice:    streamPrice,
		DownloadPrice:  downloadPrice,
		RightsHolders:  rhs,
		TotalDownloads: 0,
		TotalStreams:   0,
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
