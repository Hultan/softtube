package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	gtkHelper "github.com/hultan/softteam/gtk"
	core "github.com/hultan/softtube/softtube.core"
	//_ "golang.org/x/image/webp"
)

// VideoList : The SoftTube video list
type VideoList struct {
	Parent          *SoftTube
	TreeView        *gtk.TreeView
	ScrolledWindow  *gtk.ScrolledWindow
	KeepScrollToEnd bool
	FilterMode      uint
}

var videos []core.Video

//var filterMode int = 0
var listStore *gtk.ListStore
var filter *gtk.TreeModelFilter

// Load : Loads the toolbar from the glade file
func (v *VideoList) Load(helper *gtkHelper.GtkHelper) error {
	v.FilterMode = 0
	// Get the tree view
	treeView, err := helper.GetTreeView("video_treeview")
	if err != nil {
		return err
	}
	v.TreeView = treeView

	// Get the scrolled window surrounding the treeview
	scroll, err := helper.GetScrolledWindow("scrolled_window")
	if err != nil {
		return err
	}
	v.ScrolledWindow = scroll

	return nil
}

// SetupEvents : Setup the list events
func (v *VideoList) SetupEvents() {
	// Send in the videolist as a user data parameter to the event
	_, err := v.TreeView.Connect("row_activated", v.rowActivated)
	if err != nil {
		logger.LogError(err)
	}
}

