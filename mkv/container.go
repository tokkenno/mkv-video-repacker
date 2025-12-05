package mkv

type TrackOperations struct {
	Delay int64 // in milliseconds
}

type ExtractedTrack struct {
	Info       Track
	Operations TrackOperations
	FilePath   string
}

type ExtractedContainer struct {
	Tracks   []ExtractedTrack
	Chapters string
}
