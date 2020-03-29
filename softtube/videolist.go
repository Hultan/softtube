package main

import (
	"errors"
	"fmt"
	"log"
	"path"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	core "github.com/hultan/softtube/softtube.core"
)

// VideoList : The SoftTube video list
type VideoList struct {
	List *gtk.TreeView
}

// Load : Loads the toolbar
func (v *VideoList) Load(builder *gtk.Builder) error {
	list, err := getTreeView(builder, "video_treeview")
	if err != nil {
		return err
	}
	v.List = list

	return nil
}

// SetupEvents : Setup the list events
func (v *VideoList) SetupEvents() {
	// v.List.Connect("button-press-event", func() {
	// 	gtk.MainQuit()
	// })
}

// Fill : Fills the video list
func (v *VideoList) Fill(db *core.Database) {
	videos, err := db.Videos.GetVideos()
	if err != nil {
		logger.LogError(err)
		panic(err)
	}
	fmt.Println("Videos loaded!", len(videos))

	v.List.SetModel(nil)
	listStore, err := gtk.ListStoreNew(gdk.PixbufGetType(), glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_INT64, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		logger.Log("Failed to create liststore!")
		logger.LogError(err)
		panic(err)
	}

	for i := 0; i < len(videos); i++ {
		video := videos[i]
		progress := 0
		progressText := ""
		var color string = ""

		thumbnail := getThumbnail(video.ID)

		iter := listStore.Append()
		err = listStore.Set(iter, []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
			[]interface{}{thumbnail,
				video.SubscriptionName,
				video.Added.Format(constDateLayout),
				video.Title,
				progress,
				color,
				video.ID,
				video.Duration,
				progressText})

		if err != nil {
			logger.Log("Failed to add row!")
			logger.LogError(err)
		}
	}

	v.List.SetModel(listStore)
}

// SetupColumns : Sets up the listview columns
func (v VideoList) SetupColumns() {
	v.List.AppendColumn(createImageColumn("Image"))
	v.List.AppendColumn(createTextColumn("Channel name", liststoreColumnChannelName, 200, 200))
	v.List.AppendColumn(createTextColumn("Date", liststoreColumnDate, 90, 200))
	v.List.AppendColumn(createTextColumn("Title", liststoreColumnTitle, 0, 600))
	v.List.AppendColumn(createTextColumn("Duration", liststoreColumnDuration, 90, 200))
	v.List.AppendColumn(createProgressColumn("Progress"))
}

func getThumbnailPath(videoID string) string {
	fmt.Println(config.ClientPaths.Thumbnails)
	fmt.Println(path.Join(config.ClientPaths.Thumbnails, fmt.Sprintf("%s.jpg", videoID)))
	return "/" + path.Join(config.ClientPaths.Thumbnails, fmt.Sprintf("%s.jpg", videoID))
}

func getThumbnail(videoID string) *gdk.Pixbuf {
	path := getThumbnailPath(videoID)

	thumbnail, err := gdk.PixbufNewFromFile(path)
	if err != nil {
		msg := fmt.Sprintf("Failed to load thumnail (%s)!", path)
		logger.Log(msg)
		thumbnail = nil
	} else {
		thumbnail, err = thumbnail.ScaleSimple(160, 90, gdk.INTERP_BILINEAR)
		if err != nil {
			msg := fmt.Sprintf("Failed to scale thumnail (%s)!", path)
			logger.Log(msg)
			thumbnail = nil
		}
	}

	return thumbnail
}

// Add a column to the tree view (during the initialization of the tree view)
func createTextColumn(title string, id int, width int, weight int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}
	cellRenderer.SetProperty("weight", weight)
	//cellRenderer.ellipsize = Pango.EllipsizeMode.END

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.AddAttribute(cellRenderer, "background", liststoreColumnBackground)
	if width == 0 {
		column.SetExpand(true)
	} else {
		column.SetFixedWidth(width)
	}

	return column
}

// Add a column to the tree view (during the initialization of the tree view)
func createImageColumn(title string) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererPixbufNew()
	if err != nil {
		log.Fatal("Unable to create pixbuf cell renderer:", err)
	}
	//cellRenderer.SetProperty("weight", weight)
	//cellRenderer.SetVisible(true)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "pixbuf", liststoreColumnImage)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetFixedWidth(160)
	column.SetVisible(true)
	column.SetExpand(false)

	return column
}

// Add a column to the tree view (during the initialization of the tree view)
func createProgressColumn(title string) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererProgressNew()
	if err != nil {
		log.Fatal("Unable to create progress cell renderer:", err)
	}
	//cellRenderer.SetVisible(true)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", liststoreColumnProgressText)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetFixedWidth(90)
	column.SetVisible(true)
	column.SetExpand(false)
	column.AddAttribute(cellRenderer, "value", liststoreColumnProgress)

	return column
}

func getTreeView(builder *gtk.Builder, name string) (*gtk.TreeView, error) {
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
