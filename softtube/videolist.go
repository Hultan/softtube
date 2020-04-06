package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	core "github.com/hultan/softtube/softtube.core"
)

// VideoList : The SoftTube video list
type VideoList struct {
	Parent          *SoftTube
	Treeview        *gtk.TreeView
	ScrolledWindow  *gtk.ScrolledWindow
	KeepScrollToEnd bool
}

var videos []core.Video
var filterMode int = 0
var listStore *gtk.ListStore
var filter *gtk.TreeModelFilter

// Load : Loads the toolbar from the glade file
func (v *VideoList) Load(builder *gtk.Builder) error {
	helper := new(GtkHelper)

	// Get the tree view
	treeview, err := helper.GetTreeView(builder, "video_treeview")
	if err != nil {
		return err
	}
	v.Treeview = treeview

	// Get the scrolled window surrounding the treeview
	scroll, err := helper.GetScrolledWindow(builder, "scrolled_window")
	if err != nil {
		return err
	}
	v.ScrolledWindow = scroll

	return nil
}

// SetupEvents : Setup the list events
func (v *VideoList) SetupEvents() {
	// Send in the videolist as a user data parameter to the event
	v.Treeview.Connect("row_activated", rowActivated, v)
}

// SetupColumns : Sets up the listview columns
func (v VideoList) SetupColumns() {
	helper := new(TreeviewHelper)
	v.Treeview.AppendColumn(helper.CreateImageColumn("Image"))
	v.Treeview.AppendColumn(helper.CreateTextColumn("Channel name", liststoreColumnChannelName, 200, 300))
	v.Treeview.AppendColumn(helper.CreateTextColumn("Date", liststoreColumnDate, 90, 300))
	v.Treeview.AppendColumn(helper.CreateTextColumn("Title", liststoreColumnTitle, 0, 600))
	v.Treeview.AppendColumn(helper.CreateTextColumn("Duration", liststoreColumnDuration, 90, 300))
	v.Treeview.AppendColumn(helper.CreateProgressColumn("Progress"))
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
func (v *VideoList) SetFilterMode(mode int) {
	filterMode = mode
	v.Refresh("")
}

// DeleteWatchedVideos : Deletes all watched videos from disk
func (v *VideoList) DeleteWatchedVideos() {
	for i := 0; i < len(videos); i++ {
		video := videos[i]
		if video.Status == constStatusWatched {
			// Delete the video from disk
			deleteVideo(&video)
		}
	}

	v.Refresh("")
}

// Refresh : Refreshes the video list
func (v *VideoList) Refresh(text string) {
	var err error

	db := v.Parent.Database
	if text == "" {
		videos, err = db.Videos.GetVideos()
	} else {
		videos, err = db.Videos.Search(text)
	}
	if err != nil {
		logger.LogError(err)
		panic(err)
	}

	v.Treeview.SetModel(nil)
	listStore, err = gtk.ListStoreNew(gdk.PixbufGetType(), glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_INT64, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		logger.Log("Failed to create liststore!")
		logger.LogError(err)
		panic(err)
	}

	for i := 0; i < len(videos); i++ {
		video := videos[i]
		addVideo(&video, listStore)
	}

	filter, err := listStore.FilterNew(&gtk.TreePath{})
	err = filter.SetVisibleFunc(filterFunc)
	if err != nil {
		logger.LogError(err)
	}
	v.Treeview.SetModel(filter)

	count := filter.IterNChildren(nil)
	v.Parent.StatusBar.UpdateVideoCount(count)

	if v.KeepScrollToEnd {
		// For some reason, we can't scroll to end in the
		// UI thread so create a goroutine that does the
		// scrolling down 50 milliseconds later
		go func() {
			select {
			case <-time.After(50 * time.Millisecond):
				v.ScrollToEnd()
			}
		}()
	}
}

//
// Private functions
//

func filterFunc(model *gtk.TreeModelFilter, iter *gtk.TreeIter, userData interface{}) bool {
	value, err := model.GetValue(iter, liststoreColumnBackground)
	if err != nil {
		// TODO : Log error
	}
	color, err := value.GetString()
	if err != nil {
		// TODO : Log error
	}

	switch filterMode {
	case constFilterModeSubscriptions:
		return true
	case constFilterModeToWatch:
		if color == constColorDownloaded {
			return true
		}
	case constFilterModeToDelete:
		if color == constColorWatched {
			return true
		}
	}

	return false
}

func deleteVideo(video *core.Video) {
	path := getVideoPath(video.ID)
	if path != "" {
		command := fmt.Sprintf("rm %s", path)
		cmd := exec.Command("/bin/bash", "-c", command)
		// Starts a sub process that deletes the video
		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		// Log that the video has been deleted
		err = db.Log.Insert(constLogDeleteVideo, video.Title)
		if err != nil {
			logger.Log("Failed to log video as watched!")
			logger.LogError(err)
		}

		// Set video status as deleted
		err = db.Videos.UpdateStatus(video.ID, constStatusDeleted)
		if err != nil {
			logger.Log("Failed to set video status to deleted!")
			logger.LogError(err)
		}
	}
}

func addVideo(video *core.Video, listStore *gtk.ListStore) {
	// Get color based on status
	backgroundColor, foregroundColor := getColor(video.Status)
	// Get the duration of the video
	duration := getDuration(video.Duration)
	// If duration is invalid, lets change color to warning
	if duration == "" {
		backgroundColor, foregroundColor = setWarningColor()
	}
	// Get progress
	progress, progressText := getProgress(video.Status)
	// Get thumbnail
	thumbnail := getThumbnail(video.ID)

	// Append video to list
	iter := listStore.Append()
	err := listStore.Set(iter, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]interface{}{thumbnail,
			video.SubscriptionName,
			video.Added.Format(constDateLayout),
			video.Title,
			progress,
			backgroundColor,
			video.ID,
			duration,
			progressText,
			foregroundColor})

	if err != nil {
		logger.Log("Failed to add row!")
		logger.LogError(err)
	}
}

