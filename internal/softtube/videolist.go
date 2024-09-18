package softtube

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hultan/softtube/internal/builder"
	"github.com/hultan/softtube/internal/softtube.database"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const constVideoDurationCommand = "%s --get-duration -- '%s'"

// videoList : The SoftTube video list
type videoList struct {
	parent          *SoftTube
	treeView        *gtk.TreeView
	scroll          *scroll
	videoFunctions  *videoFunctions
	color           *color
	keepScrollToEnd bool
	currentView     viewType
	lastViewSwitch  time.Time
	logPanel        *gtk.ScrolledWindow
	isLogExpanded   bool
}

var videos []database.Video
var listStore *gtk.ListStore

// Init : Loads the toolbar from the glade file
func (v *videoList) Init(builder *builder.Builder) error {
	v.currentView = viewSubscriptions
	v.lastViewSwitch = time.Now()

	// Get the tree view
	treeView := builder.GetObject("video_treeview").(*gtk.TreeView)
	v.treeView = treeView

	// Get the scrolled window surrounding the treeview
	s := builder.GetObject("scrolled_window").(*gtk.ScrolledWindow)
	v.scroll = &scroll{s}

	v.videoFunctions = &videoFunctions{parent: v.parent, videoList: v}
	v.color = &color{v}

	helper := &treeViewHelper{v}
	helper.Setup()

	// LOG PANEL

	// Get the scrolled window surrounding the treeview
	logPanel := builder.GetObject("log_scrolled_window").(*gtk.ScrolledWindow)
	v.logPanel = logPanel
	v.isLogExpanded = true

	return nil
}

// Search : Searches for a video
func (v *videoList) Search(text string) {
	v.Refresh(text)
}

// DeleteWatchedVideos : Deletes all watched videos from disk
func (v *videoList) DeleteWatchedVideos() {
	for i := 0; i < len(videos); i++ {
		vid := videos[i]
		if vid.Status == constStatusWatched && !vid.Saved {
			// Delete the video from disk
			v.videoFunctions.delete(&vid)
		}
	}

	v.Refresh("")
}

// Refresh : Refreshes the video list
func (v *videoList) Refresh(searchFor string) {
	var err error

	if searchFor == "" {
		videos, err = v.parent.DB.Videos.GetVideos(false, v.currentView == viewSaved)
	} else {
		videos, err = v.parent.DB.Videos.Search(searchFor)
	}
	if err != nil {
		v.parent.Logger.Error.Println(err)
		return
	}

	if listStore != nil {
		listStore.Clear()
	}

	v.treeView.SetModel(nil)
	listStore, err = gtk.ListStoreNew(
		gdk.PixbufGetType(), // Thumbnail
		glib.TYPE_STRING,    // Subscription name
		glib.TYPE_STRING,    // Added date
		glib.TYPE_STRING,    // Title
		glib.TYPE_INT64,     // Progress
		glib.TYPE_STRING,    // Background color
		glib.TYPE_STRING,    // Video ID
		glib.TYPE_STRING,    // Duration
		glib.TYPE_STRING,    // Progress text
		glib.TYPE_STRING,
	) // Foreground color
	if err != nil {
		v.parent.Logger.Error.Println(err)
		panic(err)
	}

	for i := 0; i < len(videos); i++ {
		video := videos[i]
		v.videoFunctions.addToVideoList(&video, listStore)
	}

	filter, err := listStore.FilterNew(&gtk.TreePath{})
	if err != nil {
		v.parent.Logger.Error.Println(err)
		return
	}
	filter.SetVisibleFunc(v.filterFunc)
	v.treeView.SetModel(filter)

	count := filter.IterNChildren(nil)
	v.parent.statusBar.UpdateVideoCount(count)

	if v.keepScrollToEnd {
		// For some reason, we can't scroll to end in the
		// UI thread so create a goroutine that does the
		// scrolling down 50 milliseconds later
		go func() {
			// We occasionally get an exception here:
			// fatal error: unexpected signal during runtime execution
			// [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x0]
			// Tried to switch from time.Sleep() to time.After()
			// This did not work either.
			select {
			case <-time.After(50 * time.Millisecond):
				v.scroll.toEnd()
			}

			// time.Sleep(50 * time.Millisecond)
			// v.scroll.toEnd()
		}()
	}

	// Run garbage collect after refreshing list
	go func() {
		select {
		case <-time.After(50 * time.Millisecond):
			runtime.GC()
		}
	}()
}

