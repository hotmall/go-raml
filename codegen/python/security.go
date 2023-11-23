package python

import (
	"strings"

	"path/filepath"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/security"
	"github.com/hotmall/go-raml/raml"
)

const (
	middlewareTypeOauth2 = "oauth2"
)

// python representation of a security scheme
type pythonSecurity struct {
	*security.Security
	ClassName    string
	FileName     string
	ModuleName   string
	template     string
	templateName string
}

func getServerSecurityDefs(schemes map[string]raml.SecurityScheme, templates serverTemplate) []pythonSecurity {
	var defs []pythonSecurity

	// generate oauth2 middleware
	for k, ss := range schemes {
		if ss.Type != security.Oauth2 { // only support oauth2 now
			continue
		}
		sd := security.New(ss, k, "")
		ps := pythonSecurity{Security: &sd}
		title := strings.Title(ps.Name)

		ps.ClassName = "Oauth2" + title
		ps.ModuleName = "oauth2_" + ps.Name
		ps.template = templates.middlewareFile
		ps.templateName = templates.middlewareName
		ps.FileName = ps.ModuleName + ".py"

		defs = append(defs, ps)
	}
	return defs
}

func getClientSecurityDefs(schemes map[string]raml.SecurityScheme, templates clientTemplate) []pythonSecurity {
	var defs []pythonSecurity

	// generate oauth2, Basic Authentication and Pass Through clients
	for k, ss := range schemes {
		if ss.Type != security.Oauth2 && ss.Type != security.BasicAuthentication && ss.Type != security.PassThrough {
			continue
		}
		sd := security.New(ss, k, "")
		ps := pythonSecurity{Security: &sd}
		title := strings.Title(ps.Name)

		switch ps.Type {
		case security.BasicAuthentication:
			ps.ClassName = "BasicAuthClient" + title
			ps.ModuleName = "basicauth_client_" + ps.Name
			ps.template = templates.basicAuthFile
			ps.templateName = templates.basicAuthName
		case security.Oauth2:
			ps.ClassName = "Oauth2Client" + title
			ps.ModuleName = "oauth2_client_" + ps.Name
			ps.template = templates.oauth2File
			ps.templateName = templates.oauth2Name
		case security.PassThrough:
			ps.ClassName = "PassThroughClient" + title
			ps.ModuleName = "passthrough_client_" + ps.Name
			ps.template = templates.passThroughFile
			ps.templateName = templates.passThroughName
		}
		ps.FileName = ps.ModuleName + ".py"

		defs = append(defs, ps)
	}
	return defs
}

// generate security scheme representation in python.
// security scheme is generated as a middleware
func (ps *pythonSecurity) generate(dir string) error {
	filename := filepath.Join(dir, ps.FileName)
	return commons.GenerateFile(ps, ps.template, ps.templateName, filename, true)
}

type middleware struct {
	ImportPath string
	Name       string
	Args       string
	Type       string
}

func newPythonOauth2Middleware(ss raml.DefinitionChoice) (middleware, error) {
	quotedScopes, err := security.GetQuotedScopes(ss)
	if err != nil {
		return middleware{}, err
	}

	importPath, name := pythonOauth2libImportPath(ss.Name)
	return middleware{
		ImportPath: importPath,
		Name:       name,
		Args:       strings.Join(quotedScopes, ", "),
		Type:       middlewareTypeOauth2,
	}, nil
}

// get library import path from a type
func pythonOauth2libImportPath(typ string) (string, string) {
	return libImportPath(security.SecuritySchemeName(typ), "oauth2_")
}

func generateServerSecurity(schemas map[string]raml.SecurityScheme, templates serverTemplate, dir string) error {
	securities := getServerSecurityDefs(schemas, templates)
	for _, s := range securities {
		if err := s.generate(dir); err != nil {
			return err
		}
	}
	return nil
}
