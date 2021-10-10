package softtube

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
)

// searchBar : The SoftTube search bar
type searchBar struct {
	parent      *SoftTube
	clearButton *gtk.Button
	searchEntry *gtk.Entry
}

// Load : Loads the toolbar
func (s *searchBar) Load(builder *framework.GtkBuilder) error {
	clearButton := builder.GetObject("clear_search_button").(*gtk.Button)
	s.clearButton = clearButton

	searchEntry := builder.GetObject("search_entry").(*gtk.Entry)
	s.searchEntry = searchEntry
	s.SetupEvents()

	return nil
}

// SetupEvents : Set up the toolbar events
func (s *searchBar) SetupEvents() {
	_ = s.clearButton.Connect("clicked", func() {
		s.parent.searchBar.searchEntry.SetText("")
		s.parent.videoList.Refresh("")
	})
	_ = s.searchEntry.Connect("activate", func() {
		text, _ := s.parent.searchBar.searchEntry.GetText()
		s.parent.videoList.Search(text)
	})
}
