package naming

import (
	"regexp"
	"strconv"
	"strings"
)

// Extract extracts naming information from a given filename based on patterns. Some information can be missing or incomplete.
func Extract(filename string) Name {
	name := Name{
		VideoMetadata: make([]string, 0),
		AudioMetadata: make([][]string, 0),
		Authors:       make([]string, 0),
	}

	// Extract show, using all before the first dash or underscore
	name.Show = filename
	if idx := strings.IndexAny(filename, "({[-_"); idx != -1 {
		name.Show = strings.Trim(filename[:idx], " .")
	}

	// Find episode and season numbers
	name.Episode = -1
	name.Season = -1
	complexEpisodePattern := regexp.MustCompile(`(?i)S(\d{1,2})E(\d{1,3})`)
	if matches := complexEpisodePattern.FindStringSubmatch(filename); len(matches) == 3 {
		name.Season, _ = strconv.Atoi(matches[1])
		name.Episode, _ = strconv.Atoi(matches[2])
	}
	simpleEpisodePattern := regexp.MustCompile(`(?i)(\d{1,2})x(\d{1,3})`)
	if matches := simpleEpisodePattern.FindStringSubmatch(filename); len(matches) == 3 {
		name.Season, _ = strconv.Atoi(matches[1])
		name.Episode, _ = strconv.Atoi(matches[2])
	}
	if name.Episode == -1 {
		namedEpisodePattern := regexp.MustCompile(`(?i)(Episode|Ep)\s?(\d{1,3})`)
		if matches := namedEpisodePattern.FindStringSubmatch(filename); len(matches) == 3 {
			name.Episode, _ = strconv.Atoi(matches[2])
		}
	}
	if name.Season == -1 {
		namedSeasonPattern := regexp.MustCompile(`(?i)(Season)\s?(\d{1,2})`)
		if matches := namedSeasonPattern.FindStringSubmatch(filename); len(matches) == 3 {
			name.Season, _ = strconv.Atoi(matches[2])
		}
	}

	// Back the episode detection with a more general pattern if not found
	if name.Episode == -1 {
		mostSimpleEpisodePattern := regexp.MustCompile(`[\s._-](\d{2,3})[\s._-]`)
		if matches := mostSimpleEpisodePattern.FindStringSubmatch(filename); len(matches) == 2 {
			name.Episode, _ = strconv.Atoi(matches[1])
		}
	}

	// Extract year
	yearPattern := regexp.MustCompile(`\((19|20)\d{2}\)`)
	if matches := yearPattern.FindStringSubmatch(filename); len(matches) > 0 {
		name.Year, _ = strconv.Atoi(matches[0])
	}

	// Extract the extension
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		name.Extension = filename[idx+1:]
	}

	return name
}
