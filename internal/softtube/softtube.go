package softtube

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
	"github.com/hultan/softtube/internal/softtube.core"
	"github.com/hultan/softtube/internal/softtube.database"
)

//go:embed assets/main.glade
var mainGlade string

//go:embed assets/delete.png
var deleteIcon []byte

//go:embed assets/download.png
var downloadIcon []byte

//go:embed assets/error.png
var errorIcon []byte

//go:embed assets/play.png
var playIcon []byte

//go:embed assets/set_watched.png
var setWatchedIcon []byte

//go:embed assets/set_unwatched.png
var setUnwatchedIcon []byte

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

	b, err := gtk.BuilderNewFromString(mainGlade)
	if err != nil {
		s.Logger.LogError(err)
		return err
	}
	builder := &framework.GtkBuilder{Builder: b}

	win := builder.GetObject("main_window").(*gtk.Window)
	win.SetTitle(s.getWindowTitle())
	win.Maximize()
	_ = win.Connect(
		"destroy", func() {
			gtk.MainQuit()
		},
	)

	_ = win.Connect(
		"key-press-event", func(w *gtk.Window, e *gdk.Event) {
			k := gdk.EventKeyNewFromEvent(e)

			// fmt.Println(k.State())
			// fmt.Println(k.KeyVal())

			if k.State() == 16 && k.KeyVal() == 65474 { // F5
				s.videoList.Refresh("")
			}
			if k.State() == 20 && k.KeyVal() >= 49 && k.KeyVal() <= 54 { // CTRL + 1-5
				s.videoList.switchView(viewType(k.KeyVal() - 48))
			}
			if k.State() == 20 && k.KeyVal() == 102 { // Ctrl + f
				s.searchBar.searchEntry.GrabFocus()
			}
			if k.State() == 20 && k.KeyVal() == 108 { // Ctrl + l
				s.videoList.expandCollapseLog()
			}
			if k.State() == 16 && k.KeyVal() == 65535 { // Del
				if s.videoList.currentView == viewToDelete {
					s.videoList.DeleteWatchedVideos()
				}
			}
			if k.State() == 16 && k.KeyVal() == 65360 { // Home
				s.videoList.scroll.toStart()
			}
			if k.State() == 16 && k.KeyVal() == 65367 { // End
				s.videoList.scroll.toEnd()
			}
			if k.State() == 20 && k.KeyVal() == 65367 { // Ctrl + End
				status := s.toolbar.toolbarKeepScrollToEnd.GetActive()
				s.toolbar.toolbarKeepScrollToEnd.SetActive(!status)
				// s.videoList.keepScrollToEnd = !s.videoList.keepScrollToEnd
			}
			if k.State() == 20 && k.KeyVal() == 113 { // Ctrl + q
				gtk.MainQuit()
			}
			if k.State() == 20 && k.KeyVal() == 65535 { // Ctrl + Del
				s.searchBar.Clear()
			}
		},
	)
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
