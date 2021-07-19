package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam/framework"
)

// Toolbar : The toolbar for SoftTube application
type Toolbar struct {
	Parent                 *SoftTube
	ToolbarSubscriptions   *gtk.ToggleToolButton
	ToolbarDownloads       *gtk.ToggleToolButton
	ToolbarToWatch         *gtk.ToggleToolButton
	ToolbarToDelete        *gtk.ToggleToolButton
	ToolbarSaved           *gtk.ToggleToolButton
	ToolbarScrollToStart   *gtk.ToolButton
	ToolbarScrollToEnd     *gtk.ToolButton
	ToolbarKeepScrollToEnd *gtk.ToggleToolButton
	ToolbarRefresh         *gtk.ToolButton
	ToolbarDelete          *gtk.ToolButton
	ToolbarDeleteAll       *gtk.ToolButton
	ToolbarQuit            *gtk.ToolButton
}

// Load : Loads the toolbar
func (t *Toolbar) Load(builder *framework.GtkBuilder) error {
	toggle := builder.GetObject("toolbar_subscriptions").(*gtk.ToggleToolButton)
	t.ToolbarSubscriptions = toggle

	toggle = builder.GetObject("toolbar_failed").(*gtk.ToggleToolButton)
	t.ToolbarDownloads = toggle

	toggle = builder.GetObject("toolbar_to_watch").(*gtk.ToggleToolButton)
	t.ToolbarToWatch = toggle

	toggle = builder.GetObject("toolbar_to_delete").(*gtk.ToggleToolButton)
	t.ToolbarToDelete = toggle

	toggle = builder.GetObject("toolbar_saved").(*gtk.ToggleToolButton)
	t.ToolbarSaved = toggle

	tool := builder.GetObject("toolbar_scroll_to_start").(*gtk.ToolButton)
	t.ToolbarScrollToStart = tool

	tool = builder.GetObject("toolbar_scroll_to_end").(*gtk.ToolButton)
	t.ToolbarScrollToEnd = tool

	toggle = builder.GetObject("toolbar_keep_scroll_to_end").(*gtk.ToggleToolButton)
	t.ToolbarKeepScrollToEnd = toggle

	tool = builder.GetObject("toolbar_refresh_button").(*gtk.ToolButton)
	t.ToolbarRefresh = tool

	tool = builder.GetObject("toolbar_delete_all_button").(*gtk.ToolButton)
	t.ToolbarDeleteAll = tool

	tool = builder.GetObject("toolbar_quit_button").(*gtk.ToolButton)
	t.ToolbarQuit = tool

	return nil
}

// SetupEvents : Setup the toolbar events
func (t *Toolbar) SetupEvents() {
	_ = t.ToolbarQuit.Connect("clicked", func() {
		gtk.MainQuit()
	})
	_ = t.ToolbarRefresh.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.Refresh("")
	})
	_ = t.ToolbarDeleteAll.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.DeleteWatchedVideos()
	})
	_ = t.ToolbarSubscriptions.Connect("clicked", func() {
		if t.ToolbarSubscriptions.GetActive() {
			s := t.Parent
			t.ToolbarDeleteAll.SetSensitive(false)
			t.ToolbarDownloads.SetActive(false)
			t.ToolbarToWatch.SetActive(false)
			t.ToolbarToDelete.SetActive(false)
			t.ToolbarSaved.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeSubscriptions)
		}
	})
	_ = t.ToolbarDownloads.Connect("clicked", func() {
		if t.ToolbarDownloads.GetActive() {
			s := t.Parent
			t.ToolbarDeleteAll.SetSensitive(false)
			t.ToolbarSubscriptions.SetActive(false)
			t.ToolbarToWatch.SetActive(false)
			t.ToolbarToDelete.SetActive(false)
			t.ToolbarSaved.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeDownloads)
		}
	})
	_ = t.ToolbarToWatch.Connect("clicked", func() {
		if t.ToolbarToWatch.GetActive() {
			s := t.Parent
			t.ToolbarDownloads.SetActive(false)
			t.ToolbarDeleteAll.SetSensitive(false)
			t.ToolbarSubscriptions.SetActive(false)
			t.ToolbarToDelete.SetActive(false)
			t.ToolbarSaved.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeToWatch)
		}
	})
	_ = t.ToolbarToDelete.Connect("clicked", func() {
		if t.ToolbarToDelete.GetActive() {
			s := t.Parent
			t.ToolbarDownloads.SetActive(false)
			t.ToolbarDeleteAll.SetSensitive(true)
			t.ToolbarSubscriptions.SetActive(false)
			t.ToolbarToWatch.SetActive(false)
			t.ToolbarSaved.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeToDelete)
		}
	})
	_ = t.ToolbarSaved.Connect("clicked", func() {
		if t.ToolbarSaved.GetActive() {
			s := t.Parent
			t.ToolbarDownloads.SetActive(false)
			t.ToolbarDeleteAll.SetSensitive(false)
			t.ToolbarSubscriptions.SetActive(false)
			t.ToolbarToWatch.SetActive(false)
			t.ToolbarToDelete.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeSaved)
		}
	})
	_ = t.ToolbarScrollToStart.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.ScrollToStart()
	})
	_ = t.ToolbarScrollToEnd.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.ScrollToEnd()
	})
	_ = t.ToolbarKeepScrollToEnd.Connect("clicked", func() {
		if t.ToolbarKeepScrollToEnd.GetActive() {
			s := t.Parent
			s.VideoList.KeepScrollToEnd = true
			s.VideoList.ScrollToEnd()
		} else {
			s := t.Parent
			s.VideoList.KeepScrollToEnd = false
			s.VideoList.ScrollToStart()
		}
	})
}
