package softtube

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softtube/internal/builder"
)

// statusBar : The status bar of SoftTube
type statusBar struct {
	parent        *SoftTube
	videoCount    *gtk.Label
	totalDuration *gtk.Label
}

// Init initializes the status bar
func (s *statusBar) Init(builder *builder.Builder) error {
	label := builder.GetObject("statusbar_number_of_videos").(*gtk.Label)
	s.videoCount = label
	label = builder.GetObject("statusbar_total_duration").(*gtk.Label)
	s.totalDuration = label

	return nil
}

// UpdateVideoCount : Update the video count in the status bar
func (s *statusBar) UpdateVideoCount(numberOfVideos int) {
	s.videoCount.SetText(fmt.Sprintf("Number of videos : %d", numberOfVideos))
}

// UpdateVideoDuration : Update the total duration of videos in the status bar
func (s *statusBar) UpdateVideoDuration(duration int64) {
	durationString := formatTime(duration)
	if durationString == "00:00:00" {
		s.totalDuration.SetText("")
	} else {
		s.totalDuration.SetText(fmt.Sprintf("Total duration : %s", durationString))
	}
}

func formatTime(seconds int64) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
