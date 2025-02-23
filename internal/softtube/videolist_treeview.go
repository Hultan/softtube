package softtube

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

// treeViewHelper : A helper class for a gtk treeViewHelper
type treeViewHelper struct {
	videoList *videoList
}

// Setup : Set up the video list
func (t *treeViewHelper) Setup() {
	t.setupMultiSelection()
	t.setupEvents()
	t.setupColumns()
}

func (t *treeViewHelper) setupMultiSelection() {
	// Enable multiple selection
	selection, err := t.videoList.treeView.GetSelection()
	if err != nil {
		log.Fatal("Unable to get TreeSelection:", err)
	}
	selection.SetMode(gtk.SELECTION_MULTIPLE)
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
func (t *treeViewHelper) createTextColumn(
	title string, id listStoreColumnType, width, weight int,
) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		t.videoList.parent.Logger.Error.Println("Unable to create text cell renderer:", err)
		panic(err)
	}

	// Font weight and size
	_ = cellRenderer.SetProperty("weight", weight)
	_ = cellRenderer.SetProperty("size", 12500)
	_ = cellRenderer.SetProperty("ellipsize", pango.ELLIPSIZE_END)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", int(id))
	if err != nil {
		t.videoList.parent.Logger.Error.Println("Unable to create cell column:", err)
		panic(err)
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
		t.videoList.parent.Logger.Error.Println("Unable to create pixbuf cell renderer:", err)
		panic(err)
	}
	// cellRenderer.SetProperty("weight", weight)
	// cellRenderer.SetVisible(true)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "pixbuf", int(listStoreColumnImage))
	if err != nil {
		t.videoList.parent.Logger.Error.Println("Unable to create cell column:", err)
		panic(err)
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
		t.videoList.parent.Logger.Error.Println("Unable to create progress cell renderer:", err)
		panic(err)
	}
	// cellRenderer.SetVisible(true)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", int(listStoreColumnProgressText))
	if err != nil {
		t.videoList.parent.Logger.Error.Println("Unable to create cell column:", err)
		panic(err)
	}
	column.SetFixedWidth(90)
	column.SetVisible(true)
	column.SetExpand(false)
	column.AddAttribute(cellRenderer, "value", int(listStoreColumnProgress))

	return column
}
