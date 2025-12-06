package mkv

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os/exec"
	"strings"
	"time"
)

type Chapter struct {
	NumEntries int `json:"num_entries"`
}

type ContainerProperties struct {
	ContainerType         int       `json:"container_type"`
	Date                  time.Time `json:"date_utc"`
	Duration              uint64    `json:"duration"`
	IsProvidingTimestamps bool      `json:"is_providing_timestamps"`
	MuxingApp             string    `json:"muxing_application"`
	SegmentUID            string    `json:"segment_uid"`
	WritingApp            string    `json:"writing_application"`
	TimestampScale        int       `json:"timestamp_scale"`
}

type Container struct {
	Properties ContainerProperties `json:"properties"`
	Recognized bool                `json:"recognized"`
	Supported  bool                `json:"supported"`
	Type       string              `json:"type"`
}

type AttachmentProperties struct {
	UID big.Int `json:"uid"`
}

type Attachment struct {
	ID          int    `json:"id"`
	ContentType string `json:"content_type"`
	Description string `json:"description"`
	FileName    string `json:"file_name"`
	Size        uint64 `json:"size"`
}

type Identity struct {
	Attachments []Attachment `json:"attachments"`
	Chapters    []Chapter    `json:"chapters"`
	Container   Container    `json:"container"`
	Tracks      []Track      `json:"tracks"`
}

func Scan(input string) (*Identity, error) {
	cmd := exec.Command("mkvmerge", "-J", input)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("mkvmerge -J error: %v", err)
	}

	var identity Identity
	if err := json.Unmarshal(out, &identity); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %v", err)
	}

	patchIdentity(&identity)
	return &identity, nil
}

// patchIdentity applies optional fixes to the scanned Identity
func patchIdentity(identity *Identity) {
	patchLanguagesIdentity(identity)
}

func patchLanguagesIdentity(identity *Identity) {
	for i, _ := range identity.Tracks {
		tr := &identity.Tracks[i]
		if strings.Index(tr.Properties.TrackName, "European Spanish") != -1 {
			tr.Properties.LanguageIETF, _ = FromIETFName("es-ES")
		} else if strings.Index(tr.Properties.TrackName, "Brazilian Portuguese") != -1 {
			tr.Properties.LanguageIETF, _ = FromIETFName("pt-BR")
		} else if strings.Index(tr.Properties.TrackName, "Arabic (Saudi Arabia)") != -1 {
			tr.Properties.LanguageIETF, _ = FromIETFName("ar-SA")
		} else if strings.Index(tr.Properties.TrackName, "Chinese (Taiwan)") != -1 {
			tr.Properties.LanguageIETF, _ = FromIETFName("zh-TW")
		} else if strings.Index(tr.Properties.TrackName, "Chinese (Simplified)") != -1 || strings.Index(tr.Properties.TrackName, "Chinese (Mainland China)") != -1 {
			tr.Properties.LanguageIETF, _ = FromIETFName("zh-CN")
		}
	}
}
