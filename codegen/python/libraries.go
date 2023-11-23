package python

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/libraries"
	"github.com/hotmall/go-raml/raml"
)

type library struct {
	*raml.Library
	baseDir string
	dir     string
}

func newLibrary(lib *raml.Library, baseDir string) *library {
	pl := library{
		Library: lib,
		baseDir: baseDir,
	}
	pl.Filename = libraries.StripLibRootURL(pl.Filename, globLibRootURLs)

	// package directory : filename without the extension
	relDir := libRelDir(pl.Filename)
	pl.dir = commons.NormalizePkgName(filepath.Join(pl.baseDir, relDir))

	return &pl
}

// generate code of all libraries
func generateLibraries(libraries map[string]*raml.Library, baseDir string) error {
	for _, ramlLib := range libraries {
		pl := newLibrary(ramlLib, baseDir)

		if err := pl.generate(); err != nil {
			return err
		}
	}
	return nil
}

// generate code of this library
func (l *library) generate() error {
	// write empty __init__.py in each dir
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		return generateEmptyInitPy(path)
	}

	if err := filepath.Walk(l.baseDir, walkFn); err != nil {
		return err
	}

	// json schema
	if err := generateJSONSchema(&raml.APIDefinition{
		Types: l.Types,
	}, l.dir); err != nil {
		log.Errorf("failed to generate jsonschema:%v", err)
		return err
	}

	// security schemes
	if err := generateServerSecurity(l.SecuritySchemes, templates(serverKindFlask), l.dir); err != nil {
		return err
	}

	// included libraries
	for _, ramlLib := range l.Libraries {
		childLib := newLibrary(ramlLib, l.baseDir)
		if err := childLib.generate(); err != nil {
			return err
		}
	}
	return nil
}

func libImportPath(typ, prefix string) (string, string) {
	// library use '.'
	if strings.Index(typ, ".") < 0 {
		return prefix + typ, prefix + typ
	}

	splitted := strings.Split(typ, ".")
	if len(splitted) != 2 {
		log.Fatalf("pythonLibImportPath invalid typ:" + typ)
	}
	// library name in the current document
	libName := splitted[0]

	// raml file of this lib
	libDir, libFile := globAPIDef.FindLibFile(commons.DenormalizePkgName(libName))

	if libFile == "" {
		log.Fatalf("pythonLibImportPath() can't find library : %v", libName)
	}

	libPath := libraries.JoinPath(libDir, libFile, globLibRootURLs)

	// relative lib package
	libPkg := libRelDir(libPath)

	return strings.Replace(commons.NormalizePkgName(libPkg), "/", ".", -1) + "." + prefix + splitted[1], prefix + splitted[1]
}

// get relative lib directory from library filename
// this relative directory will be used as:
// - lib package name
// - lib files directory
func libRelDir(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
