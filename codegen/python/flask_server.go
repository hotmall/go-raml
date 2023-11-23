package python

import (
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/raml"
)

// FlaskServer represents a flask server
type FlaskServer struct {
	APIDef       *raml.APIDefinition
	Title        string
	ResourcesDef []pythonResource
	WithMain     bool
	apiDocsDir   string
	Gevent       bool
	Template     serverTemplate
	targetDir    string
	Name         string
}

// NewFlaskServer creates new flask server from an RAML file
func NewFlaskServer(apiDef *raml.APIDefinition, apiDocsDir, targetDir string,
	withMain bool, libRootURLs []string, gevent bool, name string) *FlaskServer {

	// TODO : get rid of this global variables
	globAPIDef = apiDef
	globLibRootURLs = libRootURLs

	// generates resource
	var prs []pythonResource
	for _, rd := range getServerResourcesDefs(apiDef) {
		pr := newResource(rd, apiDef, serverKindFlask)
		prs = append(prs, pr)
	}

	templates := templates(serverKindFlask)

	return &FlaskServer{
		APIDef:       apiDef,
		Title:        apiDef.Title,
		apiDocsDir:   apiDocsDir,
		WithMain:     withMain,
		ResourcesDef: prs,
		Gevent:       gevent,
		Template:     templates,
		targetDir:    targetDir,
		Name:         name,
	}
}

// APIDocsDir implements generator.Server.APIDocsDir interface
func (ps *FlaskServer) APIDocsDir() string {
	return ps.apiDocsDir
}

// Generate implements generator.Server.Generate interface
func (ps *FlaskServer) Generate() error {
	// python classes and it's helper
	typesPath := filepath.Join(ps.targetDir, typesDir)

	if err := commons.GenerateFile(nil, "./templates/python/client_support.tmpl",
		"client_support", filepath.Join(typesPath, "client_support.py"), true); err != nil {
		return err
	}

	if ps.Gevent {
		if err := commons.GenerateFile(ps, "./templates/python/server_gevent.tmpl", "server_gevent", filepath.Join(ps.targetDir, "server.py"), true); err != nil {
			return err
		}

		if err := GeneratePythonCapnpClasses(ps.APIDef, typesPath); err != nil {
			return nil
		}
	} else {
		_, err := GenerateAllClasses(ps.APIDef, typesPath, ps.Gevent)
		if err != nil {
			return err
		}
	}

	// json schema
	if err := generateJSONSchema(ps.APIDef, filepath.Join(ps.targetDir, handlersDir)); err != nil {
		return err
	}

	// security scheme
	if err := generateServerSecurity(ps.APIDef.SecuritySchemes, ps.Template, ps.targetDir); err != nil {
		log.Errorf("failed to generate security scheme:%v", err)
		return err
	}

	// generate resources
	if err := generateResources(ps.ResourcesDef, ps.Template, ps.targetDir); err != nil {
		return err
	}

	// libraries
	if err := generateLibraries(ps.APIDef.Libraries, ps.targetDir); err != nil {
		return err
	}

	// requirements.txt file
	if err := commons.GenerateFile(nil, "./templates/python/requirements_python.tmpl", "requirements_python",
		filepath.Join(ps.targetDir, "requirements.txt"), false); err != nil {
		return err
	}

	// etc/conf/gunicorn.py
	if err := commons.GenerateFile(nil, "./templates/python/server_gunicorn.tmpl", "server_gunicorn",
		filepath.Join(ps.targetDir, "/etc/conf/gunicorn.py"), false); err != nil {
		return err
	}
	// utils/logging.py
	if err := commons.GenerateFile(nil, "./templates/python/server_logging.tmpl", "server_logging",
		filepath.Join(ps.targetDir, "/utils/logging.py"), true); err != nil {
		return err
	}

	// utils/__init__.py
	if err := commons.GenerateFile(nil, "./templates/python/server_logging_init.tmpl", "server_logging_init",
		filepath.Join(ps.targetDir, "/utils/__init__.py"), false); err != nil {
		return err
	}

	// generate main
	if ps.WithMain {
		// generate HTML front page
		if err := commons.GenerateFile(ps, "./templates/index.html.tmpl", "index.html", filepath.Join(ps.targetDir, "index.html"), true); err != nil {
			return err
		}
		// generate Makefile
		if err := commons.GenerateFile(ps, "./templates/python/server_makefile.tmpl", "python_server_makefile", filepath.Join(ps.targetDir, "Makefile"), false); err != nil {
			return err
		}
		// main file
		if err := commons.GenerateFile(ps, "./templates/python/server_main_flask.tmpl", "server_main_flask", filepath.Join(ps.targetDir, "app.py"), true); err != nil {
			return err
		}
	}

	return nil
}

func (ps *FlaskServer) ShowAPIDocsAndIndex() bool {
	return !resource.HasCatchAllInRootRoute(ps.APIDef)
}
