package softtube

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

// treeViewHelper : A helper class for a gtk treeViewHelper
type treeViewHelper struct {
	videoList *videoList
}

// Setup : Setup the video list
func (t *treeViewHelper) Setup() {
	t.setupEvents()
	t.setupColumns()
}

// setupEvents : Setup the list events
func (t *treeViewHelper) setupEvents() {
	// Send in the videolist as a user data parameter to the event
	_ = t.videoList.treeView.Connect("row_activated", t.videoList.rowActivated)
}

// setupColumns : Sets up the listview columns
func (t *treeViewHelper) setupColumns() {
	t.videoList.treeView.AppendColumn(t.CreateImageColumn("Image"))
	t.videoList.treeView.AppendColumn(t.CreateTextColumn("Channel name", listStoreColumnChannelName, 200, 300))
	t.videoList.treeView.AppendColumn(t.CreateTextColumn("Date", listStoreColumnDate, 90, 300))
	t.videoList.treeView.AppendColumn(t.CreateTextColumn("Title", listStoreColumnTitle, 0, 600))
	t.videoList.treeView.AppendColumn(t.CreateTextColumn("Duration", listStoreColumnDuration, 90, 300))
	t.videoList.treeView.AppendColumn(t.CreateProgressColumn("Progress"))
}

// CreateTextColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *treeViewHelper) CreateTextColumn(title string, id int, width int, weight int) *gtk.TreeViewColumn {
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
func (t *treeViewHelper) CreateImageColumn(title string) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererPixbufNew()
	if err != nil {
		log.Fatal("Unable to create pixbuf cell renderer:", err)
	}
	// cellRenderer.SetProperty("weight", weight)
	// cellRenderer.SetVisible(true)

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
func (t *treeViewHelper) CreateProgressColumn(title string) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererProgressNew()
	if err != nil {
		log.Fatal("Unable to create progress cell renderer:", err)
	}
	// cellRenderer.SetVisible(true)

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
