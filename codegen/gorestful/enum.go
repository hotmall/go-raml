package gorestful

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/raml"
)

type enumField struct {
	Name  string
	Type  string
	Value string
}
type enum struct {
	Name   string
	Type   string
	Fields []enumField
	Pkg    string
}

// creates new enum representation.
//
// structName : name of the struct on which this enum reside
// prop : ramlProperty that describe this enum
// pkg : package name
// fromStruct: true if this enum come from struct, false if come from struct's field
func newEnum(structName string, prop raml.Property, pkg string, fromStruct bool) *enum {
	e := enum{
		Name: strings.Title(structName) + strings.Title(prop.Name),
		Type: convertToGoType(prop.TypeString(), prop.Items.Type),
		Pkg:  pkg,
	}
	if !fromStruct {
		e.Name = "Enum" + e.Name
	}
	for _, v := range prop.Enum.([]interface{}) {
		e.Fields = append(e.Fields, newEnumField(v, e))
	}
	return &e
}

func newEnumFromStruct(sd *structDef) *enum {
	prop := raml.Property{
		Type: fmt.Sprint(sd.T.Type),
		Name: "",
		Enum: sd.T.Enum,
	}
	return newEnum(sd.Name, prop, sd.PackageName, true)
}

// creates new enum field
// f : enum's member
// e : the enum representation
func newEnumField(f interface{}, e enum) enumField {
	var val string

	// field name = enum name + field name
	name := fmt.Sprintf("%v%v", e.Name, f)

	// first, any characters that don't match any valid variable character are replaced with '_'
	alwaysInvalid := regexp.MustCompile("[^a-zA-Z0-9_]")
	validName := alwaysInvalid.ReplaceAllLiteralString(name, "_")

	switch v := f.(type) {
	case string:
		val = fmt.Sprintf(`"%v"`, v)
	case int:
		val = fmt.Sprintf("%v", v)
	}
	return enumField{
		Name:  validName,
		Value: val,
		Type:  e.Name,
	}
}

func (e *enum) generate(dir string) error {
	filename := filepath.Join(dir, e.Name+".go")
	return commons.GenerateFile(e, "./templates/golang/enum_go.tmpl", "enum_go", filename, true)
}
