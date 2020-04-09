package main

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	gtkhelper "github.com/hultan/softteam/gtk"
)

// StatusBar : The status bar of SoftTube
type StatusBar struct {
	Parent     *SoftTube
	VideoCount *gtk.Label
}

// Load : Loads the toolbar
func (s *StatusBar) Load(builder *gtk.Builder) error {
	helper := new(gtkhelper.GtkHelper)

	label, err := helper.GetLabel(builder, "statusbar_number_of_videos")
	if err != nil {
		return err
	}
	s.VideoCount = label

	return nil
}

// UpdateVideoCount : Update the video count in the status bar
func (s *StatusBar) UpdateVideoCount(numberOfVideos int) {
	s.VideoCount.SetText(fmt.Sprintf("Number of videos : %d", numberOfVideos))
}
