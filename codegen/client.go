package codegen

import (
	"github.com/Jumpscale/go-raml/codegen/gorestful"
	"github.com/Jumpscale/go-raml/codegen/nim"
	"github.com/Jumpscale/go-raml/codegen/python"
	"github.com/Jumpscale/go-raml/raml"
)

// ClientConfig represents client's config
type ClientConfig struct {
	Dir                      string
	PackageName              string
	Lang                     string
	RootImportPath           string
	Kind                     string
	LibRootURLs              []string
	PythonUnmarshallResponse bool
}

// GenerateClient generates client library
//func GenerateClient(apiDef *raml.APIDefinition, dir, packageName, lang, rootImportPath, kind string,
//	libRootURLs []string, pythonUnnmarshallResponse bool) error {
func GenerateClient(apiDef *raml.APIDefinition, conf ClientConfig) error {
	switch conf.Lang {
	case langGo:
		gc, err := gorestful.NewClient(apiDef, conf.Kind, conf.PackageName, conf.RootImportPath, conf.Dir,
			conf.LibRootURLs)
		if err != nil {
			return err
		}
		return gc.Generate()
	case langPython:
		pc := python.NewClient(apiDef, conf.Kind, conf.PythonUnmarshallResponse)
		return pc.Generate(conf.Dir)
	case langNim:
		nc := nim.NewClient(apiDef, conf.Dir)
		return nc.Generate()
	default:
		return errInvalidLang
	}
}
