package commons

import (
	"fmt"
	"strings"

	"github.com/hotmall/go-raml/raml"
)

// IsBuiltinType returns true if the given type
// is builtin type
func IsBuiltinType(t interface{}) bool {
	tip := raml.Type{
		Type: t,
	}
	return tip.IsBuiltin()
}

// GetBasicType returns basic type of a basic or complex type
func GetBasicType(t string) string {
	tip := raml.Type{
		Type: t,
	}

	switch {
	case tip.IsArray():
		return tip.ArrayType()
	case tip.IsBidimensiArray():
		return tip.BidimensiArrayType()
	default:
		return t
	}
}

// IsArrayType returns true if the given
// string is a RAML array
func IsArrayType(t string) bool {
	tip := raml.Type{
		Type: t,
	}
	return tip.IsArray()
}

// CheckDuplicatedTitleTypes checks that the API definition has duplicated types
// of the title case version of the types.
// title case = first letter become uppercase
// example of non duplicate:
//  - One & Two
//  - One & oNe = One & ONe -> N is uppercase
// example of duplicate
// - One & one = One & One
func CheckDuplicatedTitleTypes(apiDef *raml.APIDefinition) error {
	var title string

	for name := range apiDef.Types {
		title = strings.Title(name)
		if title == name {
			continue
		}
		if _, duplicate := apiDef.Types[title]; duplicate {
			return fmt.Errorf("types conflict: %s with %v", name, title)
		}
	}
	return nil
}
