package mkv

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Merge(output string, cont ExtractedContainer) error {
	args := []string{"-o", output}
	for _, track := range cont.Tracks {
		trackIndex := 0

		name := strings.Trim(track.Info.Properties.TrackName, " ")
		args = append(args, "--track-name", fmt.Sprintf("%d:%s", trackIndex, name))

		if len(track.TimeMapPath) > 0 {
			args = append(args, "--timecodes", fmt.Sprintf("%d:%s", trackIndex, track.TimeMapPath))
		}

		if track.Info.Properties.LanguageIETF.String() != "" && track.Info.Properties.LanguageIETF.String() != "und" {
			args = append(args, "--language", fmt.Sprintf("%d:%s", trackIndex, track.Info.Properties.LanguageIETF.String()))
		}

		if track.Info.Properties.DefaultTrack {
			args = append(args, "--default-track-flag", fmt.Sprintf("%d:yes", trackIndex))
		} else {
			args = append(args, "--default-track-flag", fmt.Sprintf("%d:no", trackIndex))
		}

		if track.Info.Properties.ForcedTrack {
			args = append(args, "--forced-display-flag", fmt.Sprintf("%d:yes", trackIndex))
		} else {
			args = append(args, "--forced-display-flag", fmt.Sprintf("%d:no", trackIndex))
		}

		if track.Info.Properties.FlagOriginal {
			args = append(args, "--original-flag", fmt.Sprintf("%d:yes", trackIndex))
		} else {
			args = append(args, "--original-flag", fmt.Sprintf("%d:no", trackIndex))
		}

		if track.Operations.Delay != 0 {
			args = append(args, "--sync", fmt.Sprintf("%d:%d", trackIndex, track.Operations.Delay))
		}

		args = append(args, track.FilePath)
	}

	if cont.Chapters != "" {
		args = append(args, "--chapters", cont.Chapters)
	}

	for _, attachment := range cont.Attachments {
		if attachment.Info.Description != "" {
			args = append(args, "--attachment-description", attachment.Info.Description)
		}
		if attachment.Info.ContentType != "" {
			args = append(args, "--attachment-mime-type", attachment.Info.ContentType)
		}
		if attachment.Info.FileName != "" {
			args = append(args, "--attachment-name", attachment.Info.FileName)
		}
		args = append(args, "--attach-file", attachment.FilePath)
	}

	log.Tracef("Executing mkvmerge with args: %v", args)
	cmd := exec.Command("mkvmerge", args...)
	if logStr, err := cmd.Output(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
			log.Errorf(string(logStr))
			return fmt.Errorf("mkvmerge error: %v", err)
		}
	}

	return nil
}
