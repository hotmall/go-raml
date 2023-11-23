package python

import (
	"path/filepath"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/raml"
)

// SanicServer represents a sanic asyncio server
type SanicServer struct {
	APIDef       *raml.APIDefinition
	Title        string
	ResourcesDef []pythonResource
	withMain     bool
	apiDocsDir   string
	Template     serverTemplate
	targetDir    string
}

// NewSanicServer creates new sanic server from an RAML file
func NewSanicServer(apiDef *raml.APIDefinition, apiDocsDir, targetDir string, withMain bool,
	libRootURLs []string) *SanicServer {
	var prs []pythonResource
	for _, rd := range getServerResourcesDefs(apiDef) {
		pr := newResource(rd, apiDef, serverKindSanic)
		prs = append(prs, pr)
	}

	// TODO : get rid of this global variables
	globAPIDef = apiDef
	globLibRootURLs = libRootURLs

	templates := templates(serverKindSanic)

	return &SanicServer{
		APIDef:       apiDef,
		Title:        apiDef.Title,
		apiDocsDir:   apiDocsDir,
		withMain:     withMain,
		ResourcesDef: prs,
		Template:     templates,
		targetDir:    targetDir,
	}
}

// APIDocsDir implements generator.Server.APIDocsDir interface
func (s *SanicServer) APIDocsDir() string {
	return s.apiDocsDir
}

// Generate implements generator.Server.Generate interface
func (s *SanicServer) Generate() error {
	if err := generateJSONSchema(s.APIDef, filepath.Join(s.targetDir, handlersDir)); err != nil {
		return err
	}
	if err := s.generateResources(s.targetDir); err != nil {
		return err
	}

	// security scheme
	if err := generateServerSecurity(s.APIDef.SecuritySchemes, s.Template, s.targetDir); err != nil {
		return err
	}

	// python classes and it's helper
	if err := commons.GenerateFile(nil, "./templates/python/client_support.tmpl",
		"client_support", filepath.Join(s.targetDir, "types", "client_support.py"), true); err != nil {
		return err
	}

	_, err := GenerateAllClasses(s.APIDef, filepath.Join(s.targetDir, typesDir), false)
	if err != nil {
		return err
	}

	if err := s.generateMain(s.targetDir); err != nil {
		return err
	}

	return generateEmptyInitPy(s.targetDir)
}

func (s *SanicServer) generateMain(dir string) error {
	if !s.withMain {
		return nil
	}

	// html front page
	if err := commons.GenerateFile(s, "./templates/index.html.tmpl", "index.html",
		filepath.Join(dir, "index.html"), true); err != nil {
		return err
	}

	return commons.GenerateFile(s, "./templates/python/server_main_sanic.tmpl", "server_main_sanic",
		filepath.Join(dir, "app.py"), true)
}

func (s *SanicServer) ShowAPIDocsAndIndex() bool {
	return !resource.HasCatchAllInRootRoute(s.APIDef)
}
