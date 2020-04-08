package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	core "github.com/hultan/softtube/softtube.core"
)

// Log : Handles the GUI log
type Log struct {
	Parent      *SoftTube
	TreeView    *gtk.TreeView
	ListStore   *gtk.ListStore
	ImageBuffer [6]*gdk.Pixbuf // Images for download, play, delete, set watched/unwatched and error
}

// Load : Loads the log
func (l *Log) Load(builder *gtk.Builder) {
	helper := new(GtkHelper)
	tree, err := helper.GetTreeView(builder, "log_treeview")
	if err != nil {
		logger.LogError(err)
		return
	}
	l.TreeView = tree

	listStore, err := gtk.ListStoreNew(gdk.PixbufGetType(), glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		logger.Log("Failed to create liststore!")
		logger.LogError(err)
		panic(err)
	}
	l.ListStore = listStore
}

// FillLog : Fills the log with the last n logs
func (l *Log) FillLog() {
	logs := l.getLogs()
	l.setupColumns()
	l.loadResources()

	l.TreeView.SetModel(nil)
	for _, logItem := range logs {
		l.insertLog(logItem.Type, logItem.Message, false)
	}
	l.TreeView.SetModel(l.ListStore)
}

// InsertLog : Adds a log to the GUI log
func (l *Log) InsertLog(logType int, logMessage string) {
	// Insert into the gui log
	l.TreeView.SetModel(nil)
	l.insertLog(logType, logMessage, true)
	l.TreeView.SetModel(l.ListStore)
}

func (l *Log) insertLog(logType int, logMessage string, first bool) {
	color := l.getColor(logType)
	image := l.ImageBuffer[logType]
	var iter *gtk.TreeIter

	if first {
		iter = l.ListStore.InsertAfter(nil)
	} else {
		iter = l.ListStore.InsertBefore(nil)
	}
	l.ListStore.Set(iter, []int{0, 1, 2}, []interface{}{image, logMessage, color})
}

func (l *Log) getLogs() []core.Log {
	logs, err := l.Parent.Database.Log.GetLatest()
	if err != nil {
		logger.Log("Failed to load logs!")
		logger.LogError(err)
		return nil
	}
	return logs
}

func (l *Log) setupColumns() {
	imageRenderer, _ := gtk.CellRendererPixbufNew()
	imageColumn, _ := gtk.TreeViewColumnNew()
	imageColumn.SetExpand(false)
	imageColumn.SetFixedWidth(48)
	imageColumn.SetVisible(true)
	imageColumn.SetTitle("Type")
	imageColumn.PackStart(imageRenderer, true)
	imageColumn.AddAttribute(imageRenderer, "pixbuf", 0)
	l.TreeView.AppendColumn(imageColumn)

	logtextRenderer, _ := gtk.CellRendererTextNew()
	logtextColumn, _ := gtk.TreeViewColumnNew()
	logtextColumn.SetExpand(false)
	logtextColumn.SetVisible(true)
	logtextColumn.SetTitle("Log")
	logtextColumn.PackStart(logtextRenderer, true)
	logtextColumn.AddAttribute(logtextRenderer, "text", 1)
	logtextColumn.AddAttribute(logtextRenderer, "background", 2)
	l.TreeView.AppendColumn(logtextColumn)
}

func (l *Log) loadResources() {
	resourcePath := path.Join(filepath.Dir(os.Args[0]), "..", "resources")
	for i := constLogDownload; i <= constLogError; i++ {
		fileName := l.getImageFileName(i)
		if fileName != "" {

			pic, err := gdk.PixbufNewFromFile(path.Join(resourcePath, fileName))
			if err != nil {
				logger.LogError(err)
			}
			l.ImageBuffer[i] = pic
		}
	}
}

func (l *Log) getImageFileName(index int) string {
	switch index {
	case 0:
		return "download.png"
	case 1:
		return "play.png"
	case 2:
		return "delete.png"
	case 3:
		return "set_watched.png"
	case 4:
		return "set_unwatched.png"
	case 5:
		return "error.png"
	default:
		return ""
	}
}

func (l *Log) getColor(logtype int) string {
	color := constColorNotDownloaded

	switch logtype {
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

	return color
}