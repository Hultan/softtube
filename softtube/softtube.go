package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	core "github.com/hultan/softtube/softtube.core"
)

// SoftTube : The SoftTube application object
type SoftTube struct {
	Database *core.Database

	Toolbar   Toolbar
	StatusBar StatusBar
	MenuBar   MenuBar
	SearchBar SearchBar
	VideoList VideoList
}

// StartApplication : Starts the SoftTube application
func (s SoftTube) StartApplication(db *core.Database) error {
	logger.Log("SoftTube client startup")
	defer logger.Log("SoftTube client shutdown")

	s.Database = db

	gtk.Init(nil)

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
	win.SetTitle("SoftTube!")
	win.Maximize()
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Load tool bar
	s.Toolbar = Toolbar{Parent: &s}
	err = s.Toolbar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.Toolbar.SetupEvents()

	// Load status bar
	s.StatusBar = StatusBar{Parent: &s}
	err = s.StatusBar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}

	// Load menu bar
	s.MenuBar = MenuBar{Parent: &s}
	err = s.MenuBar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.MenuBar.SetupEvents()

	// Load search bar
	s.SearchBar = SearchBar{Parent: &s}
	err = s.SearchBar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.SearchBar.SetupEvents()

	// Load video list
	s.VideoList = VideoList{Parent: &s}
	err = s.VideoList.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.VideoList.SetupColumns()
	s.VideoList.SetupEvents()
	s.VideoList.Refresh("")

	// Show the Window and all of its components.
	win.ShowAll()
	gtk.Main()

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
