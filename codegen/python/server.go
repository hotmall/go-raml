package python

import (
	log "github.com/Sirupsen/logrus"

	"github.com/Jumpscale/go-raml/codegen/generator"
	"github.com/Jumpscale/go-raml/raml"
)

var (
	globLibRootURLs []string
)

const (
	serverKindSanic       = "sanic"
	serverKindFlask       = "flask"
	serverKindGeventFlask = "gevent-flask"
	typesDir              = "pytypes"
	handlersDir           = "handlers"
)

// NewServer creates a new python server
func NewServer(kind string, apiDef *raml.APIDefinition, apiDocsDir, targetDir string,
	withMain bool, libRootURLs []string, name string) generator.Server {
	switch kind {
	case "", serverKindFlask:
		return NewFlaskServer(apiDef, apiDocsDir, targetDir, withMain, libRootURLs, false, name)
	case serverKindGeventFlask:
		return NewFlaskServer(apiDef, apiDocsDir, targetDir, true, libRootURLs, true, name)
	case serverKindSanic:
		return NewSanicServer(apiDef, apiDocsDir, targetDir, withMain, libRootURLs)

	default:
		log.Fatalf("Invalid kind of python server : %v", kind)
		return nil
	}
}
