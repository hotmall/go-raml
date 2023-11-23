package python

import (
	"path/filepath"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/raml"
)

const (
	clientNameRequests       = "requests"
	clientNameGeventRequests = "gevent-requests"
	clientNameAiohttp        = "aiohttp"
)

var (
	globAPIDef *raml.APIDefinition
)

// Client represents a python client
type Client struct {
	Name               string
	APIDef             *raml.APIDefinition
	BaseURI            string
	Classes            []string
	Services           map[string]*service
	Kind               string
	Template           clientTemplate
	UnmarshallResponse bool // true if response body should be unmarshalled into python class
	Securities         []pythonSecurity
}

// NewClient creates a python Client
func NewClient(apiDef *raml.APIDefinition, kind string, unmarshallResponse bool) Client {
	services := map[string]*service{}
	for uri, res := range apiDef.Resources {
		rd := resource.New(apiDef, &res, commons.NormalizeURITitle(apiDef.Title), false)
		services[uri] = newService(uri, &rd, unmarshallResponse)
	}
	switch kind {
	case "":
		kind = clientNameRequests
	case clientNameRequests, clientNameAiohttp, clientNameGeventRequests:
	default:
		log.Fatalf("invalid client kind:%v", kind)
	}

	c := Client{
		Name:               commons.NormalizeURI(apiDef.Title),
		APIDef:             apiDef,
		BaseURI:            apiDef.BaseURI,
		Services:           services,
		Kind:               kind,
		UnmarshallResponse: unmarshallResponse,
	}
	if strings.Index(c.BaseURI, "{version}") > 0 {
		c.BaseURI = strings.Replace(c.BaseURI, "{version}", apiDef.Version, -1)
	}
	c.initTemplates()
	c.Securities = getClientSecurityDefs(apiDef.SecuritySchemes, c.Template)

	return c
}

// generate empty __init__.py without overwrite it
func generateEmptyInitPy(dir string) error {
	return commons.GenerateFile(nil, "./templates/python/init_py.tmpl", "init_py", filepath.Join(dir, "__init__.py"), false)
}

func (c Client) ClientKind() string {
	return c.Kind
}

// Generate generates python client library files
func (c Client) Generate(dir string) error {
	globAPIDef = c.APIDef

	// generate helper
	if err := commons.GenerateFile(nil, "./templates/python/client_utils_python.tmpl", "client_utils_python", filepath.Join(dir, "client_utils.py"), true); err != nil {
		return err
	}

	// generate http client
	if err := commons.GenerateFile(nil, c.Template.httpClientFile, c.Template.httpClientName, filepath.Join(dir, "http_client.py"), true); err != nil {
		return err
	}

	if err := c.generateServices(dir); err != nil {
		return err
	}

	if err := c.generateSecurity(dir); err != nil {
		return err
	}

	// python classes
	classes, err := GenerateAllClasses(c.APIDef, dir, false)
	if err != nil {
		return err
	}

	// helper for python classes
	if err := commons.GenerateFile(nil, "./templates/python/client_support.tmpl",
		"client_support", filepath.Join(dir, "client_support.py"), true); err != nil {
		return err
	}

	if c.UnmarshallResponse {
		if err := commons.GenerateFile(nil, "./templates/python/unmarshall_error.tmpl",
			"unmarshall_error", filepath.Join(dir, "unmarshall_error.py"), true); err != nil {
			return err
		}
		if err := commons.GenerateFile(nil, "./templates/python/unhandled_api_error.tmpl",
			"unhandled_api_error", filepath.Join(dir, "unhandled_api_error.py"), true); err != nil {
			return err
		}
	}

	sort.Strings(classes)
	c.Classes = classes
	return commons.GenerateFile(c, c.Template.initFile, c.Template.initName, filepath.Join(dir, "__init__.py"), true)
}

func (c Client) generateServices(dir string) error {
	for _, s := range c.Services {
		if err := s.generate(c.Template, dir); err != nil {
			return err
		}
	}
	return nil
}

func (c Client) generateSecurity(dir string) error {
	for _, s := range c.Securities {
		if err := s.generate(dir); err != nil {
			return err
		}
	}
	return nil
}
