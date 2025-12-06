package mkv

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"
)

func sortedExtractedTracks(tracks []ExtractedTrack) []ExtractedTrack {
	sorted := make([]ExtractedTrack, len(tracks))
	copy(sorted, tracks)

	slices.SortFunc(sorted, func(a, b ExtractedTrack) int {
		if a.Info.Type != b.Info.Type {
			// Order by type: video, audio, subtitles, others
			typeOrder := map[string]int{
				"video":     0,
				"audio":     1,
				"subtitles": 2,
			}
			aOrder, aOk := typeOrder[a.Info.Type]
			bOrder, bOk := typeOrder[b.Info.Type]
			if !aOk {
				aOrder = 99
			}
			if !bOk {
				bOrder = 99
			}
			return aOrder - bOrder
		} else {
			if a.Info.Properties.LanguageIETF != b.Info.Properties.LanguageIETF {
				// Sort by language
				return strings.Compare(a.Info.Properties.LanguageIETF.String(), b.Info.Properties.LanguageIETF.String())
			}

			if a.Info.Properties.ForcedTrack != b.Info.Properties.ForcedTrack {
				// Forced tracks first
				if a.Info.Properties.ForcedTrack {
					return -1
				} else {
					return 1
				}
			}
		}
		return 0
	})

	return sorted
}

func ExtractAll(input string, output string) (*ExtractedContainer, error) {
	identity, err := Scan(input)
	if err != nil {
		return nil, fmt.Errorf("error scanning MKV: %v", err)
	}

	if len(output) == 0 || output == "tmp" {
		log.Debugf("Empty output path, extracting to temp dir")
		output, err = os.MkdirTemp(os.TempDir(), "videorepack_")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(output)
	} else {
		if _, err := os.Stat(output); os.IsNotExist(err) {
			log.Errorf("Output file does not exist: %s", output)
			return nil, fmt.Errorf("output path does not exist: %s", output)
		} else if err != nil {
			log.Errorf("Error checking output path: %v", err)
			return nil, fmt.Errorf("error checking output path: %v", err)
		}
	}

	log.Debugf("Extracting all tracks from MKV: %s", input)

	var tracks []ExtractedTrack
	for _, track := range identity.Tracks {
		outputPath := path.Join(output, fmt.Sprintf("track_%d.%s", track.ID,
			track.Properties.FileExtension()))

		timeMapPath := ""
		if track.Type == "video" {
			timeMapPath = path.Join(output, fmt.Sprintf("track_%d_timemap.txt", track.ID))
		}

		tracks = append(tracks, ExtractedTrack{
			Info:        track,
			FilePath:    outputPath,
			TimeMapPath: timeMapPath,
			Operations:  TrackOperations{},
		})
	}

	// Extraer cada pista
	for _, t := range tracks {
		cmd := exec.Command("mkvextract", "tracks", input,
			fmt.Sprintf("%d:%s", t.Info.ID, t.FilePath))

		log.Tracef("Extracting track %d (type: %s) to %s", t.Info.ID, t.Info.Type, t.FilePath)

		if _, err := cmd.Output(); err != nil {
			return nil, fmt.Errorf("mkvextract error: %v", err)
		}

		if len(t.TimeMapPath) > 0 {
			// Extraer timecodes si es pista de video
			cmdTimeMap := exec.Command("mkvextract", "timecodes_v2", input,
				fmt.Sprintf("%d:%s", t.Info.ID, t.TimeMapPath))

			log.Tracef("Extracting timecodes for track %d to %s", t.Info.ID, t.TimeMapPath)

			if _, err := cmdTimeMap.Output(); err != nil {
				return nil, fmt.Errorf("mkvextract timecodes error: %v", err)
			}
		}
	}

	// Extraer capítulos si existen
	chaptersOut := ""
	if len(identity.Chapters) > 0 {
		chaptersOut = path.Join(output, "chapters.xml")
		cmd := exec.Command("mkvextract", input, "chapters", chaptersOut)

		log.Tracef("Extracting chapters to %s", chaptersOut)
		if _, err := cmd.Output(); err != nil {
			return nil, fmt.Errorf("mkvextract chapters error: %v", err)
		}
	}

	// Extraer ficheros adjuntos si existen
	var attachments = make([]ExtractedAttachment, 0)
	if len(identity.Attachments) > 0 {
		for _, attachment := range identity.Attachments {
			attachmentOut := path.Join(output, fmt.Sprintf("attachment_%d_%s", attachment.ID, attachment.FileName))
			cmd := exec.Command("mkvextract", input, "attachments",
				fmt.Sprintf("%d:%s", attachment.ID, attachmentOut))

			log.Tracef("Extracting attachment %d to %s", attachment.ID, attachmentOut)
			if _, err := cmd.Output(); err != nil {
				return nil, fmt.Errorf("mkvextract attachment error: %v", err)
			}
			attachments = append(attachments, ExtractedAttachment{
				Info:     attachment,
				FilePath: attachmentOut,
			})
		}
	}

	return &ExtractedContainer{
		Tracks:   sortedExtractedTracks(tracks),
		Chapters: chaptersOut,
	}, nil
}
