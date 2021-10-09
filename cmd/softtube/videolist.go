package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hultan/softteam/framework"
	"github.com/hultan/softtube/internal/softtube.database"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const constVideoDurationCommand = "%s --get-duration -- '%s'"

// VideoList : The SoftTube video list
type VideoList struct {
	Parent          *SoftTube
	TreeView        *gtk.TreeView
	ScrolledWindow  *gtk.ScrolledWindow
	KeepScrollToEnd bool
	FilterMode      uint
}

var videos []database.Video
var listStore *gtk.ListStore

// Load : Loads the toolbar from the glade file
func (v *VideoList) Load(builder *framework.GtkBuilder) error {
	v.FilterMode = 0
	// Get the tree view
	treeView := builder.GetObject("video_treeview").(*gtk.TreeView)
	v.TreeView = treeView

	// Get the scrolled window surrounding the treeview
	scroll := builder.GetObject("scrolled_window").(*gtk.ScrolledWindow)
	v.ScrolledWindow = scroll

	return nil
}

// SetupEvents : Setup the list events
func (v *VideoList) SetupEvents() {
	// Send in the videolist as a user data parameter to the event
	_ = v.TreeView.Connect("row_activated", v.rowActivated)
}

// SetupColumns : Sets up the listview columns
func (v VideoList) SetupColumns() {
	helper := new(TreeviewHelper)
	v.TreeView.AppendColumn(helper.CreateImageColumn("Image"))
	v.TreeView.AppendColumn(helper.CreateTextColumn("Channel name", listStoreColumnChannelName, 200, 300))
	v.TreeView.AppendColumn(helper.CreateTextColumn("Date", listStoreColumnDate, 90, 300))
	v.TreeView.AppendColumn(helper.CreateTextColumn("Title", listStoreColumnTitle, 0, 600))
	v.TreeView.AppendColumn(helper.CreateTextColumn("Duration", listStoreColumnDuration, 90, 300))
	v.TreeView.AppendColumn(helper.CreateProgressColumn("Progress"))
}

// Search : Searches for a video
func (v *VideoList) Search(text string) {
	v.Refresh(text)
}

// ScrollToStart : Scrolls to the start of the list
func (v *VideoList) ScrollToStart() {
	var adjustment = v.ScrolledWindow.GetVAdjustment()
	adjustment.SetValue(adjustment.GetLower())
	v.ScrolledWindow.Show()
}

// ScrollToEnd : Scrolls to the end of the list
func (v *VideoList) ScrollToEnd() {
	var adjustment = v.ScrolledWindow.GetVAdjustment()
	adjustment.SetValue(adjustment.GetUpper())
	v.ScrolledWindow.Show()
}

// SetFilterMode : Changes filter mode
func (v *VideoList) SetFilterMode(mode uint) {
	v.FilterMode = mode
	v.Refresh("")
}

// DeleteWatchedVideos : Deletes all watched videos from disk
func (v *VideoList) DeleteWatchedVideos() {
	for i := 0; i < len(videos); i++ {
		video := videos[i]
		if video.Status == constStatusWatched && !video.Saved {
			// Delete the video from disk
			v.deleteVideo(&video)
		}
	}

	v.Refresh("")
}

