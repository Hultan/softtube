package softtube

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
)

// toolbar : The toolbar for SoftTube application
type toolbar struct {
	parent                 *SoftTube
	toolbarSubscriptions   *gtk.ToggleToolButton
	toolbarDownloads       *gtk.ToggleToolButton
	toolbarToWatch         *gtk.ToggleToolButton
	toolbarToDelete        *gtk.ToggleToolButton
	toolbarSaved           *gtk.ToggleToolButton
	toolbarScrollToStart   *gtk.ToolButton
	toolbarScrollToEnd     *gtk.ToolButton
	toolbarKeepScrollToEnd *gtk.ToggleToolButton
	toolbarRefresh         *gtk.ToolButton
	toolbarDelete          *gtk.ToolButton
	toolbarDeleteAll       *gtk.ToolButton
	toolbarQuit            *gtk.ToolButton
}

// Load : Loads the toolbar
func (t *toolbar) Load(builder *framework.GtkBuilder) error {
	toggle := builder.GetObject("toolbar_subscriptions").(*gtk.ToggleToolButton)
	t.toolbarSubscriptions = toggle

	toggle = builder.GetObject("toolbar_failed").(*gtk.ToggleToolButton)
	t.toolbarDownloads = toggle

	toggle = builder.GetObject("toolbar_to_watch").(*gtk.ToggleToolButton)
	t.toolbarToWatch = toggle

	toggle = builder.GetObject("toolbar_to_delete").(*gtk.ToggleToolButton)
	t.toolbarToDelete = toggle

	toggle = builder.GetObject("toolbar_saved").(*gtk.ToggleToolButton)
	t.toolbarSaved = toggle

	tool := builder.GetObject("toolbar_scroll_to_start").(*gtk.ToolButton)
	t.toolbarScrollToStart = tool

	tool = builder.GetObject("toolbar_scroll_to_end").(*gtk.ToolButton)
	t.toolbarScrollToEnd = tool

	toggle = builder.GetObject("toolbar_keep_scroll_to_end").(*gtk.ToggleToolButton)
	t.toolbarKeepScrollToEnd = toggle

	tool = builder.GetObject("toolbar_refresh_button").(*gtk.ToolButton)
	t.toolbarRefresh = tool

	tool = builder.GetObject("toolbar_delete_all_button").(*gtk.ToolButton)
	t.toolbarDeleteAll = tool

	tool = builder.GetObject("toolbar_quit_button").(*gtk.ToolButton)
	t.toolbarQuit = tool

	return nil
}

// SetupEvents : Setup the toolbar events
func (t *toolbar) SetupEvents() {
	_ = t.toolbarQuit.Connect("clicked", func() {
		gtk.MainQuit()
	})
	_ = t.toolbarRefresh.Connect("clicked", func() {
		s := t.parent
		s.videoList.Refresh("")
	})
	_ = t.toolbarDeleteAll.Connect("clicked", func() {
		s := t.parent
		s.videoList.DeleteWatchedVideos()
	})
	_ = t.toolbarSubscriptions.Connect("clicked", func() {
		if t.toolbarSubscriptions.GetActive() {
			s := t.parent
			t.toolbarDeleteAll.SetSensitive(false)
			t.toolbarDownloads.SetActive(false)
			t.toolbarToWatch.SetActive(false)
			t.toolbarToDelete.SetActive(false)
			t.toolbarSaved.SetActive(false)
			s.videoList.SetFilterMode(constFilterModeSubscriptions)
		}
	})
	_ = t.toolbarDownloads.Connect("clicked", func() {
		if t.toolbarDownloads.GetActive() {
			s := t.parent
			t.toolbarDeleteAll.SetSensitive(false)
			t.toolbarSubscriptions.SetActive(false)
			t.toolbarToWatch.SetActive(false)
			t.toolbarToDelete.SetActive(false)
			t.toolbarSaved.SetActive(false)
			s.videoList.SetFilterMode(constFilterModeDownloads)
		}
	})
	_ = t.toolbarToWatch.Connect("clicked", func() {
		if t.toolbarToWatch.GetActive() {
			s := t.parent
			t.toolbarDownloads.SetActive(false)
			t.toolbarDeleteAll.SetSensitive(false)
			t.toolbarSubscriptions.SetActive(false)
			t.toolbarToDelete.SetActive(false)
			t.toolbarSaved.SetActive(false)
			s.videoList.SetFilterMode(constFilterModeToWatch)
		}
	})
	_ = t.toolbarToDelete.Connect("clicked", func() {
		if t.toolbarToDelete.GetActive() {
			s := t.parent
			t.toolbarDownloads.SetActive(false)
			t.toolbarDeleteAll.SetSensitive(true)
			t.toolbarSubscriptions.SetActive(false)
			t.toolbarToWatch.SetActive(false)
			t.toolbarSaved.SetActive(false)
			s.videoList.SetFilterMode(constFilterModeToDelete)
		}
	})
	_ = t.toolbarSaved.Connect("clicked", func() {
		if t.toolbarSaved.GetActive() {
			s := t.parent
			t.toolbarDownloads.SetActive(false)
			t.toolbarDeleteAll.SetSensitive(false)
			t.toolbarSubscriptions.SetActive(false)
			t.toolbarToWatch.SetActive(false)
			t.toolbarToDelete.SetActive(false)
			s.videoList.SetFilterMode(constFilterModeSaved)
		}
	})
	_ = t.toolbarScrollToStart.Connect("clicked", func() {
		s := t.parent
		s.videoList.ScrollToStart()
	})
	_ = t.toolbarScrollToEnd.Connect("clicked", func() {
		s := t.parent
		s.videoList.ScrollToEnd()
	})
	_ = t.toolbarKeepScrollToEnd.Connect("clicked", func() {
		if t.toolbarKeepScrollToEnd.GetActive() {
			s := t.parent
			s.videoList.keepScrollToEnd = true
			s.videoList.ScrollToEnd()
		} else {
			s := t.parent
			s.videoList.keepScrollToEnd = false
			s.videoList.ScrollToStart()
		}
	})
}
