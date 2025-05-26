package softtube

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	database "github.com/hultan/softtube/internal/softtube.database"
)

const youtubeDLPath = "yt-dlp"

type videoFunctions struct {
	parent    *SoftTube
	videoList *videoList
}

func (v *videoFunctions) delete(video *database.Video) {
	pathForDeletion := v.getPathForDeletion(video.ID)
	if pathForDeletion == "" {
		return
	}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		// Remove the actual video file
		command := fmt.Sprintf("rm %s", pathForDeletion)
		cmd := exec.Command("/bin/bash", "-c", command)
		// Starts a sub process that deletes the video
		err := cmd.Start()
		if err != nil {
			v.videoList.parent.Logger.Error.Printf("Error deleting video : %v\n", err)
			return
		}
		err = cmd.Wait()
		if err != nil {
			v.videoList.parent.Logger.Error.Printf("Error waiting for process to stop : %v\n", err)
			return
		}
		wg.Done()
	}()

	go func() {
		// Log that the video has been deleted in the database
		err := v.videoList.parent.DB.Log.Insert(constLogDelete, video.Title)
		if err != nil {
			v.videoList.parent.Logger.Error.Println("Failed to log video as watched!")
			v.videoList.parent.Logger.Error.Println(err)
		}
		wg.Done()
	}()

	go func() {
		// Log that the video has been deleted in the GUI
		v.videoList.parent.activityLog.addLog(constLogDelete, video.Title)
		wg.Done()
	}()

	go func() {
		// Set video status as deleted
		err := v.videoList.parent.DB.Videos.UpdateStatus(video.ID, constStatusDeleted)
		if err != nil {
			v.videoList.parent.Logger.Error.Println("Failed to set video status to deleted!")
			v.videoList.parent.Logger.Error.Println(err)
		}
		wg.Done()
	}()

	wg.Wait()
}

func (v *videoFunctions) addToVideoList(video *database.Video, listStore *gtk.ListStore) {
	// Get color based on status
	backgroundColor, foregroundColor := v.videoList.color.getColor(video)
	// Get the duration of the video
	duration := v.videoList.removeInvalidDurations(video.Duration)

	// Get progress
	progress, progressText := v.videoList.getProgress(video.Status)
	// Get thumbnail
	thumbnail := v.getThumbnail(video.ID)

	// Append video to list
	iter := listStore.Append()
	err := listStore.Set(
		iter, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		[]interface{}{
			thumbnail,
			video.SubscriptionName,
			video.Added.Format(constDateLayout),
			video.Title,
			progress,
			string(backgroundColor),
			video.ID,
			duration,
			progressText,
			string(foregroundColor),
			video.Seconds,
		},
	)

	if err != nil {
		v.videoList.parent.Logger.Error.Println("Failed to add row!")
		v.videoList.parent.Logger.Error.Println(err)
	}
}

func (v *videoFunctions) play(video *database.Video) {
	v.videoList.parent.activityLog.fillLog()

	// Mark the selected video with watched color
	v.videoList.color.setRowColor(v.videoList.treeView, constColorWatched)

	// Start Video Player
	v.startSMPlayer(video)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		// Log that the video has been watched in the database
		err := v.videoList.parent.DB.Log.Insert(constLogPlay, video.Title)
		if err != nil {
			v.videoList.parent.Logger.Error.Println("Failed to log video as watched!")
			v.videoList.parent.Logger.Error.Println(err)
		}
		wg.Done()
	}()

	go func() {
		// Log that the video has been watched in the GUI
		v.videoList.parent.activityLog.addLog(constLogPlay, video.Title)
		wg.Done()
	}()

	go func() {
		// Set video status as watched
		err := v.videoList.parent.DB.Videos.UpdateStatus(video.ID, constStatusWatched)
		if err != nil {
			v.videoList.parent.Logger.Error.Println("Failed to set video status to watched!")
			v.videoList.parent.Logger.Error.Println(err)
		}
		wg.Done()
	}()

	wg.Wait()

	v.videoList.Refresh("")

	v.setFocusBackToSoftPlan()
}

func (v *videoFunctions) startSMPlayer(video *database.Video) {
	videoPath := v.getPath(video.ID)
	if videoPath == "" {
		msg := fmt.Sprintf("Failed to find video : %s (%s)", video.Title, video.ID)
		v.videoList.parent.Logger.Error.Println(msg)
		return
	}

	command := fmt.Sprintf("smplayer '%s' &", videoPath)
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
	// Starts a sub process (smplayer)
	// Did not get this to work, but read the following, and maybe I can get
	// this to work in the future
	// https://forum.golangbridge.org/t/starting-new-processes-with-exec-command/24956
	err := cmd.Run()
	if err != nil {
		v.videoList.parent.Logger.Error.Println(err)
	}
}

