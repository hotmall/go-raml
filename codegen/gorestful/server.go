package gorestful

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/resource"
	"github.com/Jumpscale/go-raml/raml"
)

// global variables
// it is needed for libraries support
var (
	// root import path
	globRootImportPath string

	// global value of API definition
	globAPIDef *raml.APIDefinition

	globLibRootURLs []string

	globKind string
)

// Server represents a Go server
type Server struct {
	apiDef         *raml.APIDefinition
	ResourcesDef   []*goResource
	PackageName    string // Name of the package this server resides in
	apiDocsDir     string // apidocs directory. apidocs won't be generated if it is empty
	withMain       bool   // true if we need to generate main file
	RootImportPath string
	TargetDir      string   // root directory of the generated code
	libsRootURLs   []string // root URLs of the libraries
}

// NewServer creates a new Golang server
func NewServer(kind string, apiDef *raml.APIDefinition, packageName, apiDocsDir, rootImportPath string,
	withMain bool, targetDir string, libsRootURLs []string) *Server {
	// global variables
	rootImportPath = setRootImportPath(rootImportPath, targetDir)
	globAPIDef = apiDef
	globRootImportPath = rootImportPath
	globLibRootURLs = libsRootURLs
	globKind = kind

	return &Server{
		apiDef:         apiDef,
		PackageName:    packageName,
		apiDocsDir:     apiDocsDir,
		withMain:       withMain,
		RootImportPath: rootImportPath,
		TargetDir:      targetDir,
		libsRootURLs:   libsRootURLs,
	}
}

// APIDocsDir implements codegen.Server.APIDocsDir interface
func (gs *Server) APIDocsDir() string {
	return gs.apiDocsDir
}

// Generate implements codegen.Server.Generate interface
func (gs *Server) Generate() error {
	if err := commons.CheckDuplicatedTitleTypes(gs.apiDef); err != nil {
		return err
	}
	if gs.RootImportPath == "" {
		return fmt.Errorf("invalid import path = empty. please set --import-path or set target dir under gopath")
	}
	// helper package
	gh := goramlHelper{
		rootImportPath: gs.RootImportPath,
		packageName:    "goraml",
		packageDir:     "goraml",
		command:        "server",
		kind:           globKind,
	}
	if err := gh.generate(gs.TargetDir); err != nil {
		return err
	}

	if err := generateAllStructs(gs.apiDef, gs.TargetDir); err != nil {
		return err
	}

	if err := generateErrorStruct(gs.apiDef, gs.TargetDir); err != nil {
		return err
	}

	if err := generateAnyStruct(gs.apiDef, gs.TargetDir); err != nil {
		return err
	}

	// security scheme
	if err := generateSecurity(gs.apiDef.SecuritySchemes, gs.TargetDir, gs.PackageName); err != nil {
		log.Errorf("failed to generate security scheme:%v", err)
		return err
	}

	// genereate resources
	rds, err := gs.generateServerResources(gs.TargetDir)
	if err != nil {
		return err
	}
	gs.ResourcesDef = rds

	// libraries
	if err := generateLibraries(gs.apiDef.Libraries, gs.TargetDir, gs.libsRootURLs); err != nil {
		return err
	}

	// generate main
	if gs.withMain {
		// HTML front page
		if err := commons.GenerateFile(gs, "./templates/index.html.tmpl", "index.html", filepath.Join(gs.TargetDir, "index.html"), false); err != nil {
			return err
		}
		// .gitignore file
		if err := commons.GenerateFile(gs, "./templates/golang/gorestful_gitignore.tmpl", "gorestful_gitignore", filepath.Join(gs.TargetDir, ".gitignore"), false); err != nil {
			return err
		}
		// Makefile
		if err := commons.GenerateFile(gs, "./templates/golang/gorestful_makefile.tmpl", "gorestful_makefile", filepath.Join(gs.TargetDir, "Makefile"), false); err != nil {
			return err
		}
		// main file
		return commons.GenerateFile(gs, "./templates/golang/gorestful_main.tmpl", "gorestful_main", filepath.Join(gs.TargetDir, "main.go"), false)
	}

	return nil
}

// Title returns title of this server
func (gs Server) Title() string {
	return gs.apiDef.Title
}

// RouteImports returns all resources routes
func (gs Server) RouteImports() []string {
	imports := make(map[string]struct{})

	baseAPIDir := gs.RootImportPath + "/" + serverAPIDir
	for _, rd := range gs.ResourcesDef {
		imports[baseAPIDir+"/"+rd.PackageName] = struct{}{}
	}
	return commons.MapToSortedStrings(imports)
}

// ShowAPIDocsAndIndex returns if show api docs and index
func (gs Server) ShowAPIDocsAndIndex() bool {
	return !resource.HasCatchAllInRootRoute(gs.apiDef)
}
