package commands

import (
	"github.com/hotmall/go-raml/codegen"
	"github.com/hotmall/go-raml/raml"

	log "github.com/Sirupsen/logrus"
)

//ClientCommand is executed to generate client from a RAML specification
type DocsCommand struct {
	Format     string
	OutputFile string //target dir
	RamlFile   string //raml file
}

//Execute generates a client from a RAML specification
func (command *DocsCommand) Execute() error {
	log.Debug("Generating api docs in %s", command.Format)
	apiDef := new(raml.APIDefinition)
	err := raml.ParseFile(command.RamlFile, apiDef)
	if err != nil {
		return err
	}
	return codegen.GenerateDocs(apiDef, command.Format, command.OutputFile)
}
