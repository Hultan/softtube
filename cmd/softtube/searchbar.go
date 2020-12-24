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
func (s *SearchBar) Load(helper *GtkHelper) error {
	clearButton, err := helper.GetButton("clear_search_button")
	if err != nil {
		return err
	}
	s.ClearButton = clearButton

	searchEntry, err := helper.GetEntry("search_entry")
	if err != nil {
		return err
	}
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
