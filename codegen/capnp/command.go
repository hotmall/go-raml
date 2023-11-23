package capnp

import (
	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/types"
	"github.com/hotmall/go-raml/raml"
)

func GenerateCapnp(apiDef *raml.APIDefinition, dir, lang, pkg string) error {
	structs := []Struct{}

	if err := commons.CheckDuplicatedTitleTypes(apiDef); err != nil {
		return err
	}

	// generate types
	for name, t := range types.AllTypes(apiDef, "") {
		var ramlType raml.Type

		switch apiType := t.Type.(type) {
		case string:
			// a reference to a defined/built-in type, no need to create a struct
			// @todo could also be a case of inheritance, do we need to handle this?
			continue
		case types.TypeInBody:
			name = types.PascalCaseTypeName(apiType)
			ramlType = raml.Type{Type: "object"}
			properties := map[string]interface{}{}
			for k, v := range apiType.Properties {
				prop := raml.ToProperty(k, v)
				properties[k] = prop
			}
			ramlType.Properties = properties
		case raml.Type:
			if apiType.IsAlias() {
				// no need to generate capnp of type alias.
				// it will use capnp of the aliased type
				continue
			}
			ramlType = apiType
		}

		s, err := NewStruct(ramlType, name, lang, pkg)
		if err != nil {
			return err
		}
		structs = append(structs, s)
	}

	for _, s := range structs {
		registerType(s.Name)
	}

	for _, s := range structs {
		if err := s.Generate(dir); err != nil {
			return err
		}
	}
	return nil
}