func setWarningColor() (string, string) {
	return constColorWarning, "Black"
}

func getDuration(duration sql.NullString) string {
	if duration.Valid && len(strings.Trim(duration.String, " \n")) <= 1 {
		return ""
	}
	return duration.String
}

func getProgress(status int) (int, string) {
	if status == constStatusWatched {
		return 100, "watched"
	} else if status == constStatusDownloaded {
		return 50, "downloaded"
	}

	return 0, ""
}

func getColor(status int) (string, string) {
	if status == constStatusDeleted {
		return constColorDeleted, "Black"
	} else if status == constStatusWatched {
		return constColorWatched, "Black"
	} else if status == constStatusDownloaded {
		return constColorDownloaded, "Black"
	} else if status == constStatusDownloading {
		return constColorDownloading, "Black"
	}
	return constColorNotDownloaded, "White"
}

func getThumbnailPath(videoID string) string {
	// fmt.Println(config.ClientPaths.Thumbnails)
	// fmt.Println(path.Join(config.ClientPaths.Thumbnails, fmt.Sprintf("%s.jpg", videoID)))
	return "/" + path.Join(config.ClientPaths.Thumbnails, fmt.Sprintf("%s.jpg", videoID))
}

func getThumbnail(videoID string) *gdk.Pixbuf {
	path := getThumbnailPath(videoID)

	thumbnail, err := gdk.PixbufNewFromFile(path)
	if err != nil {
		thumbnail = nil
	} else {
		thumbnail, err = thumbnail.ScaleSimple(160, 90, gdk.INTERP_BILINEAR)
		if err != nil {
			msg := fmt.Sprintf("Failed to scale thumnail (%s)!", path)
			logger.Log(msg)
			thumbnail = nil
		}
	}

	return thumbnail
}

func rowActivated(treeView *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn, v *VideoList) {
	video := getSelectedVideo(treeView)
	if video == nil {
		return
	}

	if video.Status == constStatusDownloaded || video.Status == constStatusWatched || video.Status == constStatusSaved {
		playVideo(video)
		// Mark the selected video with watched color
		setRowColor(treeView, constColorWatched)
		v.Refresh("")
	} else if video.Status == constStatusNotDownloaded {
		downloadVideo(video)
		// Mark the selected video with downloading color
		setRowColor(treeView, constColorDownloading)
	}
}

func setRowColor(treeView *gtk.TreeView, color string) {
	selection, _ := treeView.GetSelection()
	rows := selection.GetSelectedRows(listStore)
	path := rows.Data().(*gtk.TreePath)
	iter, _ := listStore.GetIter(path)
	_ = listStore.SetValue(iter, liststoreColumnBackground, color)
}

func playVideo(video *core.Video) {
	path := getVideoPath(video.ID)
	if path == "" {
		msg := fmt.Sprintf("Failed to find video : %s (%s)", video.Title, video.ID)
		logger.Log(msg)
		return
	}
	command := fmt.Sprintf("smplayer '%s'", path)
	cmd := exec.Command("/bin/bash", "-c", command)
	// Starts a sub process (smplayer)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	err = db.Log.Insert(constLogPlayVideo, video.Title)
	if err != nil {
		logger.Log("Failed to log video as watched!")
		logger.LogError(err)
	}

	// Set video status as watched
	err = db.Videos.UpdateStatus(video.ID, constStatusWatched)
	if err != nil {
		logger.Log("Failed to set video status to watched!")
		logger.LogError(err)
	}
}

func getVideoPath(videoID string) string {
	tryPath := path.Join(config.ClientPaths.Videos, videoID+".mkv")
	if _, err := os.Stat(tryPath); err == nil {
		return tryPath
	}

	tryPath = path.Join(config.ClientPaths.Videos, videoID+".mp4")
	if _, err := os.Stat(tryPath); err == nil {
		return tryPath
	}

	tryPath = path.Join(config.ClientPaths.Videos, videoID+".webm")
	if _, err := os.Stat(tryPath); err == nil {
		return tryPath
	}

	return ""
}

// Download a youtube video
func downloadVideo(video *core.Video) error {
	// Set the video to be downloaded
	err := db.Download.Insert(video.ID)
	if err != nil {
		logger.Log("Failed to set video to be downloaded!")
		logger.LogError(err)
		return err
	}

	err = db.Log.Insert(constLogDownloadVideo, video.Title)
	if err != nil {
		logger.Log("Failed to log video as watched!")
		logger.LogError(err)
	}

	// Set video status as downloading
	err = db.Videos.UpdateStatus(video.ID, constStatusDownloading)
	if err != nil {
		logger.Log("Failed to set video status to downloading!")
		logger.LogError(err)
	}

	return nil
}

func getYoutubePath() string {
	return path.Join(config.ServerPaths.YoutubeDL, "youtube-dl")
}

func getSelectedVideo(treeView *gtk.TreeView) *core.Video {
	selection, err := treeView.GetSelection()
	if err != nil {
		return nil
	}
	model, iter, ok := selection.GetSelected()
	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, liststoreColumnVideoID)
		if err != nil {
			return nil
		}
		videoID, err := value.GetString()
		if err != nil {
			return nil
		}
		for i := 0; i < len(videos); i++ {
			video := videos[i]
			if video.ID == videoID {
				return &video
			}
		}
		return nil
	}

	return nil
}
