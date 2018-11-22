package gorestful

import (
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/raml"
)

const (
	typeDir     = "types" // directory of the generated types
	typePackage = "types" // package name of the generated types
)

var (
	dateMap = map[string]string{
		"date":          "Date",
		"date-only":     "DateOnly",
		"time-only":     "TimeOnly",
		"datetime-only": "DatetimeOnly",
		"datetime":      "DateTime",
	}

	typeMap = map[string]string{
		"string":  "string",
		"file":    "string",
		"number":  "float64",
		"integer": "int",
		"boolean": "bool",
		"object":  "json.RawMessage",
		"int8":    "int8",
		"int16":   "int16",
		"int32":   "int32",
		"int64":   "int64",
		"int":     "int",
		"long":    "int64",
		"float":   "float64",
		"double":  "float64",
	}
)

func mapDate(tip string) string {
	v, ok := dateMap[tip]
	if !ok {
		return ""
	}
	return goramlPkgDir() + v
}

// returns true if `tip` is:
// - builtin go type
// - builtin goraml type
func isBuiltinOrGoramlType(tip string) bool {
	switch {
	case strings.HasPrefix(tip, "[][]"):
		tip = strings.TrimPrefix(tip, "[][]")
	case strings.HasPrefix(tip, "[]"):
		tip = strings.TrimPrefix(tip, "[]")
	}
	for _, t := range typeMap {
		if t == tip {
			return true
		}
	}
	if mapDate(tip) != "" {
		return true
	}
	return false
}

func goramlPkgDir() string {
	if globGoramlPkgDir == "" {
		return "goraml."
	}
	return globGoramlPkgDir + "."
}

// convert from raml type to go type
func convertToGoType(tip, items string) string {
	// check for raml builtin type
	if v, ok := typeMap[tip]; ok {
		return v
	}

	if v := mapDate(tip); v != "" {
		return v
	}

	// creates raml.Type object for this type,
	// so we can use raml.Type methods
	ramlType := raml.Type{
		Type:  tip,
		Items: items,
	}

	// other types that need some processing
	switch {
	case ramlType.IsBidimensiArray():
		return "[][]" + ramlType.BidimensiArrayType()
	case ramlType.IsArray():
		return "[]" + convertToGoType(ramlType.ArrayType(), "")
	case ramlType.IsUnion():
		return unionNewName(tip)
	case ramlType.IsMultipleInheritance():
		parents, _ := ramlType.MultipleInheritance()
		return multipleInheritanceNewName(parents)
	}
	return commons.NormalizePkgName(tip)
}
