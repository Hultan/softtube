package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// PopupMenu : Handler the video list popupmenu
type PopupMenu struct {
	Parent                *SoftTube
	PopupMenu             *gtk.Menu
	PopupRefresh          *gtk.MenuItem
	PopupDownload         *gtk.MenuItem
	PopupRedownload        *gtk.MenuItem
	PopupRedownloadVideo   *gtk.MenuItem
	PopupRedownloadVideos  *gtk.MenuItem
	PopupPlay              *gtk.MenuItem
	PopupGetDuration       *gtk.MenuItem
	PopupGetVideoID        *gtk.MenuItem
	PopupGetThumbnail      *gtk.MenuItem
	PopupDeleteAll         *gtk.MenuItem
	PopupUnwatch           *gtk.MenuItem
	PopupSave              *gtk.MenuItem
	PopupViewSubscriptions *gtk.MenuItem
	PopupViewDownloads     *gtk.MenuItem
	PopupViewToWatch       *gtk.MenuItem
	PopupViewToDelete      *gtk.MenuItem
	PopupViewSaved         *gtk.MenuItem
}

// Load : Loads the popup menu
func (p *PopupMenu) Load(builder *SoftBuilder) error {
	menu := builder.getObject("popupmenu").(*gtk.Menu)
	p.PopupMenu = menu

	menuItem := builder.getObject("popup_refresh").(*gtk.MenuItem)
	p.PopupRefresh = menuItem

	menuItem = builder.getObject("popup_download").(*gtk.MenuItem)
	p.PopupDownload = menuItem

	menuItem = builder.getObject("popup_redownload").(*gtk.MenuItem)
	p.PopupRedownload = menuItem

	menuItem = builder.getObject("popup_redownload_failedvideo").(*gtk.MenuItem)
	p.PopupRedownloadVideo = menuItem

	menuItem = builder.getObject("popup_redownload_failedvideos").(*gtk.MenuItem)
	p.PopupRedownloadVideos = menuItem

	menuItem = builder.getObject("popup_play").(*gtk.MenuItem)
	p.PopupPlay = menuItem

	menuItem = builder.getObject("popup_get_duration").(*gtk.MenuItem)
	p.PopupGetDuration = menuItem

	menuItem = builder.getObject("popup_get_videoid").(*gtk.MenuItem)
	p.PopupGetVideoID = menuItem

	menuItem = builder.getObject("popup_get_thumbnail").(*gtk.MenuItem)
	p.PopupGetThumbnail = menuItem

	menuItem = builder.getObject("popup_delete_all").(*gtk.MenuItem)
	p.PopupDeleteAll = menuItem

	menuItem = builder.getObject("popup_unwatch").(*gtk.MenuItem)
	p.PopupUnwatch = menuItem

	menuItem = builder.getObject("popup_save").(*gtk.MenuItem)
	p.PopupSave = menuItem

	menuItem = builder.getObject("popup_view_subscriptions").(*gtk.MenuItem)
	p.PopupViewSubscriptions = menuItem

	menuItem = builder.getObject("popup_view_failed").(*gtk.MenuItem)
	p.PopupViewDownloads = menuItem

	menuItem = builder.getObject("popup_view_to_watch").(*gtk.MenuItem)
	p.PopupViewToWatch = menuItem

	menuItem = builder.getObject("popup_view_to_delete").(*gtk.MenuItem)
	p.PopupViewToDelete = menuItem

	menuItem = builder.getObject("popup_view_saved").(*gtk.MenuItem)
	p.PopupViewSaved = menuItem

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
				p.PopupViewSubscriptions.SetSensitive(false)
				p.PopupViewDownloads.SetSensitive(true)
				p.PopupViewToWatch.SetSensitive(true)
				p.PopupViewToDelete.SetSensitive(true)
				p.PopupViewSaved.SetSensitive(true)
			case constFilterModeDownloads:
				p.PopupDownload.SetSensitive(false)
				p.PopupRedownload.SetSensitive(true)
				p.PopupRedownloadVideo.SetSensitive(videoSelected)
				p.PopupPlay.SetSensitive(false)
				p.PopupGetDuration.SetSensitive(videoSelected)
				p.PopupGetVideoID.SetSensitive(videoSelected)
				p.PopupGetThumbnail.SetSensitive(videoSelected)
				p.PopupDeleteAll.SetVisible(false)
				p.PopupUnwatch.SetSensitive(true)
				p.PopupUnwatch.SetLabel(constSetAsNotDownloaded)
				p.PopupSave.SetSensitive(false)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupViewSubscriptions.SetSensitive(true)
				p.PopupViewDownloads.SetSensitive(false)
				p.PopupViewToWatch.SetSensitive(true)
				p.PopupViewToDelete.SetSensitive(true)
				p.PopupViewSaved.SetSensitive(true)
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
				p.PopupViewSubscriptions.SetSensitive(true)
				p.PopupViewDownloads.SetSensitive(true)
				p.PopupViewToWatch.SetSensitive(false)
				p.PopupViewToDelete.SetSensitive(true)
				p.PopupViewSaved.SetSensitive(true)
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
				p.PopupViewSubscriptions.SetSensitive(true)
				p.PopupViewDownloads.SetSensitive(true)
				p.PopupViewToWatch.SetSensitive(true)
				p.PopupViewToDelete.SetSensitive(false)
				p.PopupViewSaved.SetSensitive(true)
			case constFilterModeSaved:
				p.PopupDownload.SetSensitive(false)
				p.PopupRedownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(videoSelected)
				p.PopupGetDuration.SetSensitive(false)
				p.PopupGetVideoID.SetSensitive(videoSelected)
				p.PopupGetThumbnail.SetSensitive(true)
				p.PopupDeleteAll.SetVisible(false)
				p.PopupUnwatch.SetSensitive(false)
				p.PopupUnwatch.SetLabel(constSetAsWatched)
				p.PopupSave.SetSensitive(true)
				p.PopupSave.SetLabel(constSetAsNotSaved)
				p.PopupViewSubscriptions.SetSensitive(true)
				p.PopupViewDownloads.SetSensitive(true)
				p.PopupViewToWatch.SetSensitive(true)
				p.PopupViewToDelete.SetSensitive(true)
				p.PopupViewSaved.SetSensitive(false)
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
		for key, _ := range videos {
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
			p.Parent.VideoList.downloadVideoDuration(video)
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
			p.Parent.VideoList.downloadVideoThumbnail(video)
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
			p.Parent.VideoList.setVideoAsSaved(video, mode)
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
			p.Parent.VideoList.setVideoAsWatched(video, mode)
		}
	})
	_, _ = p.PopupViewSubscriptions.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeSubscriptions)
	})
	_, _ = p.PopupViewDownloads.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeDownloads)
	})
	_, _ = p.PopupViewToWatch.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeToWatch)
	})
	_, _ = p.PopupViewToDelete.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeToDelete)
	})
	_, _ = p.PopupViewSaved.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeSaved)
	})
}
