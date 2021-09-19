package gorestful

import (
	"fmt"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/raml"
	"github.com/Jumpscale/go-raml/utils"
)

// FieldDef defines a field of a struct
type fieldDef struct {
	Name          string // field name
	fieldType     string // field type
	IsComposition bool   // composition type
	IsOmitted     bool   // omitted empty
	UniqueItems   bool
	Enum          *enum // not nil if this field contains enum
	Default       interface{}

	Validators string
}

// newFieldDef creates new struct field from raml property.
func newFieldDef(apiDef *raml.APIDefinition, structName string, prop raml.Property, pkg string) fieldDef {
	var (
		fieldType = prop.TypeString()               // the field type
		basicType = commons.GetBasicType(fieldType) // basic type of the field type
	)

	// for the types, check first if it is user defined type
	if _, ok := apiDef.Types[basicType]; ok {
		titledType := strings.Title(basicType)

		// check if it is a recursive type
		if titledType == strings.Title(structName) {
			titledType = "*" + titledType // add `pointer`, otherwise compiler will complain
		}

		// if it is not array type and is not required
		if !commons.IsArrayType(fieldType) && !prop.Required {
			titledType = "*" + titledType // add `pointer`
		}

		// use strings.Replace instead of simple assignment because the fieldType
		// might be an array
		fieldType = strings.Replace(fieldType, basicType, titledType, 1)
	}
	fieldType = convertToGoType(fieldType, prop.Items.Type)

	fd := fieldDef{
		Name:      formatFieldName(prop.Name),
		fieldType: fieldType,
		IsOmitted: !prop.Required,
	}

	fd.buildValidators(prop)

	if prop.IsEnum() {
		fd.Enum = newEnum(structName, prop, pkg, false)
		fd.fieldType = fd.Enum.Name
	}

	return fd
}

func (fd fieldDef) Type() string {
	// doesn't have "." -> doesnt import from other package
	if !strings.Contains(fd.fieldType, ".") {
		return fd.fieldType
	}

	elems := strings.Split(fd.fieldType, ".")

	// import goraml or json package
	if elems[0] == "goraml" || elems[0] == "json" {
		// if property is optional, add 'pointer'
		if fd.IsOmitted {
			return "*" + fd.fieldType
		}
		return fd.fieldType
	}

	return fmt.Sprintf("%v_%v.%v", elems[0], typePackage, elems[1])
}

func (fd fieldDef) IsArray() bool {
	return strings.HasPrefix(fd.fieldType, "[]")
}

func (fd *fieldDef) buildValidators(p raml.Property) {
	validators := []string{}
	addVal := func(s string) {
		validators = append(validators, s)
	}
	// string
	if p.MinLength != nil {
		addVal(fmt.Sprintf("min=%v", *p.MinLength))
	}
	if p.MaxLength != nil {
		addVal(fmt.Sprintf("max=%v", *p.MaxLength))
	}
	if p.Pattern != nil {
		// `\` need to be escaped with 2 backslashes `\\`
		pattern := strings.ReplaceAll(*p.Pattern, "\\", "\\\\")
		// Commas need to be escaped with 2 backslashes `\\`.
		pattern = strings.ReplaceAll(pattern, ",", "\\\\,")
		addVal(fmt.Sprintf("regexp=%v", pattern))
	}

	// Number
	if p.Minimum != nil {
		addVal(fmt.Sprintf("min=%v", *p.Minimum))
	}

	if p.Maximum != nil {
		addVal(fmt.Sprintf("max=%v", *p.Maximum))
	}

	if p.MultipleOf != nil {
		addVal(fmt.Sprintf("multipleOf=%v", *p.MultipleOf))
	}

	//if p.Format != nil {
	//}

	// Array & Map
	if p.MinItems != nil {
		addVal(fmt.Sprintf("min=%v", *p.MinItems))
	}
	if p.MaxItems != nil {
		addVal(fmt.Sprintf("max=%v", *p.MaxItems))
	}
	if p.UniqueItems {
		fd.UniqueItems = true
	}

	// Required
	if !fd.IsOmitted && fd.fieldType != "bool" {
		addVal("nonzero")
	}

	fd.Validators = strings.Join(validators, ",")
}

// format struct's field name
// - Title it
// - replace '-' with camel case version
func formatFieldName(name string) string {
	return utils.Camelize(name)
}
