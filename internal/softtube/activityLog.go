package softtube

import (
	"unicode/utf8"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softtube/internal/softtube.database"
)

// activityLog handles the SoftTube activity log
type activityLog struct {
	parent      *SoftTube
	treeView    *gtk.TreeView
	listStore   *gtk.ListStore
	imageBuffer [6]*gdk.Pixbuf // Images for download, play, delete, set watched/unwatched and error
}

// Init initiates the log
func (al *activityLog) Init() error {
	al.treeView = GetObject[*gtk.TreeView]("log_treeview")

	store, err := gtk.ListStoreNew(gdk.PixbufGetType(), glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		return err
	}
	al.listStore = store
	al.setupColumns()
	al.loadResources()
	al.FillLog()

	return nil
}

// FillLog fills the log with the last n logs
func (al *activityLog) FillLog() {
	logs := al.getLogs()

	al.treeView.SetModel(nil)
	for _, logItem := range logs {
		al.insertLog(logItem.Type, logItem.Message, false)
	}
	al.treeView.SetModel(al.listStore)
}

// AddLog adds a log to the GUI log
func (al *activityLog) AddLog(logType database.LogType, logMessage string) {
	// Insert into the gui log
	al.treeView.SetModel(nil)
	al.insertLog(logType, logMessage, true)
	al.treeView.SetModel(al.listStore)
}

func (al *activityLog) insertLog(logType database.LogType, logMessage string, first bool) {
	col := al.getColor(logType)
	img := al.imageBuffer[logType]
	var iter *gtk.TreeIter

	if first {
		iter = al.listStore.InsertAfter(nil)
	} else {
		iter = al.listStore.InsertBefore(nil)
	}
	_ = al.listStore.Set(iter, []int{0, 1, 2}, []interface{}{img, al.shortenString(logMessage), col})
}

func (al *activityLog) shortenString(text string) string {
	if utf8.RuneCountInString(text) > 50 {
		r := []rune(text)
		return string(r[:47]) + "..."
	}
	//if len(text) > 50 {
	//	return text[:47] + "..."
	//}
	return text
}

func (al *activityLog) getLogs() []database.Log {
	logs, err := al.parent.DB.Log.GetLatest()
	if err != nil {
		al.parent.Logger.Error.Println("Failed to load logs!")
		al.parent.Logger.Error.Println(err)
		return nil
	}
	return logs
}

func (al *activityLog) setupColumns() {
	imageRenderer, _ := gtk.CellRendererPixbufNew()
	imageColumn, _ := gtk.TreeViewColumnNew()
	imageColumn.SetExpand(false)
	imageColumn.SetFixedWidth(48)
	imageColumn.SetVisible(true)
	imageColumn.SetTitle("Type")
	imageColumn.PackStart(imageRenderer, true)
	imageColumn.AddAttribute(imageRenderer, "pixbuf", 0)
	al.treeView.AppendColumn(imageColumn)

	logTextRenderer, _ := gtk.CellRendererTextNew()
	logTextColumn, _ := gtk.TreeViewColumnNew()
	logTextColumn.SetExpand(false)
	logTextColumn.SetVisible(true)
	logTextColumn.SetTitle("Log")
	logTextColumn.PackStart(logTextRenderer, true)
	logTextColumn.AddAttribute(logTextRenderer, "text", 1)
	logTextColumn.AddAttribute(logTextRenderer, "background", 2)
	al.treeView.AppendColumn(logTextColumn)
}

func (al *activityLog) loadResources() {
	al.imageBuffer[0] = al.createPixbuf(downloadIcon)
	al.imageBuffer[1] = al.createPixbuf(playIcon)
	al.imageBuffer[2] = al.createPixbuf(deleteIcon)
	al.imageBuffer[3] = al.createPixbuf(setWatchedIcon)
	al.imageBuffer[4] = al.createPixbuf(setUnwatchedIcon)
	al.imageBuffer[5] = al.createPixbuf(errorIcon)
}

func (al *activityLog) createPixbuf(bytes []byte) *gdk.Pixbuf {
	pic, err := gdk.PixbufNewFromBytesOnly(bytes)
	if err != nil {
		al.parent.Logger.Error.Println(err)
	}
	return pic
}

func (al *activityLog) getColor(logType database.LogType) string {
	col := constColorNotDownloaded

	switch logType {
	case constLogDownload:
		col = constColorDownloaded
		break
	case constLogPlay:
		col = constColorWatched
		break
	case constLogDelete:
		col = constColorDeleted
		break
	case constLogSetWatched:
		col = constColorNotDownloaded
		break
	case constLogSetUnwatched:
		col = constColorNotDownloaded
		break
	case constLogError:
		col = constColorWarning
		break
	default:
		break
	}

	return string(col)
}
