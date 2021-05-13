package main

import (
	"github.com/gotk3/gotk3/gtk"
)

// SearchBar : The SoftTube search bar
type SearchBar struct {
	Parent      *SoftTube
	ClearButton *gtk.Button
	SearchEntry *gtk.Entry
}

// Load : Loads the toolbar
func (s *SearchBar) Load(builder *SoftBuilder) error {
	clearButton := builder.getObject("clear_search_button").(*gtk.Button)
	s.ClearButton = clearButton

	searchEntry := builder.getObject("search_entry").(*gtk.Entry)
	s.SearchEntry = searchEntry

	return nil
}

// SetupEvents : Setup the toolbar events
func (s *SearchBar) SetupEvents() {
	_,_ = s.ClearButton.Connect("clicked", func() {
		s.Parent.SearchBar.SearchEntry.SetText("")
		s.Parent.VideoList.Refresh("")
	})
	_,_ = s.SearchEntry.Connect("activate", func() {
		text, _ := s.Parent.SearchBar.SearchEntry.GetText()
		s.Parent.VideoList.Search(text)
	})
}
