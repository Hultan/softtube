package main

import (
	"errors"

	"github.com/gotk3/gotk3/gtk"
)

// MenuBar : The SoftTube menu bar
type MenuBar struct {
	Parent        *SoftTube
	MenuFileQuit  *gtk.MenuItem
	MenuHelpAbout *gtk.MenuItem
}

// Load : Loads the toolbar
func (m *MenuBar) Load(builder *gtk.Builder) error {
	menuItem, err := getMenuItem(builder, "menu_file_quit")
	if err != nil {
		return err
	}
	m.MenuFileQuit = menuItem

	menuItem, err = getMenuItem(builder, "menu_help_about")
	if err != nil {
		return err
	}
	m.MenuHelpAbout = menuItem

	return nil
}

// SetupEvents : Setup the toolbar events
func (m *MenuBar) SetupEvents() {
	m.MenuFileQuit.Connect("activate", func() {
		gtk.MainQuit()
	})
	m.MenuHelpAbout.Connect("activate", func() {
	})
}

func getMenuItem(builder *gtk.Builder, name string) (*gtk.MenuItem, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if menuItem, ok := obj.(*gtk.MenuItem); ok {
		return menuItem, nil
	}

	return nil, errors.New("not a gtk menu item")
}
