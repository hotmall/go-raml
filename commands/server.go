package commands

import (
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/hotmall/go-raml/codegen"
)

// ServerCommand is executed to generate a go server from a RAML specification
type ServerCommand struct {
	Language         string // target language
	Kind             string
	Dir              string //target dir
	RamlFile         string //raml file
	PackageName      string //package name in the generated go source files
	NoMainGeneration bool   //do not generate a main.go file
	ImportPath       string // root import path of the code, such as : github.com/jumpscale/restapi
	NoAPIDocs        bool   // do not generate API Docs in /apidocs/ endpoint

	// Root URL of the libraries.
	// Usefull if we want to use remote libraries.
	// Example:
	//   root url     = http://localhost.com/lib
	//   library file = http://localhost.com/lib/libraries/security.raml
	//	 the library file is going to treated the same as local : libraries/security.raml
	LibRootURLs string
}

// Execute generates a Go server from an RAML specification
func (command *ServerCommand) Execute() error {
	var apiDocsDir string

	log.Infof("Generating a %v server", command.Language)

	if !command.NoAPIDocs {
		apiDocsDir = "apidocs"
	}

	cs := codegen.Server{
		RAMLFile:       command.RamlFile,
		Kind:           command.Kind,
		Dir:            command.Dir,
		PackageName:    command.PackageName,
		Lang:           command.Language,
		APIDocsDir:     apiDocsDir,
		RootImportPath: command.ImportPath,
		WithMain:       !command.NoMainGeneration,
		LibRootURLs:    strings.Split(command.LibRootURLs, ","),
	}
	return cs.Generate()
}
