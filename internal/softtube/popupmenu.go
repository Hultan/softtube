package softtube

import (
	"fmt"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
)

// popupMenu : Handler the video list popupmenu
type popupMenu struct {
	parent                 *SoftTube
	popupMenu              *gtk.Menu
	popupRefresh           *gtk.MenuItem
	popupDownload          *gtk.MenuItem
	popupRedownload        *gtk.MenuItem
	popupRedownloadVideo   *gtk.MenuItem
	popupRedownloadVideos  *gtk.MenuItem
	popupPlay              *gtk.MenuItem
	popupGetDuration       *gtk.MenuItem
	popupGetVideoID        *gtk.MenuItem
	popupGetThumbnail      *gtk.MenuItem
	popupDeleteAll         *gtk.MenuItem
	popupUnwatch           *gtk.MenuItem
	popupSave              *gtk.MenuItem
	popupViewSubscriptions *gtk.MenuItem
	popupViewDownloads     *gtk.MenuItem
	popupViewToWatch       *gtk.MenuItem
	popupViewToDelete      *gtk.MenuItem
	popupViewSaved         *gtk.MenuItem
}

// Init : Loads the popup menu
func (p *popupMenu) Init(builder *framework.GtkBuilder) error {
	p.popupMenu = builder.GetObject("popupmenu").(*gtk.Menu)
	p.popupRefresh = builder.GetObject("popup_refresh").(*gtk.MenuItem)
	p.popupDownload = builder.GetObject("popup_download").(*gtk.MenuItem)
	p.popupRedownload = builder.GetObject("popup_redownload").(*gtk.MenuItem)
	p.popupRedownloadVideo = builder.GetObject("popup_redownload_failedvideo").(*gtk.MenuItem)
	p.popupRedownloadVideos = builder.GetObject("popup_redownload_failedvideos").(*gtk.MenuItem)
	p.popupPlay = builder.GetObject("popup_play").(*gtk.MenuItem)
	p.popupGetDuration = builder.GetObject("popup_get_duration").(*gtk.MenuItem)
	p.popupGetVideoID = builder.GetObject("popup_get_videoid").(*gtk.MenuItem)
	p.popupGetThumbnail = builder.GetObject("popup_get_thumbnail").(*gtk.MenuItem)
	p.popupDeleteAll = builder.GetObject("popup_delete_all").(*gtk.MenuItem)
	p.popupUnwatch = builder.GetObject("popup_unwatch").(*gtk.MenuItem)
	p.popupSave = builder.GetObject("popup_save").(*gtk.MenuItem)
	p.popupViewSubscriptions = builder.GetObject("popup_view_subscriptions").(*gtk.MenuItem)
	p.popupViewDownloads = builder.GetObject("popup_view_failed").(*gtk.MenuItem)
	p.popupViewToWatch = builder.GetObject("popup_view_to_watch").(*gtk.MenuItem)
	p.popupViewToDelete = builder.GetObject("popup_view_to_delete").(*gtk.MenuItem)
	p.popupViewSaved = builder.GetObject("popup_view_saved").(*gtk.MenuItem)

	p.SetupEvents()

	return nil
}

