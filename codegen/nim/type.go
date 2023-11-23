package nim

import (
	"github.com/hotmall/go-raml/raml"
)

var (
	typeMap = map[string]string{
		"string":    "string",
		"file":      "string",
		"number":    "float",
		"integer":   "int",
		"boolean":   "bool",
		"datetime":  "Time",
		"date-only": "Time",
		"time-only": "Time",
		"int8":      "int8",
		"int16":     "int16",
		"int32":     "int32",
		"int64":     "int64",
		"int":       "int",
		"long":      "int64",
		"float":     "float",
		"double":    "float64",
	}
)

func toNimType(t, items string) string {
	if v, ok := typeMap[t]; ok {
		return v
	}
	// other types that need some processing
	ramlType := raml.Type{
		Type:  t,
		Items: items,
	}
	switch {
	case ramlType.IsBidimensiArray():
		return "seq[seq[" + toNimType(ramlType.BidimensiArrayType(), "") + "]]"
	case ramlType.IsArray():
		return "seq[" + toNimType(ramlType.ArrayType(), "") + "]"
	case ramlType.IsMultipleInheritance():
		parents, _ := ramlType.MultipleInheritance()
		return multipleInheritanceNewName(parents)
	}
	return t
}
