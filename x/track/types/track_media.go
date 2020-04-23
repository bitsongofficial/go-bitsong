package types

import "fmt"

type Content struct {
	Path        string     `json:"path" yaml:"path"` // /ipfs/Qm.....
	ContentType string     `json:"content_type" yaml:"content_type"`
	Duration    uint8      `json:"duration" yaml:"duration"`
	Attributes  Attributes `json:"attributes" yaml:"attributes"`
}

type TrackMedia struct {
	Audio Content `json:"audio" yaml:"audio"`
	Video Content `json:"video,omitempty" yaml:"video,omitempty"`
	Image Content `json:"image" yaml:"image"`
}

func (tm TrackMedia) Validate() error {
	if !PathRegEx.MatchString(tm.Audio.Path) {
		return fmt.Errorf("track audio is not a valid format")
	}

	if tm.Video.Path != "" {
		if !PathRegEx.MatchString(tm.Video.Path) {
			return fmt.Errorf("track video is not a valid format")
		}
	}

	if !PathRegEx.MatchString(tm.Image.Path) {
		return fmt.Errorf("track image is not a valid format")
	}

	return nil
}

func (tm TrackMedia) Equals(media TrackMedia) bool {
	return tm.Audio.Path == media.Audio.Path && tm.Image.Path == media.Image.Path && tm.Video.Path == media.Video.Path
}

func (tm TrackMedia) String() string {
	return fmt.Sprintf(`Audio: %s, Image: %s, Video: %s`, tm.Audio.Path, tm.Image.Path, tm.Video.Path)
}