//
// Private functions
//

func (v *videoList) filterFunc(model *gtk.TreeModel, iter *gtk.TreeIter) bool {
	value, err := model.GetValue(iter, int(listStoreColumnBackground))
	if err != nil {
		v.parent.Logger.Error.Println(err)
	}
	col, err := value.GetString()
	if err != nil {
		v.parent.Logger.Error.Println(err)
	}

	switch v.currentView {
	case viewSubscriptions:
		return true
	case viewDownloads:
		if col == constColorDownloading {
			return true
		}
	case viewToWatch:
		if col == constColorDownloaded {
			return true
		}
	case viewToDelete:
		if col == constColorWatched {
			return true
		}
	case viewSaved:
		if col == constColorSaved {
			return true
		}
	default:
		panic("unhandled default case")
	}

	return false
}

func (v *videoList) removeInvalidDurations(duration sql.NullString) string {
	if duration.Valid && len(strings.Trim(duration.String, " \n")) <= 1 {
		return ""
	}
	return duration.String
}

func (v *videoList) getProgress(status database.VideoStatusType) (int, string) {
	if status == constStatusWatched {
		return 100, "watched"
	} else if status == constStatusDownloaded {
		return 50, "downloaded"
	}

	return 0, ""
}

func (v *videoList) rowActivated(treeView *gtk.TreeView) {
	vid := v.videoFunctions.getSelected(treeView)
	if vid == nil {
		return
	}

	if vid.Status == constStatusDownloaded ||
		vid.Status == constStatusWatched ||
		vid.Status == constStatusSaved {

		v.videoFunctions.play(vid)

	} else if vid.Status == constStatusNotDownloaded {
		err := v.videoFunctions.download(vid, true)
		if err != nil {
			v.parent.Logger.Error.Println(err)
		}
	}
}

// Some .webp images are erroneously named .jpg, so
// rename them so that the converter can take care of them
func (v *videoList) renameJPG2WEBP(thumbnailPath string) {
	extension := filepath.Ext(thumbnailPath)
	if extension == ".jpg" {
		newName := thumbnailPath[:len(thumbnailPath)-len(extension)] + ".webp"
		err := os.Rename(thumbnailPath, newName)
		if err != nil {
			v.parent.Logger.Error.Println(err)
		}
	}
}

func (v *videoList) switchView(view viewType) {
	// lastViewSwitch is used to avoid this function
	// being called recursively when calling SetActive(true)
	since := time.Now().Sub(v.lastViewSwitch).Milliseconds()
	if since < 50 || v.currentView == view {
		return
	}
	v.currentView = view
	v.lastViewSwitch = time.Now()

	v.parent.toolbar.toolbarDeleteAll.SetSensitive(view == viewToDelete)

	v.parent.toolbar.toolbarSubscriptions.SetActive(view == viewSubscriptions)
	v.parent.menuBar.menuViewSubscriptions.SetActive(view == viewSubscriptions)

	v.parent.toolbar.toolbarDownloads.SetActive(view == viewDownloads)
	v.parent.menuBar.menuViewDownloads.SetActive(view == viewDownloads)

	v.parent.toolbar.toolbarToWatch.SetActive(view == viewToWatch)
	v.parent.menuBar.menuViewToWatch.SetActive(view == viewToWatch)

	v.parent.toolbar.toolbarSaved.SetActive(view == viewSaved)
	v.parent.menuBar.menuViewSaved.SetActive(view == viewSaved)

	v.parent.toolbar.toolbarToDelete.SetActive(view == viewToDelete)
	v.parent.menuBar.menuViewToDelete.SetActive(view == viewToDelete)

	v.Refresh("")
}

func (v *videoList) expandCollapseLog() {
	if v.isLogExpanded {
		v.logPanel.SetSizeRequest(60, -1)
	} else {
		v.logPanel.SetSizeRequest(410, -1)
	}
	v.isLogExpanded = !v.isLogExpanded
}
