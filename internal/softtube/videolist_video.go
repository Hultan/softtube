package softtube

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	database "github.com/hultan/softtube/internal/softtube.database"
)

func (v *videoList) deleteVideo(video *database.Video) {
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
			err = v.parent.db.Log.Insert(constLogDelete, video.Title)
			if err != nil {
				v.parent.logger.Log("Failed to log video as watched!")
				v.parent.logger.LogError(err)
			}
			wg.Done()
		}()

		go func() {
			// Log that the video has been deleted in the GUI
			v.parent.activityLog.InsertLog(constLogDelete, video.Title)
			wg.Done()
		}()

		go func() {
			// Set video status as deleted
			err = v.parent.db.Videos.UpdateStatus(video.ID, constStatusDeleted)
			if err != nil {
				v.parent.logger.Log("Failed to set video status to deleted!")
				v.parent.logger.LogError(err)
			}
			wg.Done()
		}()

		wg.Wait()
	}
}

func (v *videoList) addVideo(video *database.Video, listStore *gtk.ListStore) {
	// Get color based on status
	backgroundColor, foregroundColor := v.getColor(video)
	// Get the duration of the video
	duration := v.removeInvalidDurations(video.Duration)

	// Get progress
	progress, progressText := v.getProgress(video.Status)
	// Get thumbnail
	thumbnail := v.getVideoThumbnail(video.ID)

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
		v.parent.logger.Log("Failed to add row!")
		v.parent.logger.LogError(err)
	}
}

func (v *videoList) playVideo(video *database.Video) {
	videoPath := v.getVideoPath(video.ID)
	if videoPath == "" {
		msg := fmt.Sprintf("Failed to find video : %s (%s)", video.Title, video.ID)
		v.parent.logger.Log(msg)
		return
	}

	command := fmt.Sprintf("smplayer '%s'", videoPath)
	cmd := exec.Command("/bin/bash", "-c", command)

	// Starts a sub process (smplayer)
	// Use run (since we are using a go routine), otherwise use Start and Wait together
	// https://forum.golangbridge.org/t/starting-new-processes-with-exec-command/24956
	err := cmd.Run()
	if err != nil {
		v.parent.logger.LogError(err)
	}

	// Mark the selected video with watched color
	v.setRowColor(v.treeView, constColorWatched)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		// Log that the video has been deleted in the database
		err = v.parent.db.Log.Insert(constLogPlay, video.Title)
		if err != nil {
			v.parent.logger.Log("Failed to log video as watched!")
			v.parent.logger.LogError(err)
		}
		wg.Done()
	}()

	go func() {
		// Log that the video has been deleted in the GUI
		v.parent.activityLog.InsertLog(constLogPlay, video.Title)
		wg.Done()
	}()

	go func() {
		// Set video status as watched
		err = v.parent.db.Videos.UpdateStatus(video.ID, constStatusWatched)
		if err != nil {
			v.parent.logger.Log("Failed to set video status to watched!")
			v.parent.logger.LogError(err)
		}
		wg.Done()
	}()
	wg.Wait()

	v.Refresh("")
}

func (v *videoList) getVideoPath(videoID string) string {
	tryPath := path.Join(v.parent.config.ClientPaths.Videos, videoID+".mkv")
	if _, err := os.Stat(tryPath); err == nil {
		return tryPath
	}

	tryPath = path.Join(v.parent.config.ClientPaths.Videos, videoID+".mp4")
	if _, err := os.Stat(tryPath); err == nil {
		return tryPath
	}

	tryPath = path.Join(v.parent.config.ClientPaths.Videos, videoID+".webm")
	if _, err := os.Stat(tryPath); err == nil {
		return tryPath
	}

	return ""
}

func (v *videoList) getVideoPathForDeletion(videoID string) string {
	tryPath := path.Join(v.parent.config.ClientPaths.Videos, videoID+"*")
	return tryPath
}

