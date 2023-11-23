package codegen

import (
	"fmt"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
)

type markdownDocs struct {
	api    *raml.APIDefinition
	output string
}

func (d *markdownDocs) flatten(resources map[string]*raml.Resource, base string, out map[string]*raml.Resource) {
	for name, r := range resources {
		fullName := fmt.Sprintf("%s%s", base, name)
		if len(r.Methods) > 0 {
			out[fullName] = r
		}

		d.flatten(r.Nested, fullName, out)
	}
}

func (d *markdownDocs) generate() error {
	resources := make(map[string]*raml.Resource)
	for k, v := range d.api.Resources {
		func(resource raml.Resource) {
			resources[k] = &resource
		}(v)
	}

	flat := make(map[string]*raml.Resource)
	d.flatten(resources, "", flat)

	ctx := map[string]interface{}{
		"Api":       d.api,
		"Resources": flat,
		"Property":  raml.ToProperty,
	}

	return commons.GenerateFile(ctx, "./templates/docs_markdown.tmpl", "docs_markdown", d.output, true)
}
