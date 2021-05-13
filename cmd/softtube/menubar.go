package main

import (
	"github.com/gotk3/gotk3/gtk"
)

// MenuBar : The SoftTube menu bar
type MenuBar struct {
	Parent        *SoftTube
	MenuFileQuit  *gtk.MenuItem
	MenuHelpAbout *gtk.MenuItem
}

// Load : Loads the toolbar
func (m *MenuBar) Load(builder *SoftBuilder) error {
	menuItem := builder.getObject("menu_file_quit").(*gtk.MenuItem)
	m.MenuFileQuit = menuItem

	menuItem = builder.getObject("menu_help_about").(*gtk.MenuItem)
	m.MenuHelpAbout = menuItem

	return nil
}

// SetupEvents : Setup the toolbar events
func (m *MenuBar) SetupEvents() {
	_,_ = m.MenuFileQuit.Connect("activate", func() {
		gtk.MainQuit()
	})
	_,_ = m.MenuHelpAbout.Connect("activate", func() {
	})
}
