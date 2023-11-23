package golang

import (
	"fmt"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
)

// FieldDef defines a field of a struct
type fieldDef struct {
	Name          string // field name
	fieldType     string // field type
	IsComposition bool   // composition type
	IsOmitted     bool   // omitted empty
	UniqueItems   bool
	Enum          *enum // not nil if this field contains enum

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
	if strings.Index(fd.fieldType, ".") < 0 {
		return fd.fieldType
	}

	elems := strings.Split(fd.fieldType, ".")

	// import goraml or json package
	if elems[0] == "goraml" || elems[0] == "json" {
		return fd.fieldType
	}

	return fmt.Sprintf("%v_%v.%v", elems[0], typePackage, elems[1])
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
		addVal(fmt.Sprintf("regexp=%v", *p.Pattern))
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
	var formatted string
	for _, v := range strings.Split(name, "-") {
		formatted += strings.Title(v)
	}
	return formatted
}
