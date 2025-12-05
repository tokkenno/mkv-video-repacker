package mkv

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"videorepack/types"
)

type VideoTrackProperties struct {
	DisplayDimensions string `json:"display_dimensions"`
	DisplayUnit       int    `json:"display_unit"`
	PixelDimensions   string `json:"pixel_dimensions"`
}

type AudioTrackProperties struct {
	AudioChannels     int `json:"audio_channels"`
	AudioSamplingFreq int `json:"audio_sampling_frequency"`
}

type SubtitleTrackProperties struct {
	Encoding      string `json:"encoding"`
	TextSubtitles bool   `json:"text_subtitles"`
}

type TrackProperties struct {
	UID                big.Int        `json:"uid"`
	CodecID            string         `json:"codec_id"`
	CodecPrivateData   types.HexBytes `json:"codec_private_data"`
	CodecPrivateLength int            `json:"codec_private_length"`
	DefaultDuration    uint64         `json:"default_duration"`
	DefaultTrack       bool           `json:"default_track"`
	EnabledTrack       bool           `json:"enabled_track"`
	ForcedTrack        bool           `json:"forced_track"`
	FlagOriginal       bool           `json:"flag_original"`
	Language           string         `json:"language"`
	LanguageIETF       LocaleInfo     `json:"language_ietf"`
	MinimumTimestamp   int            `json:"minimum_timestamp"`
	Number             int            `json:"number"`
	Packetizer         string         `json:"packetizer"`
	TrackName          string         `json:"track_name"`
	NumIndexEntries    int            `json:"num_index_entries"`

	VideoTrackProperties
	AudioTrackProperties
	SubtitleTrackProperties
}

// FileExtension returns the suggested file extension for the track based on its codec.
func (tp *TrackProperties) FileExtension() string {
	if strings.Index(tp.CodecID, "V_MPEG4/ISO/AVC") != -1 {
		return "h264"
	} else if strings.Index(tp.CodecID, "V_MPEGH/ISO/HEVC") != -1 {
		return "hevc"
	} else if strings.Index(tp.CodecID, "A_AAC") != -1 {
		return "aac"
	} else if strings.Index(tp.CodecID, "A_AC3") != -1 {
		return "ac3"
	} else if strings.Index(tp.CodecID, "A_EAC3") != -1 {
		return "eac3"
	} else if strings.Index(tp.CodecID, "A_DTS") != -1 {
		return "dts"
	} else if strings.Index(tp.CodecID, "A_FLAC") != -1 {
		return "flac"
	} else if strings.Index(tp.CodecID, "S_TEXT/UTF8") != -1 {
		return "srt"
	}

	return "bin"
}

type Track struct {
	ID         int             `json:"id"`
	Type       string          `json:"type"`
	Codec      string          `json:"codec"`
	Properties TrackProperties `json:"properties"`
}

func (tp *Track) NamingMetadata() []string {
	var metadata []string

	if (tp.Type == "audio" || tp.Type == "subtitles") && (tp.Properties.LanguageIETF.String() != "" && tp.Properties.LanguageIETF.String() != "und") {
		metadata = append(metadata, fmt.Sprintf("%s (%s)", tp.Properties.LanguageIETF.LangEnglishName(), tp.Properties.LanguageIETF.String()))
	}

	if tp.Type == "video" {
		if tp.Properties.DisplayDimensions != "" {
			dimensionParts := strings.Split(tp.Properties.DisplayDimensions, "x")
			if len(dimensionParts) == 2 {
				height, _ := strconv.Atoi(dimensionParts[1])
				if height >= 4320 {
					height = 4320
				} else if height >= 2160 {
					height = 2160
				} else if height >= 1440 {
					height = 1440
				} else if height >= 1080 {
					height = 1080
				} else if height >= 720 {
					height = 720
				} else if height >= 480 {
					height = 480
				} else if height >= 360 {
					height = 360
				} else if height >= 240 {
					height = 240
				} else {
					height = 0
				}
				if height > 0 {
					metadata = append(metadata, fmt.Sprintf("%dp", height))
				}
			}
		}
	}

	if tp.Properties.CodecID != "" {
		metadata = append(metadata, strings.ToUpper(tp.Properties.FileExtension()))
	}

	return metadata
}
