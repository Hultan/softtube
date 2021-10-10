package softtube

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

// TreeviewHelper : A helper class for a gtk treeview
type TreeviewHelper struct {
}

// CreateTextColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *TreeviewHelper) CreateTextColumn(title string, id int, width int, weight int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}
	_ = cellRenderer.SetProperty("weight", weight)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.AddAttribute(cellRenderer, "background", listStoreColumnBackground)
	column.AddAttribute(cellRenderer, "foreground", listStoreColumnForeground)
	if width == 0 {
		column.SetExpand(true)
	} else {
		column.SetFixedWidth(width)
	}

	return column
}

// CreateImageColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *TreeviewHelper) CreateImageColumn(title string) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererPixbufNew()
	if err != nil {
		log.Fatal("Unable to create pixbuf cell renderer:", err)
	}
	//cellRenderer.SetProperty("weight", weight)
	//cellRenderer.SetVisible(true)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "pixbuf", listStoreColumnImage)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetFixedWidth(160)
	column.SetVisible(true)
	column.SetExpand(false)

	return column
}

// CreateProgressColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *TreeviewHelper) CreateProgressColumn(title string) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererProgressNew()
	if err != nil {
		log.Fatal("Unable to create progress cell renderer:", err)
	}
	//cellRenderer.SetVisible(true)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", listStoreColumnProgressText)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetFixedWidth(90)
	column.SetVisible(true)
	column.SetExpand(false)
	column.AddAttribute(cellRenderer, "value", listStoreColumnProgress)

	return column
}
