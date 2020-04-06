package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
)

// GtkHelper : A helper class for GTK
type GtkHelper struct {
}

// GetGladePath : Get the path to the glade external resource file
func (g *GtkHelper) GetGladePath() (string, error) {
	// Get directory from where the program is launched
	basePath := filepath.Dir(os.Args[0])

	// Check main path, works most times
	gladePath := path.Join(basePath, "resources/main.glade")
	if _, err := os.Stat(gladePath); err == nil {
		return gladePath, nil
	}
	// Check secondary path, for debug mode (when run from VS Code)
	gladePath = path.Join(basePath, "../resources/main.glade")
	if _, err := os.Stat(gladePath); err == nil {
		return gladePath, nil
	}

	return "", errors.New("Glade file is missing (%s)")
}

// GetWindow : Gets a gtk.Window from the builder
func (g *GtkHelper) GetWindow(builder *gtk.Builder, name string) (*gtk.Window, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		return nil, err
	}
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}

	return nil, errors.New("not a gtk window")
}

// GetMenuItem : Gets a menuitem from the builder
func (g *GtkHelper) GetMenuItem(builder *gtk.Builder, name string) (*gtk.MenuItem, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if menuItem, ok := obj.(*gtk.MenuItem); ok {
		return menuItem, nil
	}

	return nil, errors.New("not a gtk menu item")
}

// GetButton : Gets a button from the builder
func (g *GtkHelper) GetButton(builder *gtk.Builder, name string) (*gtk.Button, error) {
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

// GetEntry : Gets an Entry from the builder
func (g *GtkHelper) GetEntry(builder *gtk.Builder, name string) (*gtk.Entry, error) {
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

// GetLabel : Gets a label from the builder
func (g *GtkHelper) GetLabel(builder *gtk.Builder, name string) (*gtk.Label, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if label, ok := obj.(*gtk.Label); ok {
		return label, nil
	}

	return nil, errors.New("not a gtk label")
}

// GetToolButton : Gets a ToolButton from the builder
func (g *GtkHelper) GetToolButton(builder *gtk.Builder, name string) (*gtk.ToolButton, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.ToolButton); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk tool button")
}

// GetToggleToolButton : Gets a ToggleToolButton from the builder
func (g *GtkHelper) GetToggleToolButton(builder *gtk.Builder, name string) (*gtk.ToggleToolButton, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.ToggleToolButton); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk toggle tool button")
}

// GetTreeView : Gets a treeview from the builder
func (g *GtkHelper) GetTreeView(builder *gtk.Builder, name string) (*gtk.TreeView, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.TreeView); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk tree view")
}

// GetScrolledWindow : Gets a ScrolledWindow from the builder
func (g *GtkHelper) GetScrolledWindow(builder *gtk.Builder, name string) (*gtk.ScrolledWindow, error) {
	obj, err := builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.ScrolledWindow); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk scrolled window")
}
