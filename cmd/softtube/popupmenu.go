package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softtube/internal/softtube.core"
)

// PopupMenu : Handler the video list popupmenu
type PopupMenu struct {
	Parent                *SoftTube
	PopupMenu             *gtk.Menu
	PopupRefresh          *gtk.MenuItem
	PopupDownload         *gtk.MenuItem
	PopupRedownload       *gtk.MenuItem
	PopupRedownloadVideo  *gtk.MenuItem
	PopupRedownloadVideos *gtk.MenuItem
	PopupPlay             *gtk.MenuItem
	PopupGetDuration      *gtk.MenuItem
	PopupGetVideoID       *gtk.MenuItem
	PopupGetThumbnail     *gtk.MenuItem
	PopupDeleteAll        *gtk.MenuItem
	PopupUnwatch          *gtk.MenuItem
	PopupSave             *gtk.MenuItem
	PopupSubscriptions    *gtk.MenuItem
	PopupToWatch          *gtk.MenuItem
	PopupToDelete         *gtk.MenuItem
	PopupSaved            *gtk.MenuItem
}

// Load : Loads the popup menu
func (p *PopupMenu) Load(helper *core.GtkHelper) error {
	menu, err := helper.GetMenu("popupmenu")
	if err != nil {
		return err
	}
	p.PopupMenu = menu

	menuItem, err := helper.GetMenuItem("popup_refresh")
	if err != nil {
		return err
	}
	p.PopupRefresh = menuItem

	menuItem, err = helper.GetMenuItem("popup_download")
	if err != nil {
		return err
	}
	p.PopupDownload = menuItem

	menuSubMenu, err := helper.GetMenuItem("popup_redownload")
	if err != nil {
		return err
	}
	p.PopupRedownload = menuSubMenu

	menuItem, err = helper.GetMenuItem("popup_redownload_failedvideo")
	if err != nil {
		return err
	}
	p.PopupRedownloadVideo = menuItem

	menuItem, err = helper.GetMenuItem("popup_redownload_failedvideos")
	if err != nil {
		return err
	}
	p.PopupRedownloadVideos = menuItem

	menuItem, err = helper.GetMenuItem("popup_play")
	if err != nil {
		return err
	}
	p.PopupPlay = menuItem

	menuItem, err = helper.GetMenuItem("popup_get_duration")
	if err != nil {
		return err
	}
	p.PopupGetDuration = menuItem

	menuItem, err = helper.GetMenuItem("popup_get_videoid")
	if err != nil {
		return err
	}
	p.PopupGetVideoID = menuItem

	menuItem, err = helper.GetMenuItem("popup_get_thumbnail")
	if err != nil {
		return err
	}
	p.PopupGetThumbnail = menuItem

	menuItem, err = helper.GetMenuItem("popup_delete_all")
	if err != nil {
		return err
	}
	p.PopupDeleteAll = menuItem

	menuItem, err = helper.GetMenuItem("popup_unwatch")
	if err != nil {
		return err
	}
	p.PopupUnwatch = menuItem

	menuItem, err = helper.GetMenuItem("popup_save")
	if err != nil {
		return err
	}
	p.PopupSave = menuItem

	menuItem, err = helper.GetMenuItem("popup_view_subscriptions")
	if err != nil {
		return err
	}
	p.PopupSubscriptions = menuItem

	menuItem, err = helper.GetMenuItem("popup_view_to_watch")
	if err != nil {
		return err
	}
	p.PopupToWatch = menuItem

	menuItem, err = helper.GetMenuItem("popup_view_to_delete")
	if err != nil {
		return err
	}
	p.PopupToDelete = menuItem

	menuItem, err = helper.GetMenuItem("popup_view_saved")
	if err != nil {
		return err
	}
	p.PopupSaved = menuItem

	return nil
}

