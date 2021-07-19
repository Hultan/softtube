package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam/framework"
)

// SearchBar : The SoftTube search bar
type SearchBar struct {
	Parent      *SoftTube
	ClearButton *gtk.Button
	SearchEntry *gtk.Entry
}

// Load : Loads the toolbar
func (s *SearchBar) Load(builder *framework.GtkBuilder) error {
	clearButton := builder.GetObject("clear_search_button").(*gtk.Button)
	s.ClearButton = clearButton

	searchEntry := builder.GetObject("search_entry").(*gtk.Entry)
	s.SearchEntry = searchEntry

	return nil
}

// SetupEvents : Setup the toolbar events
func (s *SearchBar) SetupEvents() {
	_ = s.ClearButton.Connect("clicked", func() {
		s.Parent.SearchBar.SearchEntry.SetText("")
		s.Parent.VideoList.Refresh("")
	})
	_ = s.SearchEntry.Connect("activate", func() {
		text, _ := s.Parent.SearchBar.SearchEntry.GetText()
		s.Parent.VideoList.Search(text)
	})
}
