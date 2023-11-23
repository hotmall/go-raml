package python

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
)

type enumField struct {
	Name  string
	Value string
}
type enum struct {
	Name   string
	Type   string
	Fields []enumField
}

func newEnum(name string, prop raml.Property, fromClass bool) *enum {
	enumName := strings.Title(name) + strings.Title(prop.Name)
	if !fromClass {
		enumName = "Enum" + enumName
	}
	e := enum{
		Name: enumName,
	}

	for _, v := range prop.Enum.([]interface{}) {
		e.Fields = append(e.Fields, newEnumField(v, e))
	}
	return &e
}

func newEnumFromClass(pc *class) *enum {
	prop := raml.Property{
		Type: fmt.Sprint(pc.T.Type),
		Name: "",
		Enum: pc.T.Enum,
	}
	return newEnum(pc.Name(), prop, true)
}

func newEnumField(f interface{}, e enum) enumField {
	var name, val string

	switch v := f.(type) {
	case string:
		name = fmt.Sprintf("%v", v)
		val = fmt.Sprintf(`"%v"`, v)
	case int, float64, float32:
		name = fmt.Sprintf("e%v", v)
		val = fmt.Sprintf("%v", v)
	}
	// ensure name is a valid python variable name
	// 1. the first character must match [a-zA-Z_]
	// 2. all remaining characters must match [a-zA-Z0-9_]

	// first, any characters that don't match any valid variable character are replaced with '_'
	alwaysInvalid := regexp.MustCompile("[^a-zA-Z0-9_]")
	validName := alwaysInvalid.ReplaceAllLiteralString(name, "_")

	// next, if the first character is a number, prepend a '_'
	matched, err := regexp.MatchString("^[0-9]", validName)
	if matched && err == nil {
		validName = "_" + validName
	}
	return enumField{
		Name:  validName,
		Value: val,
	}
}

func (e *enum) generate(dir string) error {
	filename := filepath.Join(dir, e.Name+".py")
	return commons.GenerateFile(e, "./templates/python/enum_python.tmpl", "enum_python", filename, true)
}
