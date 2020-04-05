package main

import (
	"errors"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

// StatusBar : The status bar of SoftTube
type StatusBar struct {
	Parent     *SoftTube
	VideoCount *gtk.Label
}

// Load : Loads the toolbar
func (s *StatusBar) Load(builder *gtk.Builder) error {
	label, err := getLabel(builder, "statusbar_number_of_videos")
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

func getLabel(builder *gtk.Builder, name string) (*gtk.Label, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if label, ok := obj.(*gtk.Label); ok {
		return label, nil
	}

	return nil, errors.New("not a gtk label")
}
