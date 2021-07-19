package main

import (
	"errors"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam/framework"
)

type SoftBuilder struct {
	builder *gtk.Builder
}

func newSoftBuilder(fileName string) *SoftBuilder {
	builder := new(SoftBuilder)
	builder.createBuilder(fileName)
	return builder
}

func (s *SoftBuilder) createBuilder(gladeFileName string) {
	fw := framework.NewFramework()
	gladePath := fw.Resource.GetResourcePath(gladeFileName)
	if gladePath == "" {
		err := errors.New("resource path not found")
		logger.LogError(err)
		panic(err)
	}

	builder, err := gtk.BuilderNewFromFile(gladePath)
	if err != nil {
		logger.LogError(err)
		panic(err)
	}

	s.builder = builder
}

func (s *SoftBuilder) getObject(name string) glib.IObject {
	obj, err := s.builder.GetObject(name)
	if err != nil {
		panic(err)
	}

	return obj
}

