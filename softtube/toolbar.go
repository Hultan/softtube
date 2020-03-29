package main

import (
	"errors"

	"github.com/gotk3/gotk3/gtk"
)

// Toolbar : The toobar for SoftTube application
type Toolbar struct {
	ToolbarSubscriptions   *gtk.ToolButton
	ToolbarToWatch         *gtk.ToolButton
	ToolbarToDelete        *gtk.ToolButton
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
	tool, err := getToolButton(builder, "toolbar_subscriptions")
	if err != nil {
		return err
	}
	t.ToolbarSubscriptions = tool

	tool, err = getToolButton(builder, "toolbar_to_watch")
	if err != nil {
		return err
	}
	t.ToolbarToWatch = tool

	tool, err = getToolButton(builder, "toolbar_to_delete")
	if err != nil {
		return err
	}
	t.ToolbarToDelete = tool

	tool, err = getToolButton(builder, "toolbar_scroll_to_start")
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
