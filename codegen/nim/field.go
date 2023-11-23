package nim

import (
	"github.com/hotmall/go-raml/raml"
)

var (
	// list of nim keywords which need to escaped
	// it only listed keywords that commonly used as field name
	nimKeywords = map[string]bool{
		"type":   true,
		"method": true,
	}
)

// field represents a Nim object field
type field struct {
	Name string // field name
	Type string // field type
	Enum *enum
}

// EscapedName produces escaped name of Nim keyword
// or return as is if it is not Nim keyword
func (f field) EscapedName() string {
	name := f.Name
	if _, ok := nimKeywords[name]; !ok {
		return name
	}
	return "`" + name + "`"
}

func newField(objName string, prop raml.Property) field {
	f := field{
		Name: prop.Name,
		Type: toNimType(prop.TypeString(), prop.Items.Type),
	}
	if prop.IsEnum() {
		f.Enum = newEnum(objName, prop, false)
		f.Type = f.Enum.Name()
	}
	return f
}
