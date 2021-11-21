package softtube

import (
	"strings"

	"github.com/gotk3/gotk3/gtk"

	database "github.com/hultan/softtube/internal/softtube.database"
)

type color struct {
	videoList *videoList
}

func (c *color) getColor(video *database.Video) (colorType, colorType) {
	if video.Saved {
		return constColorSaved, "Black"
	}
	duration := c.videoList.removeInvalidDurations(video.Duration)
	if strings.Trim(duration, " ") == "LIVE" {
		// If duration is LIVE, lets change color to live color
		return constColorLive, "Black"
	} else if duration == "" {
		// If duration is invalid, lets change color to warning
		return constColorWarning, "Black"
	} else if video.Status == constStatusDeleted {
		return constColorDeleted, "Black"
	} else if video.Status == constStatusWatched {
		return constColorWatched, "Black"
	} else if video.Status == constStatusDownloaded {
		return constColorDownloaded, "Black"
	} else if video.Status == constStatusDownloading {
		return constColorDownloading, "Black"
	} else {
		return constColorNotDownloaded, "White"
	}
}

func (c *color) setRowColor(treeView *gtk.TreeView, color string) {
	selection, _ := treeView.GetSelection()
	rows := selection.GetSelectedRows(listStore)
	if rows == nil {
		return
	}
	treePath := rows.Data().(*gtk.TreePath)
	iter, _ := listStore.GetIter(treePath)
	_ = listStore.SetValue(iter, int(listStoreColumnBackground), color)
}
