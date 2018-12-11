package gorestful

import (
	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/raml"
)

func generateErrorStruct(apiDef *raml.APIDefinition, dir string) error {
	filename := dir + "/types/error.go"
	if err := commons.GenerateFile(apiDef, "./templates/golang/gorestful_error.tmpl",
		"gorestful_error_template", filename, true); err != nil {
		return err
	}
	return nil
}
