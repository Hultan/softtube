package softtube

import (
	"fmt"
	"os/exec"
	"path"
	"syscall"

	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softtube/internal/builder"
)

// menuBar : The SoftTube menu bar
type menuBar struct {
	parent                *SoftTube
	menuFileQuit          *gtk.MenuItem
	menuHelpAbout         *gtk.MenuItem
	menuViewSubscriptions *gtk.RadioMenuItem
	menuViewDownloads     *gtk.RadioMenuItem
	menuViewToWatch       *gtk.RadioMenuItem
	menuViewSaved         *gtk.RadioMenuItem
	menuViewToDelete      *gtk.RadioMenuItem

	menuViewLog       *gtk.MenuItem
	menuViewUpdateLog *gtk.MenuItem
}

// Init initiates the menu bar
func (m *menuBar) Init(builder *builder.Builder) error {
	m.menuFileQuit = builder.GetObject("menu_file_quit").(*gtk.MenuItem)
	m.menuHelpAbout = builder.GetObject("menu_help_about").(*gtk.MenuItem)

	m.menuViewSubscriptions = builder.GetObject("menu_view_subscriptions").(*gtk.RadioMenuItem)
	m.menuViewDownloads = builder.GetObject("menu_view_downloads").(*gtk.RadioMenuItem)
	m.menuViewToWatch = builder.GetObject("menu_view_to_watch").(*gtk.RadioMenuItem)
	m.menuViewSaved = builder.GetObject("menu_view_saved").(*gtk.RadioMenuItem)
	m.menuViewToDelete = builder.GetObject("menu_view_to_delete").(*gtk.RadioMenuItem)

	m.menuViewDownloads.JoinGroup(m.menuViewSubscriptions)
	m.menuViewToWatch.JoinGroup(m.menuViewSubscriptions)
	m.menuViewSaved.JoinGroup(m.menuViewSubscriptions)
	m.menuViewToDelete.JoinGroup(m.menuViewSubscriptions)
	m.menuViewSubscriptions.SetActive(true)

	m.menuViewLog = builder.GetObject("menu_view_log").(*gtk.MenuItem)
	m.menuViewUpdateLog = builder.GetObject("menu_view_update_log").(*gtk.MenuItem)

	m.SetupEvents()

	return nil
}

// SetupEvents sets up the menu events
func (m *menuBar) SetupEvents() {
	_ = m.menuFileQuit.Connect(
		"activate", func() {
			gtk.MainQuit()
		},
	)
	_ = m.menuHelpAbout.Connect(
		"activate", func() {
		},
	)

	_ = m.menuViewSubscriptions.Connect(
		"activate", func() {
			m.parent.videoList.switchView(viewSubscriptions)
		},
	)
	_ = m.menuViewDownloads.Connect(
		"activate", func() {
			m.parent.videoList.switchView(viewDownloads)
		},
	)
	_ = m.menuViewToWatch.Connect(
		"activate", func() {
			m.parent.videoList.switchView(viewToWatch)
		},
	)
	_ = m.menuViewSaved.Connect(
		"activate", func() {
			m.parent.videoList.switchView(viewSaved)
		},
	)
	_ = m.menuViewToDelete.Connect(
		"activate", func() {
			m.parent.videoList.switchView(viewToDelete)
		},
	)
	_ = m.menuViewLog.Connect(
		"activate", func() {
			m.openLogFile("softtube.client.log")
		},
	)
	_ = m.menuViewUpdateLog.Connect(
		"activate", func() {
			m.openLogFile("softtube.update.log")
		},
	)
}

func (m *menuBar) openLogFile(logFile string) {
	// Start Video Player
	go func() {
		p := path.Join(m.parent.Config.ClientPaths.Log, logFile)
		command := fmt.Sprintf("xed '%s' &", p)
		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Pgid:    0,
		}
		// Starts a sub process (smplayer)
		// Did not get this to work, but read the following, and maybe I can get
		// this to work in the future
		// https://forum.golangbridge.org/t/starting-new-processes-with-exec-command/24956
		_ = cmd.Run()
	}()
}
