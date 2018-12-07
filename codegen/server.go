package codegen

import (
	"errors"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/Jumpscale/go-raml/codegen/apidocs"
	"github.com/Jumpscale/go-raml/codegen/generator"
	"github.com/Jumpscale/go-raml/codegen/gorestful"
	"github.com/Jumpscale/go-raml/codegen/nim"
	"github.com/Jumpscale/go-raml/codegen/python"
	"github.com/Jumpscale/go-raml/codegen/tarantool"
	"github.com/Jumpscale/go-raml/raml"
)

var (
	errInvalidLang = errors.New("invalid language")
)

// Server represents server code
type Server struct {
	RAMLFile         string
	Dir              string // Destination directory
	PackageName      string
	Lang             string
	APIDocsDir       string // API Docs directory
	RootImportPath   string
	WithMain         bool   // true if we also generate main file
	Kind             string // currently only used by python client : sanic/flask
	APIFilePerMethod bool   // true if we want to generate one API file per API method
	LibRootURLs      []string
}

// Generate generates API server files
func (s *Server) Generate() error {
	apiDef := new(raml.APIDefinition)
	// parse the raml file
	dir, fileName := filepath.Split(s.RAMLFile)
	ramlBytes, err := raml.ParseReadFile(dir, fileName, apiDef)
	if err != nil {
		return err
	}

	var generator generator.Server

	switch s.Lang {
	case langGo:
		generator = gorestful.NewServer(s.Kind, apiDef, s.PackageName, s.APIDocsDir, s.RootImportPath, s.WithMain,
			s.Dir, s.LibRootURLs)
	case langPython:
		generator = python.NewServer(s.Kind, apiDef, s.APIDocsDir, s.Dir, s.WithMain, s.LibRootURLs)
	case langNim:
		generator = nim.NewServer(apiDef, s.APIDocsDir, s.Dir)
	case langTarantool:
		generator = tarantool.NewServer(apiDef, s.APIDocsDir, s.Dir)
	default:
		return errInvalidLang
	}
	err = generator.Generate()
	if err != nil {
		return err
	}

	if s.APIDocsDir == "" {
		return nil
	}

	if s.Lang == langNim {
		s.APIDocsDir = "public/" + s.APIDocsDir
	}

	log.Infof("Generating API Docs to %v", s.APIDocsDir)

	return apidocs.Generate(apiDef, s.RAMLFile, ramlBytes, filepath.Join(s.Dir, s.APIDocsDir), s.LibRootURLs)
}
