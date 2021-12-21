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

// Init initiates the toolbar
func (t *toolbar) Init(builder *framework.GtkBuilder) error {
	t.toolbarSubscriptions = builder.GetObject("toolbar_subscriptions").(*gtk.ToggleToolButton)
	t.toolbarDownloads = builder.GetObject("toolbar_downloads").(*gtk.ToggleToolButton)
	t.toolbarToWatch = builder.GetObject("toolbar_to_watch").(*gtk.ToggleToolButton)
	t.toolbarSaved = builder.GetObject("toolbar_saved").(*gtk.ToggleToolButton)
	t.toolbarToDelete = builder.GetObject("toolbar_to_delete").(*gtk.ToggleToolButton)
	t.toolbarScrollToStart = builder.GetObject("toolbar_scroll_to_start").(*gtk.ToolButton)
	t.toolbarScrollToEnd = builder.GetObject("toolbar_scroll_to_end").(*gtk.ToolButton)
	t.toolbarKeepScrollToEnd = builder.GetObject("toolbar_keep_scroll_to_end").(*gtk.ToggleToolButton)
	t.toolbarRefresh = builder.GetObject("toolbar_refresh_button").(*gtk.ToolButton)
	t.toolbarDeleteAll = builder.GetObject("toolbar_delete_all_button").(*gtk.ToolButton)
	t.toolbarQuit = builder.GetObject("toolbar_quit_button").(*gtk.ToolButton)

	t.SetupEvents()

	return nil
}

// SetupEvents : Set up the toolbar events
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
		t.parent.videoList.switchView(viewSubscriptions)
	})
	_ = t.toolbarDownloads.Connect("clicked", func() {
		t.parent.videoList.switchView(viewDownloads)
	})
	_ = t.toolbarToWatch.Connect("clicked", func() {
		t.parent.videoList.switchView(viewToWatch)
	})
	_ = t.toolbarSaved.Connect("clicked", func() {
		t.parent.videoList.switchView(viewSaved)
	})
	_ = t.toolbarToDelete.Connect("clicked", func() {
		t.parent.videoList.switchView(viewToDelete)
	})
	_ = t.toolbarScrollToStart.Connect("clicked", func() {
		s := t.parent
		s.videoList.scroll.toStart()
	})
	_ = t.toolbarScrollToEnd.Connect("clicked", func() {
		s := t.parent
		s.videoList.scroll.toEnd()
	})
	_ = t.toolbarKeepScrollToEnd.Connect("clicked", func() {
		if t.toolbarKeepScrollToEnd.GetActive() {
			s := t.parent
			s.videoList.keepScrollToEnd = true
			s.videoList.scroll.toEnd()
		} else {
			s := t.parent
			s.videoList.keepScrollToEnd = false
			s.videoList.scroll.toStart()
		}
	})
}
