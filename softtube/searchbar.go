package main

import (
	"errors"

	"github.com/gotk3/gotk3/gtk"
)

// SearchBar : The SoftTube search bar
type SearchBar struct {
	Parent      *SoftTube
	ClearButton *gtk.Button
	SearchEntry *gtk.Entry
}

// Load : Loads the toolbar
func (s *SearchBar) Load(builder *gtk.Builder) error {
	clearButton, err := getButton(builder, "clear_search_button")
	if err != nil {
		return err
	}
	s.ClearButton = clearButton

	searchEntry, err := getEntry(builder, "search_entry")
	if err != nil {
		return err
	}
	s.SearchEntry = searchEntry

	return nil
}

// SetupEvents : Setup the toolbar events
func (s *SearchBar) SetupEvents() {
	s.ClearButton.Connect("clicked", func() {
		s.Parent.SearchBar.SearchEntry.SetText("")
		s.Parent.VideoList.Refresh("")
	})
	s.SearchEntry.Connect("activate", func() {
		text, _ := s.Parent.SearchBar.SearchEntry.GetText()
		s.Parent.VideoList.Search(text)
	})
}

func getButton(builder *gtk.Builder, name string) (*gtk.Button, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if button, ok := obj.(*gtk.Button); ok {
		return button, nil
	}

	return nil, errors.New("not a gtk button")
}

func getEntry(builder *gtk.Builder, name string) (*gtk.Entry, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if entry, ok := obj.(*gtk.Entry); ok {
		return entry, nil
	}

	return nil, errors.New("not a gtk entry")
}
