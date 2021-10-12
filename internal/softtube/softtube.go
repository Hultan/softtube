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
	// Init toolbar
	s.toolbar = &toolbar{parent: s}
	err := s.toolbar.Init(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Init status bar
	s.statusBar = &statusBar{parent: s}
	err = s.statusBar.Init(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Init menu bar
	s.menuBar = &menuBar{parent: s}
	err = s.menuBar.Init(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Init search bar
	s.searchBar = &searchBar{parent: s}
	err = s.searchBar.Init(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Init video list
	s.videoList = &videoList{parent: s}
	err = s.videoList.Init(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	s.videoList.Refresh("")

	// Init popup menu bar
	s.popupMenu = &popupMenu{parent: s}
	err = s.popupMenu.Init(builder)
	if err != nil {
		s.Logger.LogError(err)
		panic(err)
	}

	// Init log
	s.activityLog = &activityLog{parent: s, treeView: s.videoList.treeView}
	s.activityLog.Init(builder)
}

func (s *SoftTube) getWindowTitle() string {
	return constAppTitle + " " + constAppVersion
}
