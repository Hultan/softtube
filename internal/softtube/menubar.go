package softtube

import (
	"fmt"
	"os/exec"
	"path"
	"syscall"

	"github.com/gotk3/gotk3/gtk"
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

	menuViewOpenSoftTube *gtk.MenuItem

	menuViewLog       *gtk.MenuItem
	menuViewUpdateLog *gtk.MenuItem
}

// Init initiates the menu bar
func (m *menuBar) Init() error {
	m.menuFileQuit = GetObject[*gtk.MenuItem]("menu_file_quit")
	m.menuHelpAbout = GetObject[*gtk.MenuItem]("menu_help_about")
	m.menuViewSubscriptions = GetObject[*gtk.RadioMenuItem]("menu_view_subscriptions")
	m.menuViewDownloads = GetObject[*gtk.RadioMenuItem]("menu_view_downloads")
	m.menuViewToWatch = GetObject[*gtk.RadioMenuItem]("menu_view_to_watch")
	m.menuViewSaved = GetObject[*gtk.RadioMenuItem]("menu_view_saved")
	m.menuViewToDelete = GetObject[*gtk.RadioMenuItem]("menu_view_to_delete")

	m.menuViewDownloads.JoinGroup(m.menuViewSubscriptions)
	m.menuViewToWatch.JoinGroup(m.menuViewSubscriptions)
	m.menuViewSaved.JoinGroup(m.menuViewSubscriptions)
	m.menuViewToDelete.JoinGroup(m.menuViewSubscriptions)
	m.menuViewSubscriptions.SetActive(true)

	m.menuViewOpenSoftTube = GetObject[*gtk.MenuItem]("menu_view_open_softtube")
	m.menuViewLog = GetObject[*gtk.MenuItem]("menu_view_log")
	m.menuViewUpdateLog = GetObject[*gtk.MenuItem]("menu_view_update_log")

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

	_ = m.menuViewOpenSoftTube.Connect(
		"activate", func() {
			m.openSoftTubeFolder()
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
	go func() {
		p := path.Join(m.parent.Config.ClientPaths.Log, logFile)
		command := fmt.Sprintf("xed '%s' &", p)
		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Pgid:    0,
		}
		// Starts a subprocess (smplayer)
		// Did not get this to work, but read the following, and maybe I can get
		// this to work in the future
		// https://forum.golangbridge.org/t/starting-new-processes-with-exec-command/24956
		err := cmd.Run()
		if err != nil {
			m.parent.Logger.Error.Println(err)
		}
	}()
}

func (m *menuBar) openSoftTubeFolder() {
	go func() {
		command := "nemo '/softtube' &"
		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Pgid:    0,
		}
		// Starts a subprocess (smplayer)
		// Did not get this to work, but read the following, and maybe I can get
		// this to work in the future
		// https://forum.golangbridge.org/t/starting-new-processes-with-exec-command/24956
		err := cmd.Run()
		if err != nil {
			m.parent.Logger.Error.Println(err)
		}
	}()
}
