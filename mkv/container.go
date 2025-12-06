package mkv

import "slices"

type TrackOperations struct {
	Delay int64 // in milliseconds
}

type ExtractedTrack struct {
	Info        Track
	Operations  TrackOperations
	FilePath    string
	TimeMapPath string
}

type ExtractedAttachment struct {
	Info     Attachment
	FilePath string
}

type ExtractedContainer struct {
	Tracks      []ExtractedTrack
	Attachments []ExtractedAttachment
	Chapters    string
}

func (ec *ExtractedContainer) GetDefaultTracks(mainLang LocaleInfo) []int {
	var videoTrack *ExtractedTrack
	audioTracks := make([]ExtractedTrack, 0)
	subtitleTracks := make([]ExtractedTrack, 0)
	for i, lang := range ec.Tracks {
		if lang.Info.Type == "video" {
			videoTrack = &ec.Tracks[i]
		} else if lang.Info.Type == "audio" {
			audioTracks = append(audioTracks, ec.Tracks[i])
		} else if lang.Info.Type == "subtitles" {
			subtitleTracks = append(subtitleTracks, ec.Tracks[i])
		}
	}

	var audioTrack *ExtractedTrack
	if pos := slices.IndexFunc(audioTracks, func(t ExtractedTrack) bool {
		return t.Info.Properties.LanguageIETF == mainLang
	}); pos != -1 {
		audioTrack = &audioTracks[pos]
	} else if pos := slices.IndexFunc(audioTracks, func(t ExtractedTrack) bool {
		baseAudioLang, _ := t.Info.Properties.LanguageIETF.Base()
		baseTargetLang, _ := mainLang.Base()

		return baseAudioLang == baseTargetLang
	}); pos != -1 {
		audioTrack = &audioTracks[pos]
	} else if pos := slices.IndexFunc(audioTracks, func(t ExtractedTrack) bool {
		return t.Info.Properties.DefaultTrack
	}); pos != -1 {
		audioTrack = &audioTracks[pos]
	} else if len(audioTracks) > 0 {
		audioTrack = &audioTracks[0]
	}

	audioInMainLang := audioTrack != nil && audioTrack.Info.Properties.LanguageIETF == mainLang
	var subtitleTrack *ExtractedTrack
	if pos := slices.IndexFunc(subtitleTracks, func(t ExtractedTrack) bool {
		return t.Info.Properties.LanguageIETF == mainLang && t.Info.Properties.ForcedTrack == audioInMainLang
	}); pos != -1 {
		subtitleTrack = &subtitleTracks[pos]
	}

	var selectedIndexes []int
	if videoTrack != nil {
		selectedIndexes = append(selectedIndexes, videoTrack.Info.ID)
	}
	if audioTrack != nil {
		selectedIndexes = append(selectedIndexes, audioTrack.Info.ID)
	}
	if subtitleTrack != nil {
		selectedIndexes = append(selectedIndexes, subtitleTrack.Info.ID)
	}

	return selectedIndexes
}