// SetupEvents : Set up the toolbar events
func (p *popupMenu) SetupEvents() {
	_ = p.parent.videoList.treeView.Connect("button-release-event", func(treeview *gtk.TreeView, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Button() == gdk.BUTTON_SECONDARY {
			videoSelected := p.parent.videoList.videoFunctions.getSelected(p.parent.videoList.treeView) != nil
			switch p.parent.videoList.currentView {
			case viewSubscriptions:
				p.popupDownload.SetSensitive(videoSelected)
				p.popupRedownload.SetSensitive(true)
				p.popupRedownloadVideo.SetSensitive(videoSelected)
				p.popupPlay.SetSensitive(false)
				p.popupGetDuration.SetSensitive(videoSelected)
				p.popupGetVideoID.SetSensitive(videoSelected)
				p.popupGetThumbnail.SetSensitive(videoSelected)
				p.popupDeleteAll.SetVisible(false)
				p.popupUnwatch.SetSensitive(videoSelected)
				p.popupUnwatch.SetLabel(constSetAsNotDownloaded)
				p.popupSave.SetSensitive(false)
				p.popupSave.SetLabel(constSetAsSaved)
				p.popupViewSubscriptions.SetSensitive(false)
				p.popupViewDownloads.SetSensitive(true)
				p.popupViewToWatch.SetSensitive(true)
				p.popupViewToDelete.SetSensitive(true)
				p.popupViewSaved.SetSensitive(true)
			case viewDownloads:
				p.popupDownload.SetSensitive(false)
				p.popupRedownload.SetSensitive(true)
				p.popupRedownloadVideo.SetSensitive(videoSelected)
				p.popupPlay.SetSensitive(false)
				p.popupGetDuration.SetSensitive(videoSelected)
				p.popupGetVideoID.SetSensitive(videoSelected)
				p.popupGetThumbnail.SetSensitive(videoSelected)
				p.popupDeleteAll.SetVisible(false)
				p.popupUnwatch.SetSensitive(true)
				p.popupUnwatch.SetLabel(constSetAsNotDownloaded)
				p.popupSave.SetSensitive(false)
				p.popupSave.SetLabel(constSetAsSaved)
				p.popupViewSubscriptions.SetSensitive(true)
				p.popupViewDownloads.SetSensitive(false)
				p.popupViewToWatch.SetSensitive(true)
				p.popupViewToDelete.SetSensitive(true)
				p.popupViewSaved.SetSensitive(true)
			case viewToWatch:
				p.popupDownload.SetSensitive(false)
				p.popupRedownload.SetSensitive(false)
				p.popupPlay.SetSensitive(videoSelected)
				p.popupGetDuration.SetSensitive(false)
				p.popupGetVideoID.SetSensitive(videoSelected)
				p.popupGetThumbnail.SetSensitive(false)
				p.popupDeleteAll.SetVisible(false)
				p.popupUnwatch.SetSensitive(videoSelected)
				p.popupUnwatch.SetLabel(constSetAsWatched)
				p.popupSave.SetSensitive(videoSelected)
				p.popupSave.SetLabel(constSetAsSaved)
				p.popupViewSubscriptions.SetSensitive(true)
				p.popupViewDownloads.SetSensitive(true)
				p.popupViewToWatch.SetSensitive(false)
				p.popupViewToDelete.SetSensitive(true)
				p.popupViewSaved.SetSensitive(true)
			case viewToDelete:
				p.popupDownload.SetSensitive(false)
				p.popupRedownload.SetSensitive(false)
				p.popupPlay.SetSensitive(videoSelected)
				p.popupGetDuration.SetSensitive(false)
				p.popupGetVideoID.SetSensitive(videoSelected)
				p.popupGetThumbnail.SetSensitive(false)
				p.popupDeleteAll.SetVisible(true)
				p.popupUnwatch.SetSensitive(videoSelected)
				p.popupUnwatch.SetLabel(constSetAsUnwatched)
				p.popupSave.SetSensitive(videoSelected)
				p.popupSave.SetLabel(constSetAsSaved)
				p.popupViewSubscriptions.SetSensitive(true)
				p.popupViewDownloads.SetSensitive(true)
				p.popupViewToWatch.SetSensitive(true)
				p.popupViewToDelete.SetSensitive(false)
				p.popupViewSaved.SetSensitive(true)
			case viewSaved:
				p.popupDownload.SetSensitive(false)
				p.popupRedownload.SetSensitive(false)
				p.popupPlay.SetSensitive(videoSelected)
				p.popupGetDuration.SetSensitive(false)
				p.popupGetVideoID.SetSensitive(videoSelected)
				p.popupGetThumbnail.SetSensitive(true)
				p.popupDeleteAll.SetVisible(false)
				p.popupUnwatch.SetSensitive(false)
				p.popupUnwatch.SetLabel(constSetAsWatched)
				p.popupSave.SetSensitive(true)
				p.popupSave.SetLabel(constSetAsNotSaved)
				p.popupViewSubscriptions.SetSensitive(true)
				p.popupViewDownloads.SetSensitive(true)
				p.popupViewToWatch.SetSensitive(true)
				p.popupViewToDelete.SetSensitive(true)
				p.popupViewSaved.SetSensitive(false)
			}
			p.popupMenu.PopupAtPointer(event)
		}
	})

	_ = p.popupRefresh.Connect("activate", func() {
		p.parent.videoList.Refresh("")
	})

	_ = p.popupDownload.Connect("activate", func() {
		treeview := p.parent.videoList.treeView
		vid := p.parent.videoList.videoFunctions.getSelected(treeview)
		if vid != nil {
			_ = p.parent.videoList.videoFunctions.download(vid, true)
		}
	})
	_ = p.popupRedownloadVideo.Connect("activate", func() {
		treeview := p.parent.videoList.treeView
		vid := p.parent.videoList.videoFunctions.getSelected(treeview)
		if vid != nil {
			_ = p.parent.videoList.videoFunctions.download(vid, false)
		}
	})
	_ = p.popupRedownloadVideos.Connect("activate", func() {
		videos, err := p.parent.DB.Videos.GetVideos(true)
		if err != nil {
			p.parent.Logger.LogError(err)
			return
		}
		for key := range videos {
			vid := &videos[key]
			if vid != nil {
				_ = p.parent.videoList.videoFunctions.download(vid, false)
			}
		}
	})

	_ = p.popupPlay.Connect("activate", func() {
		treeview := p.parent.videoList.treeView
		vid := p.parent.videoList.videoFunctions.getSelected(treeview)
		if vid != nil {
			p.parent.videoList.videoFunctions.play(vid)
		}
	})

	_ = p.popupGetDuration.Connect("activate", func() {
		treeview := p.parent.videoList.treeView
		vid := p.parent.videoList.videoFunctions.getSelected(treeview)
		if vid != nil {
			p.parent.videoList.videoFunctions.downloadDuration(vid)
		}
	})

	_ = p.popupGetVideoID.Connect("activate", func() {
		treeview := p.parent.videoList.treeView
		vid := p.parent.videoList.videoFunctions.getSelected(treeview)
		if vid != nil {
			clipboard, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
			if err != nil {
				fmt.Println("Clipboard error!")
				return
			}
			clipboard.SetText(vid.ID)
		}
	})

	_ = p.popupGetThumbnail.Connect("activate", func() {
		treeview := p.parent.videoList.treeView
		vid := p.parent.videoList.videoFunctions.getSelected(treeview)
		if vid != nil {
			_, _ = p.parent.videoList.videoFunctions.downloadThumbnail(vid)
		}
	})

	_ = p.popupDeleteAll.Connect("activate", func() {
		p.parent.videoList.DeleteWatchedVideos()
	})

	_ = p.popupSave.Connect("activate", func() {
		treeview := p.parent.videoList.treeView
		vid := p.parent.videoList.videoFunctions.getSelected(treeview)
		if vid != nil {
			mode := p.popupSave.GetLabel() == constSetAsSaved
			p.parent.videoList.videoFunctions.setAsSaved(vid, mode)
		}
	})

	_ = p.popupUnwatch.Connect("activate", func() {
		treeview := p.parent.videoList.treeView
		vid := p.parent.videoList.videoFunctions.getSelected(treeview)
		if vid != nil {
			var mode int
			switch p.popupUnwatch.GetLabel() {
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
			p.parent.videoList.videoFunctions.setAsWatched(vid, mode)
		}
	})
	_ = p.popupViewSubscriptions.Connect("activate", func() {
		p.parent.videoList.switchView(viewSubscriptions)
	})
	_ = p.popupViewDownloads.Connect("activate", func() {
		p.parent.videoList.switchView(viewDownloads)
	})
	_ = p.popupViewToWatch.Connect("activate", func() {
		p.parent.videoList.switchView(viewToWatch)
	})
	_ = p.popupViewSaved.Connect("activate", func() {
		p.parent.videoList.switchView(viewSaved)
	})
	_ = p.popupViewToDelete.Connect("activate", func() {
		p.parent.videoList.switchView(viewToDelete)
	})
}
