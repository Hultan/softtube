package softtube

import (
	"log"
	"sync"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// popupMenu handles of the video list popupmenu
type popupMenu struct {
	parent                 *SoftTube
	popupMenu              *gtk.Menu
	popupRefresh           *gtk.MenuItem
	popupDownload          *gtk.MenuItem
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
	popupSearchChannel     *gtk.MenuItem
	popupSearchVideo       *gtk.MenuItem
}

// Init loads the popup menu
func (p *popupMenu) Init() error {
	p.popupMenu = GetObject[*gtk.Menu]("popupmenu")
	p.popupRefresh = GetObject[*gtk.MenuItem]("popup_refresh")
	p.popupDownload = GetObject[*gtk.MenuItem]("popup_download")
	p.popupRedownloadVideo = GetObject[*gtk.MenuItem]("popup_redownload_failedvideo")
	p.popupRedownloadVideos = GetObject[*gtk.MenuItem]("popup_redownload_failedvideos")
	p.popupPlay = GetObject[*gtk.MenuItem]("popup_play")
	p.popupGetDuration = GetObject[*gtk.MenuItem]("popup_get_duration")
	p.popupGetVideoID = GetObject[*gtk.MenuItem]("popup_get_videoid")
	p.popupGetThumbnail = GetObject[*gtk.MenuItem]("popup_get_thumbnail")
	p.popupDeleteAll = GetObject[*gtk.MenuItem]("popup_delete_all")
	p.popupUnwatch = GetObject[*gtk.MenuItem]("popup_unwatch")
	p.popupSave = GetObject[*gtk.MenuItem]("popup_save")
	p.popupViewSubscriptions = GetObject[*gtk.MenuItem]("popup_view_subscriptions")
	p.popupViewDownloads = GetObject[*gtk.MenuItem]("popup_view_downloads")
	p.popupViewToWatch = GetObject[*gtk.MenuItem]("popup_view_to_watch")
	p.popupViewSaved = GetObject[*gtk.MenuItem]("popup_view_saved")
	p.popupViewToDelete = GetObject[*gtk.MenuItem]("popup_view_to_delete")
	p.popupSearchChannel = GetObject[*gtk.MenuItem]("popup_search_channel_name")
	p.popupSearchVideo = GetObject[*gtk.MenuItem]("popup_search_video_title")

	p.SetupEvents()

	return nil
}

// SetupEvents sets up the toolbar events
func (p *popupMenu) SetupEvents() {
	_ = p.parent.videoList.treeView.Connect(
		"button-press-event", func(treeview *gtk.TreeView, event *gdk.Event) bool {
			// This code solves the problem with the last selected row
			// getting deselected when you open the context menu
			selection, err := treeview.GetSelection()
			if err != nil {
				log.Fatal("Unable to get TreeSelection:", err)
			}

			buttonEvent := gdk.EventButtonNewFromEvent(event)
			if buttonEvent.Button() == gdk.BUTTON_SECONDARY {
				// Check if the clicked row is already selected
				path, _, _, _, _ := treeview.GetPathAtPos(int(buttonEvent.X()), int(buttonEvent.Y()))
				if path != nil {
					isSelected := selection.PathIsSelected(path)
					if !isSelected {
						selection.SelectPath(path) // Select the row if it’s not already selected
					}
				}
				// END: This code solves the problem with the last selected row
				// getting deselected when you open the context menu

				videoSelected := p.parent.videoList.videoFunctions.getSelectedVideos(p.parent.videoList.treeView) != nil
				view := func(subscription, download, toWatch, saved, toDelete bool) {
					p.popupViewSubscriptions.SetSensitive(subscription)
					p.popupViewDownloads.SetSensitive(download)
					p.popupViewToWatch.SetSensitive(toWatch)
					p.popupViewSaved.SetSensitive(saved)
					p.popupViewToDelete.SetSensitive(toDelete)
				}

				p.popupDownload.SetSensitive(false)
				p.popupRedownloadVideo.SetSensitive(false)
				p.popupRedownloadVideos.SetSensitive(false)

				switch p.parent.videoList.currentView {
				case viewSubscriptions:
					p.popupDownload.SetSensitive(videoSelected)

					p.popupPlay.SetSensitive(false)
					p.popupGetDuration.SetSensitive(videoSelected)
					p.popupGetVideoID.SetSensitive(videoSelected)
					p.popupGetThumbnail.SetSensitive(videoSelected)
					p.popupDeleteAll.SetVisible(false)
					p.popupUnwatch.SetSensitive(videoSelected)
					p.popupUnwatch.SetLabel(constSetAsNotDownloaded)
					p.popupSave.SetSensitive(false)
					p.popupSave.SetLabel(constSetAsSaved)

					view(false, true, true, true, true)
				case viewDownloads:
					p.popupRedownloadVideo.SetSensitive(videoSelected)
					p.popupRedownloadVideos.SetSensitive(videoSelected)

					p.popupPlay.SetSensitive(false)
					p.popupGetDuration.SetSensitive(videoSelected)
					p.popupGetVideoID.SetSensitive(videoSelected)
					p.popupGetThumbnail.SetSensitive(videoSelected)
					p.popupDeleteAll.SetVisible(false)
					p.popupUnwatch.SetSensitive(true)
					p.popupUnwatch.SetLabel(constSetAsNotDownloaded)
					p.popupSave.SetSensitive(false)
					p.popupSave.SetLabel(constSetAsSaved)

					view(true, false, true, true, true)
				case viewToWatch:
					p.popupPlay.SetSensitive(videoSelected)
					p.popupGetDuration.SetSensitive(false)
					p.popupGetVideoID.SetSensitive(videoSelected)
					p.popupGetThumbnail.SetSensitive(videoSelected)
					p.popupDeleteAll.SetVisible(false)
					p.popupUnwatch.SetSensitive(videoSelected)
					p.popupUnwatch.SetLabel(constSetAsWatched)
					p.popupSave.SetSensitive(videoSelected)
					p.popupSave.SetLabel(constSetAsSaved)

					view(true, true, false, true, true)
				case viewToDelete:
					p.popupPlay.SetSensitive(videoSelected)
					p.popupGetDuration.SetSensitive(false)
					p.popupGetVideoID.SetSensitive(videoSelected)
					p.popupGetThumbnail.SetSensitive(false)
					p.popupDeleteAll.SetVisible(true)
					p.popupUnwatch.SetSensitive(videoSelected)
					p.popupUnwatch.SetLabel(constSetAsUnwatched)
					p.popupSave.SetSensitive(videoSelected)
					p.popupSave.SetLabel(constSetAsSaved)

					view(true, true, true, true, false)
				case viewSaved:
					p.popupPlay.SetSensitive(videoSelected)
					p.popupGetDuration.SetSensitive(false)
					p.popupGetVideoID.SetSensitive(videoSelected)
					p.popupGetThumbnail.SetSensitive(videoSelected)
					p.popupDeleteAll.SetVisible(false)
					p.popupUnwatch.SetSensitive(true)
					p.popupUnwatch.SetLabel(constSetAsWatched)
					p.popupSave.SetSensitive(true)
					p.popupSave.SetLabel(constSetAsNotSaved)

					view(true, true, true, false, true)
				}

				p.popupMenu.PopupAtPointer(event)

				// By returning true/false here, we stop event propagation
				// so that the row the user clicks does not get deselected
				// by the right click.
				return true
			}

			return false
		},
	)

	_ = p.popupRefresh.Connect(
		"activate", func() {
			p.parent.videoList.Refresh("")
		},
	)

	_ = p.popupDownload.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				for _, video := range selectedVideos {
					_ = p.parent.videoList.videoFunctions.download(video, true)
				}
			}
		},
	)
	_ = p.popupRedownloadVideo.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				for _, video := range selectedVideos {
					_ = p.parent.videoList.videoFunctions.download(video, false)
				}
			}
		},
	)
	_ = p.popupRedownloadVideos.Connect(
		"activate", func() {
			vids, err := p.parent.DB.Videos.GetVideos(true, p.parent.videoList.currentView == viewSaved)
			if err != nil {
				p.parent.Logger.Error.Println(err)
				return
			}
			for key := range vids {
				vid := &vids[key]
				if vid != nil {
					_ = p.parent.videoList.videoFunctions.download(vid, false)
				}
			}
		},
	)

	_ = p.popupPlay.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				// You can only play one video at the time so we only play the first
				p.parent.videoList.videoFunctions.play(selectedVideos[0])
			}
		},
	)

	_ = p.popupGetDuration.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				p.parent.downloadDurations(selectedVideos)
			}
		},
	)

	_ = p.popupGetVideoID.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				// We can only copy one ID at a time, so we only copy the first
				clipboard, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
				if err != nil {
					p.parent.Logger.Error.Println(err)
					return
				}
				clipboard.SetText(selectedVideos[0].ID)
			}
		},
	)

	_ = p.popupGetThumbnail.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				p.parent.downloadThumbnails(selectedVideos)
			}
		},
	)

	_ = p.popupDeleteAll.Connect(
		"activate", func() {
			p.parent.videoList.DeleteWatchedVideos()
		},
	)

	_ = p.popupSave.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				for _, video := range selectedVideos {
					mode := p.popupSave.GetLabel() == constSetAsSaved
					p.parent.videoList.videoFunctions.setAsSaved(video, mode)
				}
			}
		},
	)

	_ = p.popupUnwatch.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				// TODO : Change to a map?
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

				var wg sync.WaitGroup
				wg.Add(len(selectedVideos))
				for _, video := range selectedVideos {
					p.parent.videoList.videoFunctions.setVideoStatus(video, mode, &wg)
				}
				wg.Wait()
				p.parent.videoList.Refresh("")
			}
		},
	)
	_ = p.popupViewSubscriptions.Connect(
		"activate", func() {
			p.parent.videoList.switchView(viewSubscriptions)
		},
	)
	_ = p.popupViewDownloads.Connect(
		"activate", func() {
			p.parent.videoList.switchView(viewDownloads)
		},
	)
	_ = p.popupViewToWatch.Connect(
		"activate", func() {
			p.parent.videoList.switchView(viewToWatch)
		},
	)
	_ = p.popupViewSaved.Connect(
		"activate", func() {
			p.parent.videoList.switchView(viewSaved)
		},
	)
	_ = p.popupViewToDelete.Connect(
		"activate", func() {
			p.parent.videoList.switchView(viewToDelete)
		},
	)
	_ = p.popupSearchChannel.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				p.parent.searchBar.searchEntry.SetText(selectedVideos[0].SubscriptionName)
				p.parent.videoList.Search(selectedVideos[0].SubscriptionName)
			}
		},
	)
	_ = p.popupSearchVideo.Connect(
		"activate", func() {
			treeview := p.parent.videoList.treeView
			selectedVideos := p.parent.videoList.videoFunctions.getSelectedVideos(treeview)
			if selectedVideos != nil {
				p.parent.searchBar.searchEntry.SetText(selectedVideos[0].Title)
				p.parent.videoList.Search(selectedVideos[0].Title)
			}
		},
	)
}
