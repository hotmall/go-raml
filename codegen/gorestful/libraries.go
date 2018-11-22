package gorestful

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/libraries"
	"github.com/Jumpscale/go-raml/raml"
)

// library defines an RAML library
// it is implemented as package in Go.
// Library defined in RAML using `uses` keyword.
// - key become library's package name
// - if value contain directory, the directories become root directory
//   of the generated package
// example :
// files: libraries/files.raml -> generated as `files` package in `libraries` directory
// types-lib: lib-types.raml  -> generated as `types_lib` package in current directory
type goLibrary struct {
	name string
	*raml.Library
	PackageName   string // package name of the library
	rootTargetDir string // root target directory of the generated code
	parentDir     string // directory of raml docs that including this library
	libRootURLs   []string
}

// create new library instance

func newGoLibrary(name string, lib *raml.Library, parentDir, rootTargetDir string,
	libRootURLs []string) *goLibrary {
	gl := &goLibrary{
		Library:       lib,
		PackageName:   commons.NormalizePkgName(name),
		parentDir:     parentDir,
		rootTargetDir: rootTargetDir,
		libRootURLs:   libRootURLs,
		name:          name,
	}
	return gl
}

// target dir of this go library
func (gl *goLibrary) targetDir() string {
	baseDir := gl.parentDir
	if libraries.IsRemote(gl.Filename) {
		baseDir = gl.rootTargetDir
	}

	fileDir := goLibPackageDir(gl.name, libraries.StripLibRootURL(gl.Filename, gl.libRootURLs))

	return filepath.Join(baseDir, fileDir)
}

// generate code of all libraries
func generateLibraries(libraries map[string]*raml.Library, baseDir string, libsRootURLs []string) error {
	for name, ramlLib := range libraries {
		//l := newGoLibrary(name, ramlLib, baseDir, baseDir, libsRootURLs)
		l := newGoLibrary(name, ramlLib, baseDir, baseDir, libsRootURLs)
		if err := l.generate(); err != nil {
			return err
		}
	}
	return nil
}

// generate code of this library
func (gl *goLibrary) generate() error {
	// generate all Type structs
	if err := generateStructs(gl.Types, gl.targetDir()); err != nil {
		return err
	}

	// security schemes
	if err := generateSecurity(gl.SecuritySchemes, gl.targetDir(), gl.PackageName); err != nil {
		return err
	}

	// included libraries
	for name, ramlLib := range gl.Libraries {
		childLib := newGoLibrary(name, ramlLib, gl.targetDir(), gl.rootTargetDir, globLibRootURLs)
		if err := childLib.generate(); err != nil {
			return err
		}
	}
	return nil
}

// get library import path from a type
func libImportPath(rootImportPath, typ string, libRootURLs []string) string {
	// all library use '.',
	// return nothing if it is not a library
	if strings.Index(typ, ".") < 0 {
		return ""
	}

	if strings.Index(typ, "json.RawMessage") >= 0 {
		return `"encoding/json"`
	}

	// library name in the current document
	libName := strings.Split(typ, ".")[0]

	if libName == "goraml" { // special package name, reserved for goraml
		return fmt.Sprintf(`"%s"`, joinImportPath(rootImportPath, "goraml"))
	}

	// raml file of this lib
	libDir, libFile := globAPIDef.FindLibFile(commons.DenormalizePkgName(libName))

	if libFile == "" {
		return ""
	}

	libPath := libraries.JoinPath(libDir, libFile, libRootURLs)

	importPath := joinImportPath(rootImportPath, path.Join(goLibPackageDir(libName, libPath), typeDir))
	return aliasLibTypeImportPath(importPath)
}

// generate import line with alias style
// example:
//  a.com/libraries/libname/types ->  libname_types "a.com/libraries/libname/types"
// because we generate all raml Types under `types` directory and package,
// all types from libraries is going to have import line with this format:
//    root_import_path/libraries/libname/types
// we create the alias to avoid conflict with other `types` directory/package.
func aliasLibTypeImportPath(path string) string {
	elems := strings.Split(path, "/")
	n := len(elems)

	return fmt.Sprintf(`%v_%v "%v"`, elems[n-2], elems[n-1], path)
}

// returns Go package directory of a library
// name is library name. filename is library file name.
// for the rule, see comment of `type goLibrary struct`
func goLibPackageDir(name, filename string) string {
	dir := filepath.Join(filepath.Dir(filename), name)

	// escape last dir element
	elems := strings.Split(dir, "/")
	elems[len(elems)-1] = commons.NormalizeIdentifier(elems[len(elems)-1])
	return strings.Join(elems, "/")
}

// joinImportPath should always use `/`.
// this is not the case in windows.
// so we need to split it based on file separator and then join again using `/`
func joinImportPath(baseDir, dir string) string {
	baseDirElems := strings.Split(baseDir, string(filepath.Separator))
	dirElems := strings.Split(dir, string(filepath.Separator))
	elems := append(baseDirElems, dirElems...)
	return path.Join(elems...)
}