func (v *videoFunctions) setFocusBackToSoftPlan() {
	go func() {
		time.Sleep(5 * time.Second)

		appClass := "Softtube"

		// Get all matching window IDs
		winIDBytes, err := exec.Command("xdotool", "search", "--class", appClass).Output()
		if err != nil || len(winIDBytes) == 0 {
			fmt.Println("No window found for class:", appClass)
			return
		}

		// Extract valid window IDs
		winIDs := strings.Fields(string(winIDBytes))
		if len(winIDs) > 1 {
			// Activate the second ID, since the first is likely a helper window
			exec.Command("xdotool", "windowmap", winIDs[1]).Run()
			exec.Command("xdotool", "windowactivate", winIDs[1]).Run()
		} else if len(winIDs) == 1 {
			// If only one window is found, activate it
			exec.Command("xdotool", "windowmap", winIDs[0]).Run()
			exec.Command("xdotool", "windowactivate", winIDs[0]).Run()
		}
	}()
}

func (v *videoFunctions) getPath(videoID string) string {
	tryPath := path.Join(v.videoList.parent.Config.ClientPaths.Videos, videoID+".mkv")
	if _, err := os.Stat(tryPath); err == nil {
		return tryPath
	}

	tryPath = path.Join(v.videoList.parent.Config.ClientPaths.Videos, videoID+".mp4")
	if _, err := os.Stat(tryPath); err == nil {
		return tryPath
	}

	tryPath = path.Join(v.videoList.parent.Config.ClientPaths.Videos, videoID+".webm")
	if _, err := os.Stat(tryPath); err == nil {
		return tryPath
	}

	return ""
}

func (v *videoFunctions) getPathForDeletion(videoID string) string {
	tryPath := path.Join(v.videoList.parent.Config.ClientPaths.Videos, videoID+"*")
	return tryPath
}

