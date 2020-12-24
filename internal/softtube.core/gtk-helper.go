package core

import (
	"errors"
	"os"
	"path"

	"github.com/gotk3/gotk3/gtk"
)

// GtkHelper : A helper class for GTK
type GtkHelper struct {
	builder *gtk.Builder
}

func GtkHelperNew(builder *gtk.Builder) *GtkHelper {
	helper := new(GtkHelper)
	helper.builder = builder
	return helper
}

// exists returns whether the given file or directory exists
func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func GtkHelperNewFromFile(fileName string) *GtkHelper {
	helper := new(GtkHelper)

	resources := new(Resources)
	exePath := resources.GetExecutablePath()
	gladePath := path.Join(exePath, fileName)
	if !exists(gladePath) {
		gladePath = path.Join(exePath, "assets", fileName)
		if !exists(gladePath) {
			gladePath = path.Join(exePath, "../assets", fileName)
			if !exists(gladePath) {
				return nil
			}
		}
	}

	builder, err := gtk.BuilderNewFromFile(gladePath)
	if err != nil {
		panic(err)
	}

	helper.builder = builder
	return helper
}

func (g *GtkHelper) SetBuilder(builder *gtk.Builder) {
	g.builder = builder
}

// GetGladePath : Get the path to the glade external resource file
func (g *GtkHelper) GetGladePath(fileName string) (string, error) {
	// Check main path, works most times
	resources := new(Resources)
	if fileName == "" {
		fileName = "main.glade"
	}

	gladePath := resources.GetResourcePath(fileName)
	if gladePath == "" {
		return "", errors.New("glade file is missing")
	}
	if _, err := os.Stat(gladePath); err == nil {
		return gladePath, nil
	}

	return "", errors.New("glade file is missing")
}

// GetWindow : Gets a gtk.Window from the builder
func (g *GtkHelper) GetWindow(name string) (*gtk.Window, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		return nil, err
	}
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}

	return nil, errors.New("not a gtk window")
}

// GetDrawingArea : Gets a gtk.DrawingArea from the builder
func (g *GtkHelper) GetDrawingArea(name string) (*gtk.DrawingArea, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		return nil, err
	}
	if win, ok := obj.(*gtk.DrawingArea); ok {
		return win, nil
	}

	return nil, errors.New("not a gtk drawing area")
}

// GetMenu : Gets a popup menu from the builder
func (g *GtkHelper) GetMenu(name string) (*gtk.Menu, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if menu, ok := obj.(*gtk.Menu); ok {
		return menu, nil
	}

	return nil, errors.New("not a gtk menu")
}

// GetButton : Gets a button from the builder
func (g *GtkHelper) GetButton(name string) (*gtk.Button, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if button, ok := obj.(*gtk.Button); ok {
		return button, nil
	}

	return nil, errors.New("not a gtk button")
}

// GetMenuItem : Gets a menuitem from the builder
func (g *GtkHelper) GetMenuItem(name string) (*gtk.MenuItem, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if menuItem, ok := obj.(*gtk.MenuItem); ok {
		return menuItem, nil
	}

	return nil, errors.New("not a gtk menu item")
}

// GetEventBox : Gets an event box from the builder
func (g *GtkHelper) GetEventBox(name string) (*gtk.EventBox, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if button, ok := obj.(*gtk.EventBox); ok {
		return button, nil
	}

	return nil, errors.New("not a gtk event box")
}

// GetEntry : Gets an Entry from the builder
func (g *GtkHelper) GetEntry(name string) (*gtk.Entry, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if entry, ok := obj.(*gtk.Entry); ok {
		return entry, nil
	}

	return nil, errors.New("not a gtk entry")
}

// GetCalendar : Gets a Calendar from the builder
func (g *GtkHelper) GetCalendar(name string) (*gtk.Calendar, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if entry, ok := obj.(*gtk.Calendar); ok {
		return entry, nil
	}

	return nil, errors.New("not a gtk calendar")
}

// GetLabel : Gets a label from the builder
func (g *GtkHelper) GetLabel(name string) (*gtk.Label, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if label, ok := obj.(*gtk.Label); ok {
		return label, nil
	}

	return nil, errors.New("not a gtk label")
}

