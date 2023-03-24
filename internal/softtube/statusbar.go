package softtube

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softtube/internal/builder"
)

// statusBar : The status bar of SoftTube
type statusBar struct {
	parent     *SoftTube
	videoCount *gtk.Label
}

// Load : Loads the toolbar
func (s *statusBar) Init(builder *builder.Builder) error {
	label := builder.GetObject("statusbar_number_of_videos").(*gtk.Label)
	s.videoCount = label

	return nil
}

// UpdateVideoCount : Update the video count in the status bar
func (s *statusBar) UpdateVideoCount(numberOfVideos int) {
	s.videoCount.SetText(fmt.Sprintf("Number of videos : %d", numberOfVideos))
}
