package softtube

import (
	"strings"

	database "github.com/hultan/softtube/internal/softtube.database"
)

func (v *videoList) getColor(video *database.Video) (string, string) {
	if video.Saved {
		return constColorSaved, "Black"
	}
	duration := v.removeInvalidDurations(video.Duration)
	if strings.Trim(duration, " ") == "LIVE" {
		// If duration is LIVE, lets change color to live color
		return constColorLive, "Black"
	} else if duration == "" {
		// If duration is invalid, lets change color to warning
		return constColorWarning, "Black"
	} else if video.Status == constStatusDeleted {
		return constColorDeleted, "Black"
	} else if video.Status == constStatusWatched {
		return constColorWatched, "Black"
	} else if video.Status == constStatusDownloaded {
		return constColorDownloaded, "Black"
	} else if video.Status == constStatusDownloading {
		return constColorDownloading, "Black"
	} else {
		return constColorNotDownloaded, "White"
	}
}
