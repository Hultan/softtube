package softtube

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam/framework"
)

// MenuBar : The SoftTube menu bar
type MenuBar struct {
	Parent        *SoftTube
	MenuFileQuit  *gtk.MenuItem
	MenuHelpAbout *gtk.MenuItem
}

// Load : Loads the toolbar
func (m *MenuBar) Load(builder *framework.GtkBuilder) error {
	menuItem := builder.GetObject("menu_file_quit").(*gtk.MenuItem)
	m.MenuFileQuit = menuItem

	menuItem = builder.GetObject("menu_help_about").(*gtk.MenuItem)
	m.MenuHelpAbout = menuItem

	return nil
}

// SetupEvents : Setup the toolbar events
func (m *MenuBar) SetupEvents() {
	_= m.MenuFileQuit.Connect("activate", func() {
		gtk.MainQuit()
	})
	_ = m.MenuHelpAbout.Connect("activate", func() {
	})
}
