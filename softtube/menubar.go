package main

import (
	"github.com/gotk3/gotk3/gtk"
	gtkhelper "github.com/hultan/softteam/gtk"
)

// MenuBar : The SoftTube menu bar
type MenuBar struct {
	Parent        *SoftTube
	MenuFileQuit  *gtk.MenuItem
	MenuHelpAbout *gtk.MenuItem
}

// Load : Loads the toolbar
func (m *MenuBar) Load(helper *gtkhelper.GtkHelper) error {
	menuItem, err := helper.GetMenuItem("menu_file_quit")
	if err != nil {
		return err
	}
	m.MenuFileQuit = menuItem

	menuItem, err = helper.GetMenuItem("menu_help_about")
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
