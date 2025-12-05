package naming

import (
	"fmt"
	"strings"
)

type Name struct {
	Show          string
	Title         string
	Year          int
	Episode       int
	Season        int
	Origin        string
	VideoMetadata []string
	AudioMetadata [][]string
	Authors       []string
	Extension     string
}

func (n *Name) FileName() string {
	filename := n.Show

	if n.Year > 0 {
		filename += fmt.Sprintf(" (%4d)", n.Year)
	}

	if n.Season >= 0 && n.Episode > 0 {
		filename += fmt.Sprintf(" - s%02de%02d", n.Season, n.Episode)
	} else if n.Episode > 0 {
		filename += fmt.Sprintf(" - e%02d", n.Episode)
	}

	if len(n.Title) > 0 {
		filename += " - " + n.Title
	}

	filename += " "

	if len(n.Origin) > 0 {
		filename += "[" + n.Origin + "]"
	}

	if len(n.VideoMetadata) > 0 {
		filename += "[" + strings.Join(n.VideoMetadata, "; ") + "]"
	}

	for _, v := range n.AudioMetadata {
		if len(v) > 0 {
			filename += "[" + strings.Join(v, "; ") + "]"
		}
	}

	if len(n.Authors) > 0 {
		filename += "[" + strings.Join(n.Authors, "; ") + "]"
	}

	return filename + "." + n.Extension
}
