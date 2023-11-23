package codegen

import (
	"github.com/hotmall/go-raml/codegen/python"
	"github.com/hotmall/go-raml/raml"
)

func GeneratePythonCapnp(apiDef *raml.APIDefinition, dir string) error {
	return python.GeneratePythonCapnpClasses(apiDef, dir)
}
