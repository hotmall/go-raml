package gorestful

import (
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/Jumpscale/go-raml/codegen/commons"
)

var (
	globGoramlPkgDir string // global variable, hold `goraml` package dir
)

// goramlHelper represents helper package
// that is not described in .raml file but needed by
// generated code.
type goramlHelper struct {
	rootImportPath string // only used by server
	packageName    string
	packageDir     string
	command        string
	kind           string
}

func (gh goramlHelper) generate(dir string) error {
	globGoramlPkgDir = gh.packageDir
	pkgDir := filepath.Join(dir, gh.packageDir)

	/// dates
	d := dateGen{PackageName: gh.packageName}
	if err := d.generate(pkgDir); err != nil {
		log.Errorf("generate() failed to generate date files:%v", err)
		return err
	}

	// generate struct validator
	if err := generateInputValidator(gh.packageName, pkgDir); err != nil {
		return err
	}

	// error helper
	ctx := struct {
		PackageName string
	}{
		PackageName: gh.packageName,
	}
	if gh.command == "client" && gh.kind == "requests" {
		if err := commons.GenerateFile(ctx, "./templates/golang/requests_go_api_error.tmpl", "requests_go_api_error",
			filepath.Join(pkgDir, "api_error.go"), true); err != nil {
			return err
		}
	} else {
		if err := commons.GenerateFile(ctx, "./templates/golang/goraml_go_api_error.tmpl", "goraml_go_api_error",
			filepath.Join(pkgDir, "api_error.go"), true); err != nil {
			return err
		}
	}
	return nil
}
