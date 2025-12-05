package mkv

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
