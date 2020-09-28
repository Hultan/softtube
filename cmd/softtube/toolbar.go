package main

import (
	"github.com/gotk3/gotk3/gtk"
	gtkHelper "github.com/hultan/softteam-tools/pkg/gtk-helper"
)

// Toolbar : The toolbar for SoftTube application
type Toolbar struct {
	Parent                 *SoftTube
	ToolbarSubscriptions   *gtk.ToggleToolButton
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
func (t *Toolbar) Load(helper *gtkHelper.GtkHelper) error {
	toggle, err := helper.GetToggleToolButton("toolbar_subscriptions")
	if err != nil {
		return err
	}
	t.ToolbarSubscriptions = toggle

	toggle, err = helper.GetToggleToolButton("toolbar_to_watch")
	if err != nil {
		return err
	}
	t.ToolbarToWatch = toggle

	toggle, err = helper.GetToggleToolButton("toolbar_to_delete")
	if err != nil {
		return err
	}
	t.ToolbarToDelete = toggle

	toggle, err = helper.GetToggleToolButton("toolbar_saved")
	if err != nil {
		return err
	}
	t.ToolbarSaved = toggle

	tool, err := helper.GetToolButton("toolbar_scroll_to_start")
	if err != nil {
		return err
	}
	t.ToolbarScrollToStart = tool

	tool, err = helper.GetToolButton("toolbar_scroll_to_end")
	if err != nil {
		return err
	}
	t.ToolbarScrollToEnd = tool

	toggleTool, err := helper.GetToggleToolButton("toolbar_keep_scroll_to_end")
	if err != nil {
		return err
	}
	t.ToolbarKeepScrollToEnd = toggleTool

	tool, err = helper.GetToolButton("toolbar_refresh_button")
	if err != nil {
		return err
	}
	t.ToolbarRefresh = tool

	tool, err = helper.GetToolButton("toolbar_delete_all_button")
	if err != nil {
		return err
	}
	t.ToolbarDeleteAll = tool

	tool, err = helper.GetToolButton("toolbar_quit_button")
	if err != nil {
		return err
	}
	t.ToolbarQuit = tool

	return nil
}

// SetupEvents : Setup the toolbar events
func (t *Toolbar) SetupEvents() {
	_,_ = t.ToolbarQuit.Connect("clicked", func() {
		gtk.MainQuit()
	})
	_,_ = t.ToolbarRefresh.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.Refresh("")
	})
	_,_ = t.ToolbarDeleteAll.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.DeleteWatchedVideos()
	})
	_,_ = t.ToolbarSubscriptions.Connect("clicked", func() {
		if t.ToolbarSubscriptions.GetActive() {
			s := t.Parent
			t.ToolbarDeleteAll.SetSensitive(false)
			t.ToolbarToWatch.SetActive(false)
			t.ToolbarToDelete.SetActive(false)
			t.ToolbarSaved.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeSubscriptions)
		}
	})
	_,_ = t.ToolbarToWatch.Connect("clicked", func() {
		if t.ToolbarToWatch.GetActive() {
			s := t.Parent
			t.ToolbarDeleteAll.SetSensitive(false)
			t.ToolbarSubscriptions.SetActive(false)
			t.ToolbarToDelete.SetActive(false)
			t.ToolbarSaved.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeToWatch)
		}
	})
	_,_ = t.ToolbarToDelete.Connect("clicked", func() {
		if t.ToolbarToDelete.GetActive() {
			s := t.Parent
			t.ToolbarDeleteAll.SetSensitive(true)
			t.ToolbarSubscriptions.SetActive(false)
			t.ToolbarToWatch.SetActive(false)
			t.ToolbarSaved.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeToDelete)
		}
	})
	_,_ = t.ToolbarSaved.Connect("clicked", func() {
		if t.ToolbarSaved.GetActive() {
			s := t.Parent
			t.ToolbarDeleteAll.SetSensitive(false)
			t.ToolbarSubscriptions.SetActive(false)
			t.ToolbarToWatch.SetActive(false)
			t.ToolbarToDelete.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeSaved)
		}
	})
	_,_ = t.ToolbarScrollToStart.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.ScrollToStart()
	})
	_,_ = t.ToolbarScrollToEnd.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.ScrollToEnd()
	})
	_,_ = t.ToolbarKeepScrollToEnd.Connect("clicked", func() {
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
