package softtube

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
	core "github.com/hultan/softtube/internal/softtube.core"
	"github.com/hultan/softtube/internal/softtube.database"
)

var (
	config    *core.Config
	logger    *core.Logger
)

// SoftTube : The SoftTube application object
type SoftTube struct {
	Database  *database.Database

	Toolbar   *Toolbar
	StatusBar *StatusBar
	MenuBar   *MenuBar
	PopupMenu *PopupMenu
	SearchBar *SearchBar
	VideoList *VideoList
	Log       *SoftTubeLog
}

// StartApplication : Starts the SoftTube application
func (s *SoftTube) StartApplication(db *database.Database, c *core.Config, l *core.Logger) error {
	logger.Log("SoftTube client startup")
	defer logger.Log("SoftTube client shutdown")

	s.Database = db
	config = c
	logger = l

	gtk.Init(nil)

	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("main.glade")
	if err != nil {
		logger.LogError(err)
		return err
	}

	win := builder.GetObject("main_window").(*gtk.Window)
	win.SetTitle(s.getWindowTitle())
	win.Maximize()
	_ = win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetIconName("video-display")

	// Load tool bar
	s.Toolbar = &Toolbar{Parent: s}
	err = s.Toolbar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.Toolbar.SetupEvents()

	// Load status bar
	s.StatusBar = &StatusBar{Parent: s}
	err = s.StatusBar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}

	// Load menu bar
	s.MenuBar = &MenuBar{Parent: s}
	err = s.MenuBar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.MenuBar.SetupEvents()

	// Load search bar
	s.SearchBar = &SearchBar{Parent: s}
	err = s.SearchBar.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.SearchBar.SetupEvents()

	// Load video list
	s.VideoList = &VideoList{Parent: s}
	err = s.VideoList.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.VideoList.SetupColumns()
	s.VideoList.SetupEvents()
	s.VideoList.Refresh("")

	// Load popup menu bar
	s.PopupMenu = &PopupMenu{Parent: s}
	err = s.PopupMenu.Load(builder)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	s.PopupMenu.SetupEvents()

	// Load log
	s.Log = &SoftTubeLog{Parent: s, TreeView: s.VideoList.TreeView}
	s.Log.Load(builder)
	s.Log.FillLog()

	// Show the Window and all of its components.
	win.ShowAll()
	gtk.Main()

	return nil
}

func (s *SoftTube) getWindowTitle() string {
	return constAppTitle + " " + constAppVersion
}
