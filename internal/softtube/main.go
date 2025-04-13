package softtube

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/dialog"
	"github.com/hultan/softtube/internal/logger"

	"github.com/hultan/softtube/internal/builder"
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
	Logger *logger.Logger
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
	s.Logger.Info.Println("SoftTube client startup")
	defer s.Logger.Info.Println("SoftTube client shutdown")

	gtk.Init(nil)

	b := builder.NewBuilder(mainGlade)

	win := b.GetObject("main_window").(*gtk.Window)
	win.SetTitle(s.getWindowTitle())
	win.Maximize()
	_ = win.Connect(
		"destroy", func() {
			gtk.MainQuit()
		},
	)

	_ = win.Connect(
		"key-press-event", func(w *gtk.Window, e *gdk.Event) {
			s.onKeyPressed(e)
		},
	)
	win.SetIconName("video-display")

	s.setupControls(b)

	// Show the Window and all of its components.
	win.ShowAll()

	go func() {
		s.videoList.Refresh("")
	}()

	s.showStats()

	gtk.Main()

	return nil
}

func (s *SoftTube) setupControls(builder *builder.Builder) {
	// Init toolbar
	s.toolbar = &toolbar{parent: s}
	err := s.toolbar.Init(builder)
	if err != nil {
		s.Logger.Error.Println("setupControls : toolbar failed!")
		s.Logger.Error.Println(err)
		panic(err)
	}

	// Init status bar
	s.statusBar = &statusBar{parent: s}
	err = s.statusBar.Init(builder)
	if err != nil {
		s.Logger.Error.Println("setupControls : statusbar failed!")
		s.Logger.Error.Println(err)
		panic(err)
	}

	// Init menu bar
	s.menuBar = &menuBar{parent: s}
	err = s.menuBar.Init(builder)
	if err != nil {
		s.Logger.Error.Println("setupControls : menubar failed!")
		s.Logger.Error.Println(err)
		panic(err)
	}

	// Init search bar
	s.searchBar = &searchBar{parent: s}
	err = s.searchBar.Init(builder)
	if err != nil {
		s.Logger.Error.Println("setupControls : searchbar failed!")
		s.Logger.Error.Println(err)
		panic(err)
	}

	// Init video list
	s.videoList = &videoList{parent: s}
	err = s.videoList.Init(builder)
	if err != nil {
		s.Logger.Error.Println("setupControls : videolist failed!")
		s.Logger.Error.Println(err)
		panic(err)
	}

	// Init popup menu bar
	s.popupMenu = &popupMenu{parent: s}
	err = s.popupMenu.Init(builder)
	if err != nil {
		s.Logger.Error.Println("setupControls : popupmenu failed!")
		s.Logger.Error.Println(err)
		panic(err)
	}

	// Init log
	s.activityLog = &activityLog{parent: s, treeView: s.videoList.treeView}
	err = s.activityLog.Init(builder)
	if err != nil {
		s.Logger.Error.Println("setupControls : activitylog failed!")
		s.Logger.Error.Println(err)
		panic(err)
	}
}

func (s *SoftTube) getWindowTitle() string {
	return constAppTitle + " " + constAppVersion
}

func (s *SoftTube) onKeyPressed(e *gdk.Event) {
	k := gdk.EventKeyNewFromEvent(e)

	ctrl := (k.State() & gdk.CONTROL_MASK) != 0
	special := (k.State() & gdk.MOD2_MASK) != 0 // Used for special keys like F5, DELETE, HOME in X11 etc

	// Control + key
	if ctrl {
		switch k.KeyVal() {
		case gdk.KEY_s: // Ctrl + s
			s.showStats()
		case gdk.KEY_f: // Ctrl + f
			s.searchBar.searchEntry.GrabFocus()
		case gdk.KEY_l: // Ctrl + l
			s.videoList.expandCollapseLog()
		case gdk.KEY_q: // Ctrl + q
			gtk.MainQuit()
		case gdk.KEY_d: // Ctrl + d
			selectedVideos := s.videoList.videoFunctions.getSelectedVideos(s.videoList.treeView)
			if selectedVideos != nil {
				s.downloadDurations(selectedVideos)
			}
		case gdk.KEY_t: // Ctrl + t
			selectedVideos := s.videoList.videoFunctions.getSelectedVideos(s.videoList.treeView)
			if selectedVideos != nil {
				s.downloadThumbnails(selectedVideos)
			}
		case gdk.KEY_Delete: // Ctrl + Del
			s.searchBar.Clear()
		case gdk.KEY_End: // Ctrl + End
			status := s.toolbar.toolbarKeepScrollToEnd.GetActive()
			s.toolbar.toolbarKeepScrollToEnd.SetActive(!status)
		default: // Ctrl + 1-5
			if k.KeyVal() >= gdk.KEY_1 && k.KeyVal() <= gdk.KEY_5 { // Change view
				s.videoList.switchView(viewType(k.KeyVal() - gdk.KEY_0))
			}
		}
	}

	// Special keys
	if special {
		switch k.KeyVal() {
		case gdk.KEY_F5: // F5
			s.videoList.Refresh("")
		case gdk.KEY_Delete: // Del
			if s.videoList.currentView == viewToDelete {
				s.videoList.DeleteWatchedVideos()
			}
		case gdk.KEY_Home: // Home
			s.videoList.scroll.toStart()
		case gdk.KEY_End: // End
			s.videoList.scroll.toEnd()
		}
	}
}

func (s *SoftTube) downloadDurations(selectedVideos []*database.Video) {
	errorChan := make(chan error, len(selectedVideos)) // Buffered channel for errors

	for _, video := range selectedVideos {
		go func() { s.videoList.videoFunctions.downloadDuration(video.ID, errorChan) }()
	}

	// Collect errors from all goroutines
	for i := 0; i < len(selectedVideos); i++ {
		if err := <-errorChan; err != nil {
			_, _ = dialog.Title("Failed to get duration").
				Text("An error occurred while trying to get duration of the video").
				ExtraExpand(err.Error()).ExtraHeight(100).
				Width(500).ErrorIcon().OkButton().Show()
		}
	}

	close(errorChan) // Close the channel when done
}

func (s *SoftTube) downloadThumbnails(selectedVideos []*database.Video) {
	errorChan := make(chan error, len(selectedVideos)) // Buffered channel for errors

	for _, video := range selectedVideos {
		go func() { s.videoList.videoFunctions.downloadThumbnail(video.ID, errorChan) }()
	}

	// Collect errors from all goroutines
	for i := 0; i < len(selectedVideos); i++ {
		if err := <-errorChan; err != nil {
			_, _ = dialog.Title("Failed to get thumbnail").
				Text("An error occurred while trying to get the thumbnail of the video").
				ExtraExpand(err.Error()).ExtraHeight(100).
				Width(500).ErrorIcon().OkButton().Show()
		}
	}

	close(errorChan) // Close the channel when done
}

func (s *SoftTube) onWindowDestroy() {
	s.DB.Close()
}