// SetupEvents : Setup the toolbar events
func (p *PopupMenu) SetupEvents() {
	_, _ = p.Parent.VideoList.TreeView.Connect("button-release-event", func(treeview *gtk.TreeView, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Button() == gdk.BUTTON_SECONDARY {
			videoSelected := p.Parent.VideoList.getSelectedVideo(p.Parent.VideoList.TreeView) != nil
			switch p.Parent.VideoList.FilterMode {
			case constFilterModeSubscriptions:
				p.PopupDownload.SetSensitive(videoSelected)
				p.PopupRedownload.SetSensitive(true)
				p.PopupRedownloadVideo.SetSensitive(videoSelected)
				p.PopupPlay.SetSensitive(false)
				p.PopupGetDuration.SetSensitive(videoSelected)
				p.PopupGetVideoID.SetSensitive(videoSelected)
				p.PopupGetThumbnail.SetSensitive(videoSelected)
				p.PopupDeleteAll.SetVisible(false)
				p.PopupUnwatch.SetSensitive(videoSelected)
				p.PopupUnwatch.SetLabel(constSetAsNotDownloaded)
				p.PopupSave.SetSensitive(false)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupSubscriptions.SetSensitive(false)
				p.PopupToWatch.SetSensitive(true)
				p.PopupToDelete.SetSensitive(true)
				p.PopupSaved.SetSensitive(true)
			case constFilterModeToWatch:
				p.PopupDownload.SetSensitive(false)
				p.PopupRedownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(videoSelected)
				p.PopupGetDuration.SetSensitive(false)
				p.PopupGetVideoID.SetSensitive(videoSelected)
				p.PopupGetThumbnail.SetSensitive(false)
				p.PopupDeleteAll.SetVisible(false)
				p.PopupUnwatch.SetSensitive(videoSelected)
				p.PopupUnwatch.SetLabel(constSetAsWatched)
				p.PopupSave.SetSensitive(videoSelected)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupSubscriptions.SetSensitive(true)
				p.PopupToWatch.SetSensitive(false)
				p.PopupToDelete.SetSensitive(true)
				p.PopupSaved.SetSensitive(true)
			case constFilterModeToDelete:
				p.PopupDownload.SetSensitive(false)
				p.PopupRedownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(videoSelected)
				p.PopupGetDuration.SetSensitive(false)
				p.PopupGetVideoID.SetSensitive(videoSelected)
				p.PopupGetThumbnail.SetSensitive(false)
				p.PopupDeleteAll.SetVisible(true)
				p.PopupUnwatch.SetSensitive(videoSelected)
				p.PopupUnwatch.SetLabel(constSetAsUnwatched)
				p.PopupSave.SetSensitive(videoSelected)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupSubscriptions.SetSensitive(true)
				p.PopupToWatch.SetSensitive(true)
				p.PopupToDelete.SetSensitive(false)
				p.PopupSaved.SetSensitive(true)
			case constFilterModeSaved:
				p.PopupDownload.SetSensitive(false)
				p.PopupRedownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(videoSelected)
				p.PopupGetDuration.SetSensitive(false)
				p.PopupGetVideoID.SetSensitive(videoSelected)
				p.PopupGetThumbnail.SetSensitive(false)
				p.PopupDeleteAll.SetVisible(false)
				p.PopupUnwatch.SetSensitive(false)
				p.PopupUnwatch.SetLabel(constSetAsWatched)
				p.PopupSave.SetSensitive(false)
				p.PopupSave.SetLabel(constSetAsNotSaved)
				p.PopupSubscriptions.SetSensitive(true)
				p.PopupToWatch.SetSensitive(true)
				p.PopupToDelete.SetSensitive(true)
				p.PopupSaved.SetSensitive(false)
			}
			p.PopupMenu.PopupAtPointer(event)
		}
	})

	_, _ = p.PopupRefresh.Connect("activate", func() {
		p.Parent.VideoList.Refresh("")
	})

	_, _ = p.PopupDownload.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video != nil {
			_ = p.Parent.VideoList.downloadVideo(video)
		}
	})
	_, _ = p.PopupRedownloadVideo.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video != nil {
			_ = p.Parent.VideoList.downloadVideo(video)
		}
	})
	_, _ = p.PopupRedownloadVideos.Connect("activate", func() {
		videos, err := db.Videos.GetVideos(true)
		if err != nil {
			logger.LogError(err)
			return
		}
		for key,_ := range videos {
			video := &videos[key]
			if video != nil {
				_ = p.Parent.VideoList.downloadVideo(video)
			}
		}
	})

	_, _ = p.PopupPlay.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video != nil {
			p.Parent.VideoList.playVideo(video)
		}
	})

	_, _ = p.PopupGetDuration.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video != nil {
			p.Parent.VideoList.downloadDuration(video)
		}
	})

	_, _ = p.PopupGetVideoID.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video != nil {
			clipboard, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
			if err != nil {
				fmt.Println("Clipboard error!")
				return
			}
			clipboard.SetText(video.ID)
		}
	})

	_, _ = p.PopupGetThumbnail.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video != nil {
			p.Parent.VideoList.downloadThumbnail(video)
		}
	})

	_, _ = p.PopupDeleteAll.Connect("activate", func() {
		p.Parent.VideoList.DeleteWatchedVideos()
	})

	_, _ = p.PopupSave.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video != nil {
			mode := p.PopupSave.GetLabel() == constSetAsSaved
			p.Parent.VideoList.setAsSaved(video, mode)
		}
	})

	_, _ = p.PopupUnwatch.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video != nil {
			var mode int
			switch p.PopupUnwatch.GetLabel() {
			case constSetAsNotDownloaded:
				mode = 0
				break
			case constSetAsWatched:
				mode = 1
				break
			case constSetAsUnwatched:
				mode = 2
				break
			}
			p.Parent.VideoList.setAsWatched(video, mode)
		}
	})
	_, _ = p.PopupSubscriptions.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeSubscriptions)
	})
	_, _ = p.PopupToWatch.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeToWatch)
	})
	_, _ = p.PopupToDelete.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeToDelete)
	})
	_, _ = p.PopupSaved.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeSaved)
	})
}