// Refresh : Refreshes the video list
func (v *VideoList) Refresh(text string) {
	var err error

	db := v.Parent.Database
	if text == "" {
		videos, err = db.Videos.GetVideos(false)
	} else {
		videos, err = db.Videos.Search(text)
	}
	if err != nil {
		logger.LogError(err)
		return
	}

	if listStore != nil {
		listStore.Clear()
	}

	v.TreeView.SetModel(nil)
	listStore, err = gtk.ListStoreNew(gdk.PixbufGetType(), // Thumbnail
		glib.TYPE_STRING, // Subscription name
		glib.TYPE_STRING, // Added date
		glib.TYPE_STRING, // Title
		glib.TYPE_INT64,  // Progress
		glib.TYPE_STRING, // Background color
		glib.TYPE_STRING, // Video ID
		glib.TYPE_STRING, // Duration
		glib.TYPE_STRING, // Progress text
		glib.TYPE_STRING) // Foreground color
	if err != nil {
		logger.Log("Failed to create list store!")
		logger.LogError(err)
		panic(err)
	}

	for i := 0; i < len(videos); i++ {
		video := videos[i]
		v.addVideo(&video, listStore)
	}

	filter, err := listStore.FilterNew(&gtk.TreePath{})
	if err != nil {
		logger.LogError(err)
		return
	}
	filter.SetVisibleFunc(v.filterFunc)
	v.TreeView.SetModel(filter)

	count := filter.IterNChildren(nil)
	v.Parent.StatusBar.UpdateVideoCount(count)

	if v.KeepScrollToEnd {
		// For some reason, we can't scroll to end in the
		// UI thread so create a goroutine that does the
		// scrolling down 50 milliseconds later
		go func() {

			time.Sleep(50 * time.Millisecond)
			v.ScrollToEnd()
			//select {
			//case <-time.After(50 * time.Millisecond):
			//	v.ScrollToEnd()
			//}
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

func (v *VideoList) filterFunc(model *gtk.TreeModel, iter *gtk.TreeIter) bool {
	value, err := model.GetValue(iter, listStoreColumnBackground)
	if err != nil {
		logger.LogError(err)
	}
	color, err := value.GetString()
	if err != nil {
		logger.LogError(err)
	}

	switch v.FilterMode {
	case constFilterModeSubscriptions:
		return true
	case constFilterModeDownloads:
		if color == constColorDownloading {
			return true
		}
	case constFilterModeToWatch:
		if color == constColorDownloaded {
			return true
		}
	case constFilterModeToDelete:
		if color == constColorWatched {
			return true
		}
	case constFilterModeSaved:
		if color == constColorSaved {
			return true
		}
	}

	return false
}

func (v *VideoList) removeInvalidDurations(duration sql.NullString) string {
	if duration.Valid && len(strings.Trim(duration.String, " \n")) <= 1 {
		return ""
	}
	return duration.String
}

func (v *VideoList) getYoutubePath() string {
	return "yt-dlp"
}

func (v *VideoList) getProgress(status int) (int, string) {
	if status == constStatusWatched {
		return 100, "watched"
	} else if status == constStatusDownloaded {
		return 50, "downloaded"
	}

	return 0, ""
}

func (v *VideoList) rowActivated(treeView *gtk.TreeView) {
	video := v.getSelectedVideo(treeView)
	if video == nil {
		return
	}

	if video.Status == constStatusDownloaded || video.Status == constStatusWatched || video.Status == constStatusSaved {
		v.playVideo(video)
	} else if video.Status == constStatusNotDownloaded {
		err := v.downloadVideo(video, true)
		if err != nil {
			logger.LogError(err)
		}
	}
}

// Some .webp images are erroneously named .jpg, so
// rename them so that the converter can take care of them
func (v *VideoList) renameJPG2WEBP(thumbnailPath string) {
	extension := filepath.Ext(thumbnailPath)
	if extension == ".jpg" {
		newName := thumbnailPath[:len(thumbnailPath)-len(extension)] + ".webp"
		_ = os.Rename(thumbnailPath, newName)
	}
}

func (v *VideoList) setRowColor(treeView *gtk.TreeView, color string) {
	selection, _ := treeView.GetSelection()
	rows := selection.GetSelectedRows(listStore)
	treePath := rows.Data().(*gtk.TreePath)
	iter, _ := listStore.GetIter(treePath)
	_ = listStore.SetValue(iter, listStoreColumnBackground, color)
}

//func (v *VideoList) setTooltips() {
//	iter, _ := listStore.GetIterFirst()
//	for ;iter != nil; {
//		videoIDValue, err := listStore.GetValue(iter, 6)
//		if err!=nil {
//			continue
//		}
//		//path, err := listStore.GetPath(iter)
//		//if err!=nil {
//		//	continue
//		//}
//		videoID, err := videoIDValue.GetString()
//		fmt.Println(videoID)
//		//tool := gtk.Tooltip{videoID}
//		//v.TreeView.SetTooltipCell(videoID, path, v.TreeView.GetColumn(0),0)
//		_ = listStore.IterNext(iter)
//	}
//}
