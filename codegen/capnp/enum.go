package capnp

import (
	"fmt"
	"path/filepath"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
	"github.com/pinzolo/casee"
)

type enum struct {
	ID     string
	Name   string
	Fields []field
	lang   string
	pkg    string
}

func newEnum(structName string, prop raml.Property, lang, pkg string) *enum {
	e := enum{
		ID:   getID(),
		Name: "Enum" + structName + casee.ToPascalCase(prop.Name),
		lang: lang,
		pkg:  pkg,
	}
	for k, v := range prop.Enum.([]interface{}) {
		f := field{
			Name: casee.ToCamelCase(v.(string)),
			Num:  k,
		}
		e.Fields = append(e.Fields, f)
	}
	return &e
}

func newEnumFromType(structName string, t raml.Type, lang, pkg string) *enum {
	e := enum{
		ID:   getID(),
		Name: structName,
		lang: lang,
		pkg:  pkg,
	}
	for k, v := range t.Enum.([]interface{}) {
		f := field{
			Name: casee.ToCamelCase(v.(string)),
			Num:  k,
		}
		e.Fields = append(e.Fields, f)
	}
	return &e
}

func (e *enum) generate(dir string) error {
	filename := filepath.Join(dir, e.Name+".capnp")
	return commons.GenerateFile(e, "./templates/capnp/enum_capnp.tmpl", "enum_capnp", filename, true)
}

func (e *enum) Imports() string {
	switch e.lang {
	case "go":
		return `using Go = import "/go.capnp";`
	default:
		return ""
	}

}

func (e *enum) Annotations() []string {
	switch e.lang {
	case "go":
		return e.goAnnotations()
	default:
		return []string{}
	}
}

func (e *enum) goAnnotations() []string {
	return []string{
		fmt.Sprintf(`$Go.package("%v")`, e.pkg),
		fmt.Sprintf(`$Go.import("%v")`, e.pkg),
	}
}
func isEnum(prop raml.Property) bool {
	return prop.Enum != nil
}
