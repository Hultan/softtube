package softtube

import (
	"github.com/gotk3/gotk3/gtk"
)

// toolbar is the toolbar for SoftTube application
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

// Init initiates the toolbar
func (t *toolbar) Init() error {
	t.toolbarSubscriptions = GetObject[*gtk.ToggleToolButton]("toolbar_subscriptions")
	t.toolbarDownloads = GetObject[*gtk.ToggleToolButton]("toolbar_downloads")
	t.toolbarToWatch = GetObject[*gtk.ToggleToolButton]("toolbar_to_watch")
	t.toolbarSaved = GetObject[*gtk.ToggleToolButton]("toolbar_saved")
	t.toolbarToDelete = GetObject[*gtk.ToggleToolButton]("toolbar_to_delete")
	t.toolbarScrollToStart = GetObject[*gtk.ToolButton]("toolbar_scroll_to_start")
	t.toolbarScrollToEnd = GetObject[*gtk.ToolButton]("toolbar_scroll_to_end")
	t.toolbarKeepScrollToEnd = GetObject[*gtk.ToggleToolButton]("toolbar_keep_scroll_to_end")
	t.toolbarRefresh = GetObject[*gtk.ToolButton]("toolbar_refresh_button")
	t.toolbarDeleteAll = GetObject[*gtk.ToolButton]("toolbar_delete_all_button")
	t.toolbarQuit = GetObject[*gtk.ToolButton]("toolbar_quit_button")

	t.SetupEvents()

	return nil
}

// SetupEvents sets up the toolbar events
func (t *toolbar) SetupEvents() {
	_ = t.toolbarQuit.Connect(
		"clicked", func() {
			gtk.MainQuit()
		},
	)
	_ = t.toolbarRefresh.Connect(
		"clicked", func() {
			s := t.parent
			s.videoList.Refresh("")
		},
	)
	_ = t.toolbarDeleteAll.Connect(
		"clicked", func() {
			s := t.parent
			s.videoList.DeleteWatchedVideos()
		},
	)
	_ = t.toolbarSubscriptions.Connect(
		"clicked", func() {
			t.parent.videoList.switchView(viewSubscriptions)
		},
	)
	_ = t.toolbarDownloads.Connect(
		"clicked", func() {
			t.parent.videoList.switchView(viewDownloads)
		},
	)
	_ = t.toolbarToWatch.Connect(
		"clicked", func() {
			t.parent.videoList.switchView(viewToWatch)
		},
	)
	_ = t.toolbarSaved.Connect(
		"clicked", func() {
			t.parent.videoList.switchView(viewSaved)
		},
	)
	_ = t.toolbarToDelete.Connect(
		"clicked", func() {
			t.parent.videoList.switchView(viewToDelete)
		},
	)
	_ = t.toolbarScrollToStart.Connect(
		"clicked", func() {
			s := t.parent
			s.videoList.scroll.toStart()
		},
	)
	_ = t.toolbarScrollToEnd.Connect(
		"clicked", func() {
			s := t.parent
			s.videoList.scroll.toEnd()
		},
	)
	_ = t.toolbarKeepScrollToEnd.Connect(
		"clicked", func() {
			if t.toolbarKeepScrollToEnd.GetActive() {
				s := t.parent
				s.videoList.keepScrollToEnd = true
				s.videoList.scroll.toEnd()
			} else {
				s := t.parent
				s.videoList.keepScrollToEnd = false
				s.videoList.scroll.toStart()
			}
		},
	)
}
