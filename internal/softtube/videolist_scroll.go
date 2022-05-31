package softtube

import (
	"github.com/gotk3/gotk3/gtk"
)

type scroll struct {
	scrolledWindow *gtk.ScrolledWindow
}

// toStart : Scrolls to the start of the list
func (s *scroll) toStart() {
	var adjustment = s.scrolledWindow.GetVAdjustment()
	// Possible solution to SIGSEGV, Segmentation fault
	if adjustment == nil {
		return
	}
	adjustment.SetValue(adjustment.GetLower())
	s.scrolledWindow.Show()
}

// toEnd : Scrolls to the end of the list
func (s *scroll) toEnd() {
	var adjustment = s.scrolledWindow.GetVAdjustment()
	// Possible solution to SIGSEGV, Segmentation fault
	if adjustment == nil {
		return
	}
	adjustment.SetValue(adjustment.GetUpper())
	s.scrolledWindow.Show()
}
