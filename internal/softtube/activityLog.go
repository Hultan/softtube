package softtube

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
	"github.com/hultan/softtube/internal/softtube.database"
)

// activityLog : Handles the SoftTube activity log
type activityLog struct {
	parent      *SoftTube
	treeView    *gtk.TreeView
	listStore   *gtk.ListStore
	imageBuffer [6]*gdk.Pixbuf // Images for download, play, delete, set watched/unwatched and error
}

// Init : Loads the log
func (a *activityLog) Init(builder *framework.GtkBuilder) {
	tree := builder.GetObject("log_treeview").(*gtk.TreeView)
	a.treeView = tree

	store, err := gtk.ListStoreNew(gdk.PixbufGetType(), glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		a.parent.Logger.Log("Failed to create liststore!")
		a.parent.Logger.LogError(err)
		panic(err)
	}
	a.listStore = store
	a.FillLog(true)
}

// FillLog : Fills the log with the last n logs
func (a *activityLog) FillLog(init bool) {
	logs := a.getLogs()
	if init {
		a.setupColumns()
		a.loadResources()
	}

	a.treeView.SetModel(nil)
	for _, logItem := range logs {
		a.insertLog(logItem.Type, logItem.Message, false)
	}
	a.treeView.SetModel(a.listStore)
}

// AddLog : Adds a log to the GUI log
func (a *activityLog) AddLog(logType database.LogType, logMessage string) {
	// Insert into the gui log
	a.treeView.SetModel(nil)
	a.insertLog(logType, logMessage, true)
	a.treeView.SetModel(a.listStore)
}

func (a *activityLog) insertLog(logType database.LogType, logMessage string, first bool) {
	color := a.getColor(logType)
	image := a.imageBuffer[logType]
	var iter *gtk.TreeIter

	if first {
		iter = a.listStore.InsertAfter(nil)
	} else {
		iter = a.listStore.InsertBefore(nil)
	}
	_ = a.listStore.Set(iter, []int{0, 1, 2}, []interface{}{image, a.shortenString(logMessage), color})
}

func (a *activityLog) shortenString(text string) string {
	if len(text) > 50 {
		return text[:47] + "..."
	}
	return text
}

func (a *activityLog) getLogs() []database.Log {
	logs, err := a.parent.DB.Log.GetLatest()
	if err != nil {
		a.parent.Logger.Log("Failed to load logs!")
		a.parent.Logger.LogError(err)
		return nil
	}
	return logs
}

func (a *activityLog) setupColumns() {
	imageRenderer, _ := gtk.CellRendererPixbufNew()
	imageColumn, _ := gtk.TreeViewColumnNew()
	imageColumn.SetExpand(false)
	imageColumn.SetFixedWidth(48)
	imageColumn.SetVisible(true)
	imageColumn.SetTitle("Type")
	imageColumn.PackStart(imageRenderer, true)
	imageColumn.AddAttribute(imageRenderer, "pixbuf", 0)
	a.treeView.AppendColumn(imageColumn)

	logTextRenderer, _ := gtk.CellRendererTextNew()
	logTextColumn, _ := gtk.TreeViewColumnNew()
	logTextColumn.SetExpand(false)
	logTextColumn.SetVisible(true)
	logTextColumn.SetTitle("Log")
	logTextColumn.PackStart(logTextRenderer, true)
	logTextColumn.AddAttribute(logTextRenderer, "text", 1)
	logTextColumn.AddAttribute(logTextRenderer, "background", 2)
	a.treeView.AppendColumn(logTextColumn)
}

func (a *activityLog) loadResources() {
	a.imageBuffer[0] = a.createPixbuf(downloadIcon)
	a.imageBuffer[1] = a.createPixbuf(playIcon)
	a.imageBuffer[2] = a.createPixbuf(deleteIcon)
	a.imageBuffer[3] = a.createPixbuf(setWatchedIcon)
	a.imageBuffer[4] = a.createPixbuf(setUnwatchedIcon)
	a.imageBuffer[5] = a.createPixbuf(errorIcon)
}

func (a *activityLog) createPixbuf(bytes []byte) *gdk.Pixbuf {
	pic, err := gdk.PixbufNewFromBytesOnly(bytes)
	if err != nil {
		a.parent.Logger.LogError(err)
	}
	return pic
}

func (a *activityLog) getColor(logType database.LogType) string {
	color := constColorNotDownloaded

	switch logType {
	case constLogDownload:
		color = constColorDownloaded
		break
	case constLogPlay:
		color = constColorWatched
		break
	case constLogDelete:
		color = constColorDeleted
		break
	case constLogSetWatched:
		color = constColorNotDownloaded
		break
	case constLogSetUnwatched:
		color = constColorNotDownloaded
		break
	case constLogError:
		color = constColorWarning
		break
	default:
		break
	}

	return string(color)
}
