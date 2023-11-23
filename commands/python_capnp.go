package commands

import (
	log "github.com/Sirupsen/logrus"

	"github.com/hotmall/go-raml/codegen"
	"github.com/hotmall/go-raml/raml"
)

// PythonCapnp is executed to generate python class with capnp loader from RAML specification
type PythonCapnp struct {
	Dir      string //target dir
	RAMLFile string //raml file
}

//Execute generates a client from a RAML specification
func (command *PythonCapnp) Execute() error {
	var apiDef raml.APIDefinition

	log.Debug("Generating python classes with capnp conversion")

	err := raml.ParseFile(command.RAMLFile, &apiDef)
	if err != nil {
		return err
	}
	return codegen.GeneratePythonCapnp(&apiDef, command.Dir)
}
