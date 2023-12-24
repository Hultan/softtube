package softtube

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softtube/internal/builder"
)

// searchBar : The SoftTube search bar
type searchBar struct {
	parent      *SoftTube
	clearButton *gtk.Button
	searchEntry *gtk.Entry
}

// Init initializes the searchBar
func (s *searchBar) Init(builder *builder.Builder) error {
	clearButton := builder.GetObject("clear_search_button").(*gtk.Button)
	s.clearButton = clearButton

	searchEntry := builder.GetObject("search_entry").(*gtk.Entry)
	s.searchEntry = searchEntry

	s.SetupEvents()

	return nil
}

// SetupEvents : Set up the toolbar events
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
