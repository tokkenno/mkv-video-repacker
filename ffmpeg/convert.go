package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func Convert(opts ConvertOptions) error {
	// Composite ffmpeg command based on options
	args := []string{"-nostats", "-hide_banner", "-progress", "-"}

	for _, input := range opts.Inputs {
		// Check if file exists and we can read it
		if inf, err := os.Stat(input.Path); os.IsNotExist(err) {
			return err
		} else if inf.Mode().IsRegular() && inf.Mode().Perm()&(1<<(uint(7))) == 0 {
			return errors.New("cannot read input file: " + input.Path)
		}
		args = append(args, "-i", fmt.Sprintf("%s", input.Path))
		if input.TrackMap != "" {
			args = append(args, "-map", input.TrackMap)
		}
	}

	// Add track conversion options
	for _, track := range opts.Tracks {
		args = append(args, "-c:"+track.Index, track.Encoder)
	}

	// Set output path
	args = append(args, fmt.Sprintf("%s", opts.OutputPath))

	// Execute ffmpeg command
	log.WithFields(log.Fields{"process": "ffmpeg"}).Tracef("Executing ffmpeg with args: %v", args)
	cmd := exec.Command("ffmpeg", args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("ffmpeg execution error: %v", exitErr)
		} else {
			return fmt.Errorf("ffmpeg pre-execution error: %v", err)
		}
	}

	// Return the converted track info
	return nil
}
