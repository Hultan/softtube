package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	gtkhelper "github.com/hultan/softteam/gtk"
)

// PopupMenu : Handler the video list popupmenu
type PopupMenu struct {
	Parent             *SoftTube
	PopupMenu          *gtk.Menu
	PopupRefresh       *gtk.MenuItem
	PopupDownload      *gtk.MenuItem
	PopupPlay          *gtk.MenuItem
	PopupDeleteAll     *gtk.MenuItem
	PopupUnwatch       *gtk.MenuItem
	PopupSave          *gtk.MenuItem
	PopupSubscriptions *gtk.MenuItem
	PopupToWatch       *gtk.MenuItem
	PopupToDelete      *gtk.MenuItem
	PopupSaved         *gtk.MenuItem
}

// Load : Loads the popup menu
func (p *PopupMenu) Load(builder *gtk.Builder) error {
	helper := new(gtkhelper.GtkHelper)

	menu, err := helper.GetMenu(builder, "popupmenu")
	if err != nil {
		return err
	}
	p.PopupMenu = menu

	menuItem, err := helper.GetMenuItem(builder, "popup_refresh")
	if err != nil {
		return err
	}
	p.PopupRefresh = menuItem

	menuItem, err = helper.GetMenuItem(builder, "popup_download")
	if err != nil {
		return err
	}
	p.PopupDownload = menuItem

	menuItem, err = helper.GetMenuItem(builder, "popup_play")
	if err != nil {
		return err
	}
	p.PopupPlay = menuItem

	menuItem, err = helper.GetMenuItem(builder, "popup_delete_all")
	if err != nil {
		return err
	}
	p.PopupDeleteAll = menuItem

	menuItem, err = helper.GetMenuItem(builder, "popup_unwatch")
	if err != nil {
		return err
	}
	p.PopupUnwatch = menuItem

	menuItem, err = helper.GetMenuItem(builder, "popup_save")
	if err != nil {
		return err
	}
	p.PopupSave = menuItem

	menuItem, err = helper.GetMenuItem(builder, "popup_view_subscriptions")
	if err != nil {
		return err
	}
	p.PopupSubscriptions = menuItem

	menuItem, err = helper.GetMenuItem(builder, "popup_view_to_watch")
	if err != nil {
		return err
	}
	p.PopupToWatch = menuItem

	menuItem, err = helper.GetMenuItem(builder, "popup_view_to_delete")
	if err != nil {
		return err
	}
	p.PopupToDelete = menuItem

	menuItem, err = helper.GetMenuItem(builder, "popup_view_saved")
	if err != nil {
		return err
	}
	p.PopupSaved = menuItem

	return nil
}

// SetupEvents : Setup the toolbar events
func (p *PopupMenu) SetupEvents() {
	p.Parent.VideoList.Treeview.Connect("button-release-event", func(treeview *gtk.TreeView, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Button() == 3 { // 3 == Mouse right button!?
			switch p.Parent.VideoList.FilterMode {
			case constFilterModeSubscriptions:
				p.PopupDownload.SetSensitive(true)
				p.PopupPlay.SetSensitive(false)
				p.PopupDeleteAll.SetSensitive(false)
				p.PopupUnwatch.SetSensitive(false)
				p.PopupUnwatch.SetLabel(constSetAsWatched)
				p.PopupSave.SetSensitive(false)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupSubscriptions.SetSensitive(false)
				p.PopupToWatch.SetSensitive(true)
				p.PopupToDelete.SetSensitive(true)
			case constFilterModeToWatch:
				p.PopupDownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(true)
				p.PopupDeleteAll.SetSensitive(false)
				p.PopupUnwatch.SetSensitive(true)
				p.PopupUnwatch.SetLabel(constSetAsWatched)
				p.PopupSave.SetSensitive(true)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupSubscriptions.SetSensitive(true)
				p.PopupToWatch.SetSensitive(false)
				p.PopupToDelete.SetSensitive(true)
			case constFilterModeToDelete:
				p.PopupDownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(true)
				p.PopupDeleteAll.SetSensitive(true)
				p.PopupUnwatch.SetSensitive(true)
				p.PopupUnwatch.SetLabel(constSetAsUnwatched)
				p.PopupSave.SetSensitive(true)
				p.PopupSave.SetLabel(constSetAsSaved)
				p.PopupSubscriptions.SetSensitive(true)
				p.PopupToWatch.SetSensitive(true)
				p.PopupToDelete.SetSensitive(false)
			case constFilterModeSaved:
				p.PopupDownload.SetSensitive(false)
				p.PopupPlay.SetSensitive(true)
				p.PopupDeleteAll.SetSensitive(false)
				p.PopupUnwatch.SetSensitive(false)
				p.PopupUnwatch.SetLabel(constSetAsWatched)
				p.PopupSave.SetSensitive(true)
				p.PopupSave.SetLabel(constSetAsNotSaved)
				p.PopupSubscriptions.SetSensitive(true)
				p.PopupToWatch.SetSensitive(true)
				p.PopupToDelete.SetSensitive(false)
			}
			p.PopupMenu.PopupAtPointer(event)
		}
	})

	p.PopupRefresh.Connect("activate", func() {
		p.Parent.VideoList.Refresh("")
	})

	p.PopupDownload.Connect("activate", func() {
		treeview := p.Parent.VideoList.Treeview
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		p.Parent.VideoList.downloadVideo(video)
	})

	p.PopupPlay.Connect("activate", func() {
		treeview := p.Parent.VideoList.Treeview
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		p.Parent.VideoList.playVideo(video)
	})

	p.PopupDeleteAll.Connect("activate", func() {
		p.Parent.VideoList.DeleteWatchedVideos()
	})

	p.PopupSave.Connect("activate", func() {
		treeview := p.Parent.VideoList.Treeview
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		mode := p.PopupSave.GetLabel() == constSetAsSaved
		p.Parent.VideoList.setAsSaved(video, mode)
	})

	p.PopupUnwatch.Connect("activate", func() {
		treeview := p.Parent.VideoList.Treeview
		video := p.Parent.VideoList.getSelectedVideo(treeview)
		mode := p.PopupUnwatch.GetLabel() == constSetAsWatched
		p.Parent.VideoList.setAsWatched(video, mode)
	})
	p.PopupSubscriptions.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeSubscriptions)
	})
	p.PopupToWatch.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeToWatch)
	})
	p.PopupToDelete.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeToDelete)
	})
	p.PopupSaved.Connect("activate", func() {
		p.Parent.VideoList.SetFilterMode(constFilterModeSaved)
	})
}
