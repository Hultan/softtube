package softtube

import (
	"github.com/gotk3/gotk3/gtk"
)

// searchBar is the SoftTube search bar
type searchBar struct {
	parent      *SoftTube
	clearButton *gtk.Button
	searchEntry *gtk.Entry
}

// Init initializes the searchBar
func (s *searchBar) Init() error {
	s.clearButton = GetObject[*gtk.Button]("clear_search_button")
	s.searchEntry = GetObject[*gtk.Entry]("search_entry")

	s.SetupEvents()

	return nil
}

// SetupEvents sets up the toolbar events
func (s *searchBar) SetupEvents() {
	_ = s.clearButton.Connect(
		"clicked", func() {
			s.Clear()
		},
	)
	_ = s.searchEntry.Connect(
		"activate", func() {
			text, _ := s.parent.searchBar.searchEntry.GetText()
			s.parent.videoList.Search(text)
		},
	)
}

// Clear clears the previous search
func (s *searchBar) Clear() {
	s.parent.searchBar.searchEntry.SetText("")
	s.parent.videoList.Refresh("")
}
