package ffmpeg

type TrackConvertOptions struct {
	Index   string
	Encoder string
}

type InputFile struct {
	Path string
	// The mapping of tracks in this input file. See https://trac.ffmpeg.org/wiki/Map
	TrackMap string
}

type ConvertOptions struct {
	Inputs     []InputFile
	OutputPath string
	Tracks     []TrackConvertOptions
}
