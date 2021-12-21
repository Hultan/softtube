package softtube

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
)

// menuBar : The SoftTube menu bar
type menuBar struct {
	parent                *SoftTube
	menuFileQuit          *gtk.MenuItem
	menuHelpAbout         *gtk.MenuItem
	menuViewSubscriptions *gtk.RadioMenuItem
	menuViewDownloads     *gtk.RadioMenuItem
	menuViewToWatch       *gtk.RadioMenuItem
	menuViewSaved         *gtk.RadioMenuItem
	menuViewToDelete      *gtk.RadioMenuItem
}

// Init initiates the menu bar
func (m *menuBar) Init(builder *framework.GtkBuilder) error {
	m.menuFileQuit = builder.GetObject("menu_file_quit").(*gtk.MenuItem)
	m.menuHelpAbout = builder.GetObject("menu_help_about").(*gtk.MenuItem)

	m.menuViewSubscriptions = builder.GetObject("menu_view_subscriptions").(*gtk.RadioMenuItem)
	m.menuViewDownloads = builder.GetObject("menu_view_downloads").(*gtk.RadioMenuItem)
	m.menuViewToWatch = builder.GetObject("menu_view_to_watch").(*gtk.RadioMenuItem)
	m.menuViewSaved = builder.GetObject("menu_view_saved").(*gtk.RadioMenuItem)
	m.menuViewToDelete = builder.GetObject("menu_view_to_delete").(*gtk.RadioMenuItem)

	m.menuViewDownloads.JoinGroup(m.menuViewSubscriptions)
	m.menuViewToWatch.JoinGroup(m.menuViewSubscriptions)
	m.menuViewSaved.JoinGroup(m.menuViewSubscriptions)
	m.menuViewToDelete.JoinGroup(m.menuViewSubscriptions)
	m.menuViewSubscriptions.SetActive(true)

	m.SetupEvents()

	return nil
}

// SetupEvents sets up the menu events
func (m *menuBar) SetupEvents() {
	_ = m.menuFileQuit.Connect("activate", func() {
		gtk.MainQuit()
	})
	_ = m.menuHelpAbout.Connect("activate", func() {
	})

	_ = m.menuViewSubscriptions.Connect("activate", func() {
		m.parent.videoList.switchView(viewSubscriptions)
	})
	_ = m.menuViewDownloads.Connect("activate", func() {
		m.parent.videoList.switchView(viewDownloads)
	})
	_ = m.menuViewToWatch.Connect("activate", func() {
		m.parent.videoList.switchView(viewToWatch)
	})
	_ = m.menuViewSaved.Connect("activate", func() {
		m.parent.videoList.switchView(viewSaved)
	})
	_ = m.menuViewToDelete.Connect("activate", func() {
		m.parent.videoList.switchView(viewToDelete)
	})
}
