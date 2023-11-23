package commands

import (
	"strings"

	"github.com/hotmall/go-raml/codegen"
	"github.com/hotmall/go-raml/raml"

	log "github.com/Sirupsen/logrus"
)

//ClientCommand is executed to generate client from a RAML specification
type ClientCommand struct {
	Language    string
	Dir         string //target dir
	RamlFile    string //raml file
	PackageName string //package name in the generated go source files
	ImportPath  string
	Kind        string

	// Root URL of the libraries.
	// Usefull if we want to use remote libraries.
	// Example:
	//   root url     = http://localhost.com/lib
	//   library file = http://localhost.com/lib/libraries/security.raml
	//	 the library file is going to treated the same as local : libraries/security.raml
	LibRootURLs string

	// If true, python client will unmarshall the response
	// Other languages already unmarshall the response
	PythonUnmarshallResponse bool
}

//Execute generates a client from a RAML specification
func (command *ClientCommand) Execute() error {
	log.Debug("Generating a rest client for ", command.Language)
	apiDef := new(raml.APIDefinition)
	err := raml.ParseFile(command.RamlFile, apiDef)
	if err != nil {
		return err
	}
	conf := codegen.ClientConfig{
		Dir:                      command.Dir,
		PackageName:              command.PackageName,
		Lang:                     command.Language,
		RootImportPath:           command.ImportPath,
		Kind:                     command.Kind,
		LibRootURLs:              strings.Split(command.LibRootURLs, ","),
		PythonUnmarshallResponse: command.PythonUnmarshallResponse,
	}

	return codegen.GenerateClient(apiDef, conf)
}
