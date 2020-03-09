package main

import (
	"errors"
	"fmt"
	"os"

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
		logger.Log("SoftTube client startup")
	})

	// Connect function to application shutdown event, this is not required.
	application.Connect("shutdown", func() {
		logger.Log("SoftTube client shutdown")
		logger.LogFinished("softtube client")
		logger.Close()
	})

	// Connect function to application activate event
	application.Connect("activate", func() {
		// Get the GtkBuilder UI definition in the glade file.
		path, err := getGladePath()
		if err != nil {
			logger.LogError(err)
			panic(err)
		}
		builder, err := gtk.BuilderNewFromFile(path)
		if err != nil {
			// panic for any errors.
			logger.LogError(err)
			panic(err)
		}

		win, err := getWindow(builder, "main_window")
		if err != nil {
			logger.LogError(err)
			panic(err)
		}

		toolbar := Toolbar{}
		err = toolbar.Load(builder)
		if err != nil {
			logger.LogError(err)
			panic(err)
		}
		toolbar.SetupEvents()

		// Show the Window and all of its components.
		win.Show()
		application.AddWindow(win)
	})

	// Launch the application
	os.Exit(application.Run(os.Args))

	return nil
}

func getWindow(builder *gtk.Builder, name string) (*gtk.Window, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}

	return nil, errors.New("not a gtk window")
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	// Make type assertion (as per gtk.go).
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func getGladePath() (string, error) {
	path := "resources/main.glade"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = "../resources/main.glade"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = "main.glade"
			if _, err := os.Stat(path); os.IsNotExist(err) {
				errorMessage := fmt.Sprintf("Glade file is missing (%s)", path)
				return "", errors.New(errorMessage)
			}
			return path, nil
		}
		return path, nil
	}
	return path, nil
}
