package nim

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
)

type enumField struct {
	Name string
}

type enum struct {
	name   string
	Fields []enumField
}

func (e enum) NimType() string {
	return "enum"
}

func (e enum) Name() string {
	return e.name
}
func (e enum) FieldsMap() map[string]field {
	return map[string]field{}
}

func newEnum(objName string, prop raml.Property, fromObj bool) *enum {
	e := enum{
		name: strings.Title(objName) + strings.Title(prop.Name),
	}
	if !fromObj {
		e.name = "Enum" + e.Name()
	}
	for _, v := range prop.Enum.([]interface{}) {
		e.Fields = append(e.Fields, newEnumField(v, e))
	}
	return &e
}

func newEnumFromObject(o *object) *enum {
	prop := raml.Property{
		Type: fmt.Sprint(o.T.Type),
		Name: "",
		Enum: o.T.Enum,
	}
	return newEnum(o.Name(), prop, true)
}
func newEnumField(f interface{}, e enum) enumField {
	var name string

	switch v := f.(type) {
	case string:
		name = fmt.Sprintf("%v", v)
	case int:
		name = fmt.Sprintf("e%v=%v", v, v)
	}
	alwaysInvalid := regexp.MustCompile("[^a-zA-Z0-9_=]")
	validName := alwaysInvalid.ReplaceAllLiteralString(name, "_")
	return enumField{
		Name: validName,
	}
}

func (e *enum) generate(dir string) error {
	filename := filepath.Join(dir, e.Name()+".nim")
	return commons.GenerateFile(e, "./templates/nim/enum_nim.tmpl", "enum_nim", filename, true)
}

// FieldsStr is a string representation of all the fields
func (e *enum) FieldsStr() string {
	var names []string
	for _, f := range e.Fields {
		names = append(names, f.Name)
	}
	return strings.Join(names, ", ")
}
