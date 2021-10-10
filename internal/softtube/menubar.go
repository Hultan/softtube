package softtube

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
)

// menuBar : The SoftTube menu bar
type menuBar struct {
	parent        *SoftTube
	menuFileQuit  *gtk.MenuItem
	menuHelpAbout *gtk.MenuItem
}

// Load : Loads the toolbar
func (m *menuBar) Load(builder *framework.GtkBuilder) error {
	menuItem := builder.GetObject("menu_file_quit").(*gtk.MenuItem)
	m.menuFileQuit = menuItem

	menuItem = builder.GetObject("menu_help_about").(*gtk.MenuItem)
	m.menuHelpAbout = menuItem

	return nil
}

// SetupEvents : Setup the toolbar events
func (m *menuBar) SetupEvents() {
	_ = m.menuFileQuit.Connect("activate", func() {
		gtk.MainQuit()
	})
	_ = m.menuHelpAbout.Connect("activate", func() {
	})
}