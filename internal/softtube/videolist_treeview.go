package softtube

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

// treeViewHelper : A helper class for a gtk treeViewHelper
type treeViewHelper struct {
	videoList *videoList
}

// Setup : Set up the video list
func (t *treeViewHelper) Setup() {
	t.setupEvents()
	t.setupColumns()
}

// setupEvents : Set up the list events
func (t *treeViewHelper) setupEvents() {
	_ = t.videoList.treeView.Connect("row_activated", t.videoList.rowActivated)
}

// setupColumns : Sets up the listview columns
func (t *treeViewHelper) setupColumns() {
	tw := t.videoList.treeView

	tw.AppendColumn(t.createImageColumn("Image"))
	tw.AppendColumn(t.createTextColumn("Channel name", listStoreColumnChannelName, 200, 300))
	tw.AppendColumn(t.createTextColumn("Date", listStoreColumnDate, 120, 300))
	tw.AppendColumn(t.createTextColumn("Duration", listStoreColumnDuration, 90, 300))
	tw.AppendColumn(t.createTextColumn("Title", listStoreColumnTitle, 0, 600))
	tw.AppendColumn(t.createProgressColumn("Progress"))
}

// createTextColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *treeViewHelper) createTextColumn(title string, id listStoreColumnType, width int, weight int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}
	_ = cellRenderer.SetProperty("weight", weight)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", int(id))
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.AddAttribute(cellRenderer, "background", int(listStoreColumnBackground))
	column.AddAttribute(cellRenderer, "foreground", int(listStoreColumnForeground))
	if width == 0 {
		column.SetExpand(true)
	} else {
		column.SetFixedWidth(width)
	}

	return column
}

// createImageColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *treeViewHelper) createImageColumn(title string) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererPixbufNew()
	if err != nil {
		log.Fatal("Unable to create pixbuf cell renderer:", err)
	}
	// cellRenderer.SetProperty("weight", weight)
	// cellRenderer.SetVisible(true)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "pixbuf", int(listStoreColumnImage))
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetFixedWidth(160)
	column.SetVisible(true)
	column.SetExpand(false)

	return column
}

// createProgressColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *treeViewHelper) createProgressColumn(title string) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererProgressNew()
	if err != nil {
		log.Fatal("Unable to create progress cell renderer:", err)
	}
	// cellRenderer.SetVisible(true)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", int(listStoreColumnProgressText))
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetFixedWidth(90)
	column.SetVisible(true)
	column.SetExpand(false)
	column.AddAttribute(cellRenderer, "value", int(listStoreColumnProgress))

	return column
}
