package builder

import (
	"github.com/gotk3/gotk3/gtk"
)

type Builder struct {
	builder *gtk.Builder
}

// NewBuilder creates a gtk.Builder and wraps it in a Builder struct
func NewBuilder(glade string) *Builder {
	// Create a new builder
	b, err := gtk.BuilderNewFromString(glade)
	if err != nil {
		panic(err)
	}
	return &Builder{b}
}

//// GetObject gets a gtk object by name
//func (b *Builder) GetObject(name string) glib.IObject {
//	if b.builder == nil {
//		panic("missing builder")
//	}
//	obj, err := b.builder.GetObject(name)
//	if err != nil {
//		panic(err)
//	}
//
//	return obj
//}