// GetComboBox : Gets a ComboBox from the builder
func (g *GtkHelper) GetComboBox(name string) (*gtk.ComboBox, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if label, ok := obj.(*gtk.ComboBox); ok {
		return label, nil
	}

	return nil, errors.New("not a gtk combobox")
}

// GetToolButton : Gets a ToolButton from the builder
func (g *GtkHelper) GetToolButton(name string) (*gtk.ToolButton, error) {
	obj, err := g.builder.GetObject(name)
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
func (g *GtkHelper) GetToggleToolButton(name string) (*gtk.ToggleToolButton, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.ToggleToolButton); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk toggle tool button")
}

// GetTreeView : Gets a tree view from the builder
func (g *GtkHelper) GetTreeView(name string) (*gtk.TreeView, error) {
	obj, err := g.builder.GetObject(name)
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
func (g *GtkHelper) GetScrolledWindow(name string) (*gtk.ScrolledWindow, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.ScrolledWindow); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk scrolled window")
}

// GetRadioButton : Gets a RadioButton from the builder
func (g *GtkHelper) GetRadioButton(name string) (*gtk.RadioButton, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.RadioButton); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk radio button")
}

// GetCheckButton : Gets a CheckButton from the builder
func (g *GtkHelper) GetCheckButton(name string) (*gtk.CheckButton, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if tool, ok := obj.(*gtk.CheckButton); ok {
		return tool, nil
	}

	return nil, errors.New("not a gtk check button")
}

// GetBox : Gets a Box from the builder
func (g *GtkHelper) GetBox(name string) (*gtk.Box, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if box, ok := obj.(*gtk.Box); ok {
		return box, nil
	}

	return nil, errors.New("not a gtk box")
}

// GetApplicationWindow : Gets a gtk.ApplicationWindow from the builder
func (g *GtkHelper) GetApplicationWindow(name string) (*gtk.ApplicationWindow, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		return nil, err
	}
	if win, ok := obj.(*gtk.ApplicationWindow); ok {
		return win, nil
	}

	return nil, errors.New("not a gtk application window")
}

// GetImage : Gets a gtk.Image from the builder
func (g *GtkHelper) GetImage(name string) (*gtk.Image, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		return nil, err
	}
	if win, ok := obj.(*gtk.Image); ok {
		return win, nil
	}

	return nil, errors.New("not a gtk image")
}

// GetSpinButton : Gets a gtk.SpinButton from the builder
func (g *GtkHelper) GetSpinButton(name string) (*gtk.SpinButton, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		return nil, err
	}
	if win, ok := obj.(*gtk.SpinButton); ok {
		return win, nil
	}

	return nil, errors.New("not a gtk spin button")
}

// GetFixed : Gets a fixed from the builder
func (g *GtkHelper) GetFixed(name string) (*gtk.Fixed, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if fixed, ok := obj.(*gtk.Fixed); ok {
		return fixed, nil
	}

	return nil, errors.New("not a gtk fixed")
}

// GetStatusBar : Gets a status bar from the builder
func (g *GtkHelper) GetStatusBar(name string) (*gtk.Statusbar, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if bar, ok := obj.(*gtk.Statusbar); ok {
		return bar, nil
	}

	return nil, errors.New("not a gtk status bar")
}

// GetDialog : Gets a dialog from the builder
func (g *GtkHelper) GetDialog(name string) (*gtk.Dialog, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if dialog, ok := obj.(*gtk.Dialog); ok {
		return dialog, nil
	}

	return nil, errors.New("not a gtk dialog")
}

// GetAboutDialog : Gets an about dialog from the builder
func (g *GtkHelper) GetAboutDialog(name string) (*gtk.AboutDialog, error) {
	obj, err := g.builder.GetObject(name)
	if err != nil {
		// object not found
		return nil, err
	}
	if aboutDialog, ok := obj.(*gtk.AboutDialog); ok {
		return aboutDialog, nil
	}

	return nil, errors.New("not a gtk about dialog")
}
