package main

import (
	"errors"

	"github.com/gotk3/gotk3/gtk"
)

// Toolbar : The toobar for SoftTube application
type Toolbar struct {
	Parent                 *SoftTube
	ToolbarSubscriptions   *gtk.ToggleToolButton
	ToolbarToWatch         *gtk.ToggleToolButton
	ToolbarToDelete        *gtk.ToggleToolButton
	ToolbarScrollToStart   *gtk.ToolButton
	ToolbarScrollToEnd     *gtk.ToolButton
	ToolbarKeepScrollToEnd *gtk.ToggleToolButton
	ToolbarRefresh         *gtk.ToolButton
	ToolbarDelete          *gtk.ToolButton
	ToolbarDeleteAll       *gtk.ToolButton
	ToolbarQuit            *gtk.ToolButton
}

// Load : Loads the toolbar
func (t *Toolbar) Load(builder *gtk.Builder) error {
	toggle, err := getToggleToolButton(builder, "toolbar_subscriptions")
	if err != nil {
		return err
	}
	t.ToolbarSubscriptions = toggle

	toggle, err = getToggleToolButton(builder, "toolbar_to_watch")
	if err != nil {
		return err
	}
	t.ToolbarToWatch = toggle

	toggle, err = getToggleToolButton(builder, "toolbar_to_delete")
	if err != nil {
		return err
	}
	t.ToolbarToDelete = toggle

	tool, err := getToolButton(builder, "toolbar_scroll_to_start")
	if err != nil {
		return err
	}
	t.ToolbarScrollToStart = tool

	tool, err = getToolButton(builder, "toolbar_scroll_to_end")
	if err != nil {
		return err
	}
	t.ToolbarScrollToEnd = tool

	toggleTool, err := getToggleToolButton(builder, "toolbar_keep_scroll_to_end")
	if err != nil {
		return err
	}
	t.ToolbarKeepScrollToEnd = toggleTool

	tool, err = getToolButton(builder, "toolbar_refresh_button")
	if err != nil {
		return err
	}
	t.ToolbarRefresh = tool

	tool, err = getToolButton(builder, "toolbar_delete_button")
	if err != nil {
		return err
	}
	t.ToolbarDelete = tool

	tool, err = getToolButton(builder, "toolbar_delete_all_button")
	if err != nil {
		return err
	}
	t.ToolbarDeleteAll = tool

	tool, err = getToolButton(builder, "toolbar_quit_button")
	if err != nil {
		return err
	}
	t.ToolbarQuit = tool

	return nil
}

// SetupEvents : Setup the toolbar events
func (t *Toolbar) SetupEvents() {
	t.ToolbarQuit.Connect("clicked", func() {
		gtk.MainQuit()
	})
	t.ToolbarRefresh.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.Refresh("")
	})
	t.ToolbarDeleteAll.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.DeleteWatchedVideos()
	})
	t.ToolbarSubscriptions.Connect("clicked", func() {
		if t.ToolbarSubscriptions.GetActive() {
			s := t.Parent
			t.ToolbarDeleteAll.SetSensitive(false)
			t.ToolbarToWatch.SetActive(false)
			t.ToolbarToDelete.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeSubscriptions)
		}
	})
	t.ToolbarToWatch.Connect("clicked", func() {
		if t.ToolbarToWatch.GetActive() {
			s := t.Parent
			t.ToolbarDeleteAll.SetSensitive(false)
			t.ToolbarSubscriptions.SetActive(false)
			t.ToolbarToDelete.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeToWatch)
		}
	})
	t.ToolbarToDelete.Connect("clicked", func() {
		if t.ToolbarToDelete.GetActive() {
			s := t.Parent
			t.ToolbarDeleteAll.SetSensitive(true)
			t.ToolbarSubscriptions.SetActive(false)
			t.ToolbarToWatch.SetActive(false)
			s.VideoList.SetFilterMode(constFilterModeToDelete)
		}
	})
	t.ToolbarScrollToStart.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.ScrollToStart()
	})
	t.ToolbarScrollToEnd.Connect("clicked", func() {
		s := t.Parent
		s.VideoList.ScrollToEnd()
	})
	t.ToolbarKeepScrollToEnd.Connect("clicked", func() {
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

func getToolButton(builder *gtk.Builder, name string) (*gtk.ToolButton, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.ToolButton); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk tool button")
}

func getToggleToolButton(builder *gtk.Builder, name string) (*gtk.ToggleToolButton, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.ToggleToolButton); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk toggle tool button")
}
