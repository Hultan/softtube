package main

import (
	"github.com/gotk3/gotk3/gtk"
	gtkhelper "github.com/hultan/softteam/gtk"
	core "github.com/hultan/softtube/softtube.core"
)

// SoftTube : The SoftTube application object
type SoftTube struct {
	Database *core.Database

	Toolbar   *Toolbar
	StatusBar *StatusBar
	MenuBar   *MenuBar
	SearchBar *SearchBar
	VideoList *VideoList
	Log       *Log
}

// StartApplication : Starts the SoftTube application
func (s SoftTube) StartApplication(db *core.Database) error {
	logger.Log("SoftTube client startup")
	defer logger.Log("SoftTube client shutdown")

	s.Database = db

	gtk.Init(nil)

	helper := new(gtkhelper.GtkHelper)

	// Get the path to the glade file
	path, err := helper.GetGladePath()
	if err != nil {
		logger.LogError(err)
		panic(err)
	}

	// Create the builder from the glade file
	builder, err := gtk.BuilderNewFromFile(path)
	if err != nil {
		// panic for any errors.
		logger.LogError(err)
		panic(err)
	}

	win, err := helper.GetWindow(builder, "main_window")
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	win.SetTitle("SoftTube!")
	win.Maximize()
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetIconName("video-display")

	// Load tool bar
	s.Toolbar = &Toolbar{Parent: &s}
	err = s.Toolbar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.Toolbar.SetupEvents()

	// Load status bar
	s.StatusBar = &StatusBar{Parent: &s}
	err = s.StatusBar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}

	// Load menu bar
	s.MenuBar = &MenuBar{Parent: &s}
	err = s.MenuBar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.MenuBar.SetupEvents()

	// Load search bar
	s.SearchBar = &SearchBar{Parent: &s}
	err = s.SearchBar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.SearchBar.SetupEvents()

	// Load video list
	s.VideoList = &VideoList{Parent: &s}
	err = s.VideoList.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.VideoList.SetupColumns()
	s.VideoList.SetupEvents()
	s.VideoList.Refresh("")

	// Load log
	s.Log = &Log{Parent: &s, TreeView: s.VideoList.Treeview}
	s.Log.Load(builder)
	s.Log.FillLog()

	// Show the Window and all of its components.
	win.ShowAll()
	gtk.Main()

	return nil
}
