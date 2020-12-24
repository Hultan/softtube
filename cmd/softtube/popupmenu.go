package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// PopupMenu : Handler the video list popupmenu
type PopupMenu struct {
	Parent             *SoftTube
	PopupMenu          *gtk.Menu
	PopupRefresh       *gtk.MenuItem
	PopupDownload      *gtk.MenuItem
	PopupPlay          *gtk.MenuItem
	PopupGetDuration   *gtk.MenuItem
	PopupGetVideoID    *gtk.MenuItem
	PopupDeleteAll     *gtk.MenuItem
	PopupUnwatch       *gtk.MenuItem
	PopupSave          *gtk.MenuItem
	PopupSubscriptions *gtk.MenuItem
	PopupToWatch       *gtk.MenuItem
	PopupToDelete      *gtk.MenuItem
	PopupSaved         *gtk.MenuItem
}

// Load : Loads the popup menu
func (p *PopupMenu) Load(helper *GtkHelper) error {
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
	_,_ = p.Parent.VideoList.TreeView.Connect("button-release-event", func(treeview *gtk.TreeView, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Button() == 3 { // 3 == Mouse right button!?
			switch p.Parent.VideoList.FilterMode {
			case constFilterModeSubscriptions:
				p.PopupDownload.SetSensitive(true)
				p.PopupPlay.SetSensitive(false)
				p.PopupGetDuration.SetSensitive(true)
				p.PopupGetVideoID.SetSensitive(true)
				p.PopupDeleteAll.SetSensitive(false)
				p.PopupUnwatch.SetSensitive(true)
				p.PopupUnwatch.SetLabel(constSetAsNotDownloaded)
				p.PopupSave.SetSensitive(false)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupSubscriptions.SetSensitive(false)
				p.PopupToWatch.SetSensitive(true)
				p.PopupToDelete.SetSensitive(true)
				p.PopupSaved.SetSensitive(true)
			case constFilterModeToWatch:
				p.PopupDownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(true)
				p.PopupGetDuration.SetSensitive(false)
				p.PopupGetVideoID.SetSensitive(true)
				p.PopupDeleteAll.SetSensitive(false)
				p.PopupUnwatch.SetSensitive(true)
				p.PopupUnwatch.SetLabel(constSetAsWatched)
				p.PopupSave.SetSensitive(true)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupSubscriptions.SetSensitive(true)
				p.PopupToWatch.SetSensitive(false)
				p.PopupToDelete.SetSensitive(true)
				p.PopupSaved.SetSensitive(true)
			case constFilterModeToDelete:
				p.PopupDownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(true)
				p.PopupGetDuration.SetSensitive(false)
				p.PopupGetVideoID.SetSensitive(true)
				p.PopupDeleteAll.SetSensitive(true)
				p.PopupUnwatch.SetSensitive(true)
				p.PopupUnwatch.SetLabel(constSetAsUnwatched)
				p.PopupSave.SetSensitive(true)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupSubscriptions.SetSensitive(true)
				p.PopupToWatch.SetSensitive(true)
				p.PopupToDelete.SetSensitive(false)
				p.PopupSaved.SetSensitive(true)
			case constFilterModeSaved:
				p.PopupDownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(true)
				p.PopupGetDuration.SetSensitive(false)
				p.PopupGetVideoID.SetSensitive(true)
				p.PopupDeleteAll.SetSensitive(false)
				p.PopupUnwatch.SetSensitive(false)
				p.PopupUnwatch.SetLabel(constSetAsWatched)
				p.PopupSave.SetSensitive(true)
				p.PopupSave.SetLabel(constSetAsNotSaved)
				p.PopupSubscriptions.SetSensitive(true)
				p.PopupToWatch.SetSensitive(true)
				p.PopupToDelete.SetSensitive(true)
				p.PopupSaved.SetSensitive(false)
			}
			p.PopupMenu.PopupAtPointer(event)
		}
	})

	_,_ = p.PopupRefresh.Connect("activate", func() {
		p.Parent.VideoList.Refresh("")
	})

	_,_ = p.PopupDownload.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video!=nil {
			_ = p.Parent.VideoList.downloadVideo(video)
		}
	})

	_,_ = p.PopupPlay.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video!=nil {
			p.Parent.VideoList.playVideo(video)
		}
	})

	_,_ = p.PopupGetDuration.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video!=nil {
			p.Parent.VideoList.downloadDuration(video)
		}
	})

	_,_ = p.PopupGetVideoID.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video!=nil {
			clipboard, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
			if err != nil {
				fmt.Println("Clipboard error!")
			}
			clipboard.SetText(video.ID)
		}
	})

	_,_ = p.PopupDeleteAll.Connect("activate", func() {
		p.Parent.VideoList.DeleteWatchedVideos()
	})

	_,_ = p.PopupSave.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video!=nil {
			mode := p.PopupSave.GetLabel() == constSetAsSaved
			p.Parent.VideoList.setAsSaved(video, mode)
		}
	})

	_,_ = p.PopupUnwatch.Connect("activate", func() {
		treeview := p.Parent.VideoList.TreeView
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		if video!=nil {
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
	_,_ = p.PopupSubscriptions.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeSubscriptions)
	})
	_,_ = p.PopupToWatch.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeToWatch)
	})
	_,_ = p.PopupToDelete.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeToDelete)
	})
	_,_ = p.PopupSaved.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeSaved)
	})
}
