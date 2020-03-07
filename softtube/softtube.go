package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// SoftTube : The SoftTube application object
type SoftTube struct {
	Application *gtk.Application
	Toolbar     Toolbar
}

// StartApplication : Starts the SoftTube application
func (s SoftTube) StartApplication() error {

	// Create a new application.
	application, err := gtk.ApplicationNew(constAppID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return err
	}
	s.Application = application

	// Connect function to application startup event, this is not required.
	application.Connect("startup", func() {
		logger.Log("application startup")
	})

	return nil
}
