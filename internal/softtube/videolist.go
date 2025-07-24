package softtube

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hultan/softtube/internal/softtube.database"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// [youtube] WX-NxWCB_XY: Sign in to confirm you’re not a bot. Use --cookies-from-browser
// or --cookies for the authentication.
// See  https://github.com/yt-dlp/yt-dlp/wiki/FAQ#how-do-i-pass-cookies-to-yt-dlp for how to
// manually pass cookies.
// Also see  https://github.com/yt-dlp/yt-dlp/wiki/Extractors#exporting-youtube-cookies  for
// tips on effectively exporting YouTube cookies
const constVideoDurationCommand = "%s --get-duration --cookies-from-browser firefox -- '%s'"

// videoList is the SoftTube video list
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
var totalDuration int64

// Init initializes the toolbar from the glade file
func (v *videoList) Init() error {
	v.currentView = viewSubscriptions
	v.lastViewSwitch = time.Now()

	// Get the tree view
	v.treeView = GetObject[*gtk.TreeView]("video_treeview")
	v.scroll = &scroll{GetObject[*gtk.ScrolledWindow]("scrolled_window")}

	v.videoFunctions = &videoFunctions{parent: v.parent, videoList: v}
	v.color = &color{v}

	helper := &treeViewHelper{v}
	helper.Setup()

	// LOG PANEL

	// Get the scrolled window surrounding the treeview
	v.logPanel = GetObject[*gtk.ScrolledWindow]("log_scrolled_window")
	v.isLogExpanded = true

	return nil
}

// Search searches for a video
func (v *videoList) Search(text string) {
	v.Refresh(text)
}

// DeleteWatchedVideos deletes all watched videos from disk
func (v *videoList) DeleteWatchedVideos() {
	for i := 0; i < len(videos); i++ {
		vid := videos[i]
		if vid.Status == constStatusWatched && !vid.Saved {
			// Delete the video from the disk
			v.videoFunctions.delete(&vid)
		}
	}

	v.Refresh("")
}

// Refresh refreshes the video list
func (v *videoList) Refresh(searchFor string) {
	var err error

	totalDuration = 0

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
		glib.TYPE_STRING,    // Foreground color
		glib.TYPE_INT64,     // Seconds
	)
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
	v.parent.statusBar.UpdateVideoDuration(totalDuration)

	err = v.setNextSelectedVideo()
	if err != nil {
		v.parent.Logger.Error.Println(err)
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

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

	// Run garbage collect after refreshing the list
	go func() {
		select {
		case <-time.After(50 * time.Millisecond):
			runtime.GC()
		}
	}()

	go func() {
		del := v.parent.DB.Videos.HasVideosToDelete()

		label := "To delete"
		if del {
			label = "* " + label + " *"
		}

		glib.IdleAdd(func() {
			v.parent.toolbar.toolbarToDelete.SetLabel(label)
		})
	}()
}

func (v *videoList) setNextSelectedVideo() error {
	path, err := v.getNextVideoPath()
	if err != nil {
		return err
	}
	if path == nil {
		return errors.New("path is nil")
	}

	selection, err := v.treeView.GetSelection()
	if err != nil {
		return err
	}
	if selection == nil {
		return errors.New("selection is nil")
	}

	// This is to make sure that this code is called
	// from the main thread when this function is executed
	// from a goroutine.
	glib.IdleAdd(func() {
		// Set the selection
		selection.SelectPath(path)
		// Set the cursor to ensure it gets properly activated
		v.treeView.SetCursor(path, nil, false)
		// Ensure the TreeView is focused
		v.treeView.GrabFocus()
	})

	return nil
}

func (v *videoList) getNextVideoPath() (*gtk.TreePath, error) {
	var path *gtk.TreePath
	var err error
	var pathString string

	if v.keepScrollToEnd {
		// Select the last row
		model, err := v.treeView.GetModel() // Get the model
		if err != nil {
			return nil, err
		}
		if model != nil {
			iter, ok := model.ToTreeModel().GetIterFirst()
			if ok { // Ensure there's at least one row
				count := 1
				for model.ToTreeModel().IterNext(iter) { // Count the rows
					count++
				}
				// Get path string to the last row
				pathString = fmt.Sprintf("%d", count-1)
			}
		}
	} else {
		// Get path string to the first row
		pathString = "0"
	}
	path, err = gtk.TreePathNewFromString(pathString)
	if err != nil {
		return nil, err
	}

	return path, nil
}

//
// Private functions
//

func (v *videoList) filterFunc(model *gtk.TreeModel, iter *gtk.TreeIter) bool {
	// Get background color
	value, err := model.GetValue(iter, int(listStoreColumnBackground))
	if err != nil {
		v.parent.Logger.Error.Println(err)
	}
	col, err := value.GetString()
	if err != nil {
		v.parent.Logger.Error.Println(err)
	}

	// Get duration
	durationValue, err := model.GetValue(iter, int(listStoreColumnSeconds))
	if err != nil {
		v.parent.Logger.Error.Println(err)
	}
	durationString, err := durationValue.GoValue()
	if err != nil {
		v.parent.Logger.Error.Println(err)
	}
	duration := durationString.(int64)
	include := false
	switch v.currentView {
	case viewSubscriptions:
		include = true
	case viewDownloads:
		if col == constColorDownloading {
			include = true
		}
	case viewToWatch:
		if col == constColorDownloaded {
			include = true
		}
	case viewToDelete:
		if col == constColorWatched {
			include = true
		}
	case viewSaved:
		if col == constColorSaved {
			include = true
		}
	default:
		panic("unhandled default case")
	}

	if include {
		totalDuration += duration
	}

	return include
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
	selectedVideos := v.videoFunctions.getSelectedVideos(treeView)
	if selectedVideos == nil {
		return
	}

	// We only allow playing or downloading one video at a time
	selectedVideo := selectedVideos[0]

	if selectedVideo.Status == constStatusDownloaded ||
		selectedVideo.Status == constStatusWatched ||
		selectedVideo.Saved {

		v.videoFunctions.play(selectedVideo)

	} else if selectedVideo.Status == constStatusNotDownloaded {
		err := v.videoFunctions.download(selectedVideo, true)
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