// Download a youtube video
func (v *videoList) downloadVideo(video *database.Video, markAsDownloading bool) error {
	// Set the video to be downloaded
	err := v.parent.db.Download.Insert(video.ID)
	if err != nil {
		v.parent.logger.Log("Failed to set video to be downloaded!")
		v.parent.logger.LogError(err)
		return err
	}

	if markAsDownloading {
		// Mark the selected video with downloading color
		v.setRowColor(v.treeView, constColorDownloading)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		// Log that the video has been requested to be downloaded in the database
		err = v.parent.db.Log.Insert(constLogDownload, video.Title)
		if err != nil {
			v.parent.logger.Log("Failed to log video as watched!")
			v.parent.logger.LogError(err)
		}
		wg.Done()
	}()

	go func() {
		// Log that the video has been deleted in the GUI
		v.parent.activityLog.InsertLog(constLogDownload, video.Title)
		wg.Done()
	}()

	go func() {
		// Set video status as downloading
		err = v.parent.db.Videos.UpdateStatus(video.ID, constStatusDownloading)
		if err != nil {
			v.parent.logger.Log("Failed to set video status to downloading!")
			v.parent.logger.LogError(err)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}

func (v *videoList) getSelectedVideo(treeView *gtk.TreeView) *database.Video {
	selection, err := treeView.GetSelection()
	if err != nil {
		return nil
	}
	model, iter, ok := selection.GetSelected()
	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, listStoreColumnVideoID)
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

func (v *videoList) setVideoAsWatched(video *database.Video, mode int) {
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
	err := v.parent.db.Videos.UpdateStatus(video.ID, status)
	if err != nil {
		v.parent.logger.LogFormat("Failed to set video as downloaded/watched/unwatched! %s", video.ID)
		v.parent.logger.LogError(err)
	}
	// v.setRowColor(v.Treeview, constColorSaved)
	v.Refresh("")
}

func (v *videoList) setVideoAsSaved(video *database.Video, saved bool) {
	err := v.parent.db.Videos.UpdateSave(video.ID, saved)
	if err != nil {
		if saved {
			v.parent.logger.LogFormat("Failed to set video as saved! %s", video.ID)
		} else {
			v.parent.logger.LogFormat("Failed to set video as unsaved! %s", video.ID)
		}
		v.parent.logger.LogError(err)
	}
	// v.setRowColor(v.Treeview, constColorSaved)
	v.Refresh("")
}

func (v *videoList) getVideoThumbnailPath(videoID string) string {
	thumbnailPath := "/" + path.Join(v.parent.config.ClientPaths.Thumbnails, fmt.Sprintf("%s.jpg", videoID))
	if _, err := os.Stat(thumbnailPath); err == nil {
		return thumbnailPath
	}
	thumbnailPath = "/" + path.Join(v.parent.config.ClientPaths.Thumbnails, fmt.Sprintf("%s.webp", videoID))
	if _, err := os.Stat(thumbnailPath); err == nil {
		// YouTube started to return *.webp thumbnails instead of *.jpg thumbnails sometimes
		// Go can't read them, and getThumbnail fails to get a PixBuf, so return "" for now
		return ""
		//return thumbnailPath
	}
	return ""
}

func (v *videoList) getVideoThumbnail(videoID string) *gdk.Pixbuf {
	thumbnailPath := v.getVideoThumbnailPath(videoID)
	if thumbnailPath == "" {
		return nil
	}

	thumbnail, err := gdk.PixbufNewFromFile(thumbnailPath)
	if err != nil {
		v.renameJPG2WEBP(thumbnailPath)
		v.parent.logger.LogError(err)
		thumbnail = nil
	} else {
		thumbnail, err = thumbnail.ScaleSimple(160, 90, gdk.INTERP_BILINEAR)
		if err != nil {
			msg := fmt.Sprintf("Failed to scale thumnail (%s)!", thumbnailPath)
			v.parent.logger.Log(msg)
			thumbnail = nil
		}
	}

	return thumbnail
}

// Download a youtube video
func (v *videoList) downloadVideoDuration(video *database.Video) {
	if video == nil {
		return
	}

	go func() {
		command := fmt.Sprintf(constVideoDurationCommand, v.getYoutubePath(), video.ID)
		cmd := exec.Command("/bin/bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return
		}
		duration := string(output)
		if duration == "0" || strings.HasPrefix(duration, "ERROR: Premieres") || strings.HasPrefix(duration, "ERROR: This live event") {
			// Is it a live streaming event?
			duration = "LIVE"
		}

		_ = v.parent.db.Videos.UpdateDuration(video.ID, duration)
	}()
}

// Get the thumbnail of a youtube video
func (v *videoList) downloadVideoThumbnail(video *database.Video) (string, error) {
	// %s/%s.jpg
	thumbPath := fmt.Sprintf(constThumbnailLocation, v.parent.config.ServerPaths.Thumbnails, video.ID)

	// Don't download thumbnail if it already exists
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		command := fmt.Sprintf(constThumbnailCommand, v.getYoutubePath(), thumbPath, video.ID)
		cmd := exec.Command("/bin/bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return string(output), err
		}
	}
	return "", nil
}