// Download a YouTube video
func (v *videoFunctions) download(video *database.Video, markAsDownloading bool) error {
	// Set the video to be downloaded
	err := v.videoList.parent.DB.Download.Insert(video.ID)
	if err != nil {
		v.videoList.parent.Logger.Error.Println("Failed to set video to be downloaded!")
		v.videoList.parent.Logger.Error.Println(err)
		return err
	}

	if markAsDownloading {
		// Mark the selected video with downloading color
		v.videoList.color.setRowColor(v.videoList.treeView, constColorDownloading)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		// Log that the video has been requested to be downloaded in the database
		err = v.videoList.parent.DB.Log.Insert(constLogDownload, video.Title)
		if err != nil {
			v.videoList.parent.Logger.Error.Println("Failed to log video as watched!")
			v.videoList.parent.Logger.Error.Println(err)
		}
		wg.Done()
	}()

	go func() {
		// Log that the video has been deleted in the GUI
		v.videoList.parent.activityLog.addLog(constLogDownload, video.Title)
		wg.Done()
	}()

	go func() {
		// Set video status as downloading
		err = v.videoList.parent.DB.Videos.UpdateStatus(video.ID, constStatusDownloading)
		if err != nil {
			v.videoList.parent.Logger.Error.Println("Failed to set video status to downloading!")
			v.videoList.parent.Logger.Error.Println(err)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}

func (v *videoFunctions) getSelectedVideos(treeView *gtk.TreeView) []*database.Video {
	var selectedVideos []*database.Video

	selection, err := treeView.GetSelection()
	if err != nil {
		// Failed to get selection
		return nil
	}

	model, err := treeView.GetModel()
	if err != nil {
		// Failed to get model
		return nil
	}

	paths := selection.GetSelectedRows(model)
	if paths.Length() == 0 {
		// No rows selected
		return nil
	}

	paths.Foreach(func(path interface{}) {
		iter, _ := model.ToTreeModel().GetIter(path.(*gtk.TreePath))
		value, _ := model.ToTreeModel().GetValue(iter, int(listStoreColumnVideoID))
		videoID, err := value.GetString()
		if err != nil {
			return
		}

		for i := 0; i < len(videos); i++ {
			video := videos[i]
			if video.ID == videoID {
				selectedVideos = append(selectedVideos, &video)
				break
			}
		}
	})

	//model, iter, ok := selection.GetSelected()
	//if !ok {
	//	// No rows where selected
	//	return nil
	//}
	//value, err := model.(*gtk.TreeModel).GetValue(iter, int(listStoreColumnVideoID))
	//if err != nil {
	//	return nil
	//}
	//videoID, err := value.GetString()
	//if err != nil {
	//	return nil
	//}
	//for i := 0; i < len(videos); i++ {
	//	video := videos[i]
	//	if video.ID == videoID {
	//		return &video
	//	}
	//}
	return selectedVideos
}

func (v *videoFunctions) setVideoStatus(video *database.Video, mode int, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		var status database.VideoStatusType
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
		err := v.videoList.parent.DB.Videos.UpdateStatus(video.ID, status)
		if err != nil {
			v.videoList.parent.Logger.Error.Printf("Failed to set video as downloaded/watched/unwatched! %s", video.ID)
			v.videoList.parent.Logger.Error.Println(err)
		}
	}()
}

func (v *videoFunctions) setAsSaved(video *database.Video, saved bool) {
	err := v.videoList.parent.DB.Videos.UpdateSave(video.ID, saved)
	if err != nil {
		if saved {
			v.videoList.parent.Logger.Error.Printf("Failed to set video as saved! %s", video.ID)
		} else {
			v.videoList.parent.Logger.Error.Printf("Failed to set video as unsaved! %s", video.ID)
		}
		v.videoList.parent.Logger.Error.Println(err)
	}
	// v.videoList.setRowColor(v.videoList.Treeview, constColorSaved)
	v.videoList.Refresh("")
}

func (v *videoFunctions) getThumbnailPath(videoID string) string {
	thumbnailPath := "/" + path.Join(v.videoList.parent.Config.ClientPaths.Thumbnails, fmt.Sprintf("%s.jpg", videoID))
	if _, err := os.Stat(thumbnailPath); err == nil {
		return thumbnailPath
	}
	thumbnailPath = "/" + path.Join(v.videoList.parent.Config.ClientPaths.Thumbnails, fmt.Sprintf("%s.webp", videoID))
	if _, err := os.Stat(thumbnailPath); err == nil {
		// YouTube started to return *.webp thumbnails instead of *.jpg thumbnails sometimes
		// Go can't read them, and getThumbnail fails to get a PixBuf, so return "" for now
		return thumbnailPath
	}
	return ""
}

func (v *videoFunctions) getThumbnail(videoID string) *gdk.Pixbuf {
	thumbnailPath := v.getThumbnailPath(videoID)
	if thumbnailPath == "" {
		return nil
	}

	thumbnail, err := gdk.PixbufNewFromFile(thumbnailPath)
	if err != nil {
		v.videoList.renameJPG2WEBP(thumbnailPath)
		v.videoList.parent.Logger.Error.Println(err)
		thumbnail = nil
	} else {
		thumbnail, err = thumbnail.ScaleSimple(160, 90, gdk.INTERP_BILINEAR)
		if err != nil {
			msg := fmt.Sprintf("Failed to scale thumnail (%s)!", thumbnailPath)
			v.videoList.parent.Logger.Error.Println(msg)
			thumbnail = nil
		}
	}

	return thumbnail
}

// Download a YouTube video
func (v *videoFunctions) downloadDuration(videoId string, errorChan chan<- error) {
	command := fmt.Sprintf(constVideoDurationCommand, youtubeDLPath, videoId)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("Failed to get duration! Reason:\n\n%s\n%s\n", string(output), err)
		v.parent.Logger.Error.Println(msg)
		errorChan <- errors.New(msg)
		return
	}

	duration := strings.Trim(string(output), " \n")
	if isWarning(duration) {
		i := strings.Index(duration, "\n")
		duration = duration[i+1:]
	}
	if isLive(duration) {
		duration = "LIVE"
	}
	if isError(duration) {
		duration = "ERROR"
	}
	if isMember(duration) {
		duration = "MEMBER"
	}

	err = v.videoList.parent.DB.Videos.UpdateDuration(videoId, duration)
	if err != nil {
		v.videoList.parent.Logger.Error.Println(err.Error())
		errorChan <- err
		return
	}
	errorChan <- nil
}

// Get the thumbnail of a YouTube video
func (v *videoFunctions) downloadThumbnail(videoId string, errorChan chan<- error) {
	// %s/%s.jpg
	thumbPath := fmt.Sprintf(constThumbnailLocation, v.videoList.parent.Config.ServerPaths.Thumbnails, videoId)

	// Don't download thumbnail if it already exists
	if _, err := os.Stat(thumbPath); os.IsExist(err) {
		errorChan <- nil
		return
	}

	command := fmt.Sprintf(constThumbnailCommand, youtubeDLPath, thumbPath, videoId)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("Failed to download thumbnail! Reason:\n\n%s\n%s\n", string(output), err)
		v.parent.Logger.Error.Println(msg)
		errorChan <- errors.New(msg)
		return
	}

	errorChan <- nil
}

func isWarning(duration string) bool {
	if strings.HasPrefix(duration, "WARNING: ") {
		return true
	}

	return false
}

func isLive(duration string) bool {
	if duration == "" || duration == "0" {
		return true
	}

	if strings.Contains(duration, "Premieres") {
		return true
	}

	if strings.Contains(duration, "This live event") {
		return true
	}

	return false
}

func isError(duration string) bool {
	if strings.Contains(duration, "uploader has not") {
		return true
	}

	return false
}

func isMember(duration string) bool {
	if strings.Contains(duration, "Join this channel") {
		return true
	}

	return false
}