// SetupColumns : Sets up the listview columns
func (v VideoList) SetupColumns() {
	helper := new(TreeviewHelper)
	v.TreeView.AppendColumn(helper.CreateImageColumn("Image"))
	v.TreeView.AppendColumn(helper.CreateTextColumn("Channel name", liststoreColumnChannelName, 200, 300))
	v.TreeView.AppendColumn(helper.CreateTextColumn("Date", liststoreColumnDate, 90, 300))
	v.TreeView.AppendColumn(helper.CreateTextColumn("Title", liststoreColumnTitle, 0, 600))
	v.TreeView.AppendColumn(helper.CreateTextColumn("Duration", liststoreColumnDuration, 90, 300))
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
		videos, err = db.Videos.GetVideos()
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
	if err!=nil{
		logger.LogError(err)
	}
	err = filter.SetVisibleFunc(v.filterFunc)
	if err != nil {
		logger.LogError(err)
	}
	v.TreeView.SetModel(filter)

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

func (v *VideoList) filterFunc(model *gtk.TreeModelFilter, iter *gtk.TreeIter, userData interface{}) bool {
	value, err := model.GetValue(iter, liststoreColumnBackground)
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

func (v *VideoList) deleteVideo(video *core.Video) {
	pathForDeletion := v.getVideoPathForDeletion(video.ID)
	if pathForDeletion != "" {
		command := fmt.Sprintf("rm %s", pathForDeletion)
		cmd := exec.Command("/bin/bash", "-c", command)
		// Starts a sub process that deletes the video
		err := cmd.Start()
		if err != nil {
			//panic(err)
		}

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			// Log that the video has been deleted in the database
			err = db.Log.Insert(constLogDelete, video.Title)
			if err != nil {
				logger.Log("Failed to log video as watched!")
				logger.LogError(err)
			}
			wg.Done()
		}()

		go func() {
			// Log that the video has been deleted in the GUI
			v.Parent.Log.InsertLog(constLogDelete, video.Title)
			wg.Done()
		}()

		go func() {
			// Set video status as deleted
			err = db.Videos.UpdateStatus(video.ID, constStatusDeleted)
			if err != nil {
				logger.Log("Failed to set video status to deleted!")
				logger.LogError(err)
			}
			wg.Done()
		}()

		wg.Wait()
	}
}

func (v *VideoList) addVideo(video *core.Video, listStore *gtk.ListStore) {
	// Get color based on status
	backgroundColor, foregroundColor := v.getColor(video)
	// Get the duration of the video
	duration := v.getDuration(video.Duration)
	// If duration is invalid, lets change color to warning
	if duration == "" {
		backgroundColor, foregroundColor = v.setWarningColor()
	}
	// Get progress
	progress, progressText := v.getProgress(video.Status)
	// Get thumbnail
	thumbnail := v.getThumbnail(video.ID)

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

func (v *VideoList) setWarningColor() (string, string) {
	return constColorWarning, "Black"
}

func (v *VideoList) getDuration(duration sql.NullString) string {
	if duration.Valid && len(strings.Trim(duration.String, " \n")) <= 1 {
		return ""
	}
	return duration.String
}

func (v *VideoList) getProgress(status int) (int, string) {
	if status == constStatusWatched {
		return 100, "watched"
	} else if status == constStatusDownloaded {
		return 50, "downloaded"
	}

	return 0, ""
}

func (v *VideoList) getColor(video *core.Video) (string, string) {
	if video.Saved {
		return constColorSaved, "Black"
	}

	if video.Status == constStatusDeleted {
		return constColorDeleted, "Black"
	} else if video.Status == constStatusWatched {
		return constColorWatched, "Black"
	} else if video.Status == constStatusDownloaded {
		return constColorDownloaded, "Black"
	} else if video.Status == constStatusDownloading {
		return constColorDownloading, "Black"
	}
	return constColorNotDownloaded, "White"
}

func (v *VideoList) getThumbnailPath(videoID string) string {
	thumbnailPath := "/" + path.Join(config.ClientPaths.Thumbnails, fmt.Sprintf("%s.jpg", videoID))
	if _, err := os.Stat(thumbnailPath); err == nil {
		return thumbnailPath
	}
	thumbnailPath = "/" + path.Join(config.ClientPaths.Thumbnails, fmt.Sprintf("%s.webp", videoID))
	if _, err := os.Stat(thumbnailPath); err == nil {
		// YouTube started to return *.webp thumbnails instead of *.jpg thumbnails sometimes
		// Go can't read them, and getThumbnail fails to get a PixBuf, so return "" for now
		return ""
		//return thumbnailPath
	}
	return ""
}

func (v *VideoList) getThumbnail(videoID string) *gdk.Pixbuf {
	thumbnailPath := v.getThumbnailPath(videoID)
	if thumbnailPath=="" {
		return nil
	}

	thumbnail, err := gdk.PixbufNewFromFile(thumbnailPath)
	if err != nil {
		logger.LogError(err)
		thumbnail = nil
	} else {
		thumbnail, err = thumbnail.ScaleSimple(160, 90, gdk.INTERP_BILINEAR)
		if err != nil {
			msg := fmt.Sprintf("Failed to scale thumnail (%s)!", thumbnailPath)
			logger.Log(msg)
			thumbnail = nil
		}
	}

	return thumbnail
}

func (v *VideoList) setAsWatched(video *core.Video, mode int) {
	var status int
	switch mode {
	case 0:
		status = constStatusNotDownloaded
		break
	case 1:
		status = constStatusWatched
		break
	case 2:
		status = constStatusDownloaded
		break
	}
	err := v.Parent.Database.Videos.UpdateStatus(video.ID, status)
	if err != nil {
		logger.LogFormat("Failed to set video as downloaded/watched/unwatched! %s", video.ID)
		logger.LogError(err)
	}
	// v.setRowColor(v.Treeview, constColorSaved)
	v.Refresh("")
}

func (v *VideoList) setAsSaved(video *core.Video, saved bool) {
	err := v.Parent.Database.Videos.UpdateSave(video.ID, saved)
	if err != nil {
		if saved {
			logger.LogFormat("Failed to set video as saved! %s", video.ID)
		} else {
			logger.LogFormat("Failed to set video as unsaved! %s", video.ID)
		}
		logger.LogError(err)
	}
	// v.setRowColor(v.Treeview, constColorSaved)
	v.Refresh("")
}

func (v *VideoList) rowActivated(treeView *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
	video := v.getSelectedVideo(treeView)
	if video == nil {
		return
	}

	if video.Status == constStatusDownloaded || video.Status == constStatusWatched || video.Status == constStatusSaved {
		v.playVideo(video)
	} else if video.Status == constStatusNotDownloaded {
		err:=v.downloadVideo(video)
		if err!=nil {
			logger.LogError(err)
		}
	}
}

func (v *VideoList) setRowColor(treeView *gtk.TreeView, color string) {
	selection, _ := treeView.GetSelection()
	rows := selection.GetSelectedRows(listStore)
	treePath := rows.Data().(*gtk.TreePath)
	iter, _ := listStore.GetIter(treePath)
	_ = listStore.SetValue(iter, liststoreColumnBackground, color)
}

func (v *VideoList) playVideo(video *core.Video) {
	videoPath := v.getVideoPath(video.ID)
	if videoPath == "" {
		msg := fmt.Sprintf("Failed to find video : %s (%s)", video.Title, video.ID)
		logger.Log(msg)
		return
	}
	command := fmt.Sprintf("smplayer '%s'", videoPath)
	cmd := exec.Command("/bin/bash", "-c", command)
	// Starts a sub process (smplayer)
	err := cmd.Start()
	if err != nil {
		logger.LogError(err)
	}

	// Mark the selected video with watched color
	v.setRowColor(v.TreeView, constColorWatched)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		// Log that the video has been deleted in the database
		err = db.Log.Insert(constLogPlay, video.Title)
		if err != nil {
			logger.Log("Failed to log video as watched!")
			logger.LogError(err)
		}
		wg.Done()
	}()

	go func() {
		// Log that the video has been deleted in the GUI
		v.Parent.Log.InsertLog(constLogPlay, video.Title)
		wg.Done()
	}()

	go func() {
		// Set video status as watched
		err = db.Videos.UpdateStatus(video.ID, constStatusWatched)
		if err != nil {
			logger.Log("Failed to set video status to watched!")
			logger.LogError(err)
		}
		wg.Done()
	}()
	wg.Wait()

	v.Refresh("")
}

func (v *VideoList) getVideoPath(videoID string) string {
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

func (v *VideoList) getVideoPathForDeletion(videoID string) string {
	tryPath := path.Join(config.ClientPaths.Videos, videoID+"*")
	return tryPath
}

// Download a youtube video
func (v *VideoList) downloadVideo(video *core.Video) error {
	// Set the video to be downloaded
	err := db.Download.Insert(video.ID)
	if err != nil {
		logger.Log("Failed to set video to be downloaded!")
		logger.LogError(err)
		return err
	}

	// Mark the selected video with downloading color
	v.setRowColor(v.TreeView, constColorDownloading)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		// Log that the video has been requested to be downloaded in the database
		err = db.Log.Insert(constLogDownload, video.Title)
		if err != nil {
			logger.Log("Failed to log video as watched!")
			logger.LogError(err)
		}
		wg.Done()
	}()

	go func() {
		// Log that the video has been deleted in the GUI
		v.Parent.Log.InsertLog(constLogDownload, video.Title)
		wg.Done()
	}()

	go func() {
		// Set video status as downloading
		err = db.Videos.UpdateStatus(video.ID, constStatusDownloading)
		if err != nil {
			logger.Log("Failed to set video status to downloading!")
			logger.LogError(err)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}

// Not used???
// func (v *VideoList) getYoutubePath() string {
// 	return path.Join(config.ServerPaths.YoutubeDL, "youtube-dl")
// }

func (v *VideoList) getSelectedVideo(treeView *gtk.TreeView) *core.Video {
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
