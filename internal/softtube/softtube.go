package softtube

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
	"github.com/hultan/softtube/internal/softtube.core"
	"github.com/hultan/softtube/internal/softtube.database"
)

// SoftTube : The SoftTube application object
type SoftTube struct {
	Config *core.Config
	Logger *core.Logger
	DB     *database.Database

	toolbar     *toolbar
	statusBar   *statusBar
	menuBar     *menuBar
	popupMenu   *popupMenu
	searchBar   *searchBar
	videoList   *videoList
	activityLog *activityLog
}

// StartApplication : Starts the SoftTube application
func (s *SoftTube) StartApplication() error {
	s.Logger.Log("SoftTube client startup")
	defer s.Logger.Log("SoftTube client shutdown")

	gtk.Init(nil)

	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("main.glade")
	if err != nil {
		s.Logger.LogError(err)
		return err
	}

	win := builder.GetObject("main_window").(*gtk.Window)
	win.SetTitle(s.getWindowTitle())
	win.Maximize()
	_ = win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetIconName("video-display")

	s.setupControls(builder)

	// Show the Window and all of its components.
	win.ShowAll()
	gtk.Main()

	return nil
}

func (s *SoftTube) setupControls(builder *framework.GtkBuilder) {
	// Load tool bar
	s.toolbar = &toolbar{parent: s}
	err := s.toolbar.Load(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Load status bar
	s.statusBar = &statusBar{parent: s}
	err = s.statusBar.Load(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Load menu bar
	s.menuBar = &menuBar{parent: s}
	err = s.menuBar.Load(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Load search bar
	s.searchBar = &searchBar{parent: s}
	err = s.searchBar.Load(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Load video list
	s.videoList = &videoList{parent: s}
	err = s.videoList.Load(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	s.videoList.Refresh("")

	// Load popup menu bar
	s.popupMenu = &popupMenu{parent: s}
	err = s.popupMenu.Load(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Load log
	s.activityLog = &activityLog{parent: s, treeView: s.videoList.treeView}
	s.activityLog.Load(builder)
}

func (s *SoftTube) getWindowTitle() string {
	return constAppTitle + " " + constAppVersion
}
