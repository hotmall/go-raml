package python

import (
	"fmt"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
)

// import package and name
type pyimport struct {
	Module string
	Name   string
}

// pythons class's field
type field struct {
	Name                    string
	Type                    string
	Required                bool   // if the field itself is required
	DataType                string // the python datatype (objmap) used in the template
	HasChildProperties      bool
	RequiredChildProperties []string
	Validators              string
	Enum                    *enum
	isFormField             bool
	imports                 []pyimport
	UnionTypes              []string
	// Initializer string
	IsList     bool                // it is a list field
	validators map[string][]string // array of validators, only used to build `Validators` field
}

func newField(className string, apiDef *raml.APIDefinition, T raml.Type, propName string, propInterface interface{},
	types map[string]raml.Type, childProperties []objectProperty,
	typeHierarchy []map[string]raml.Type) (field, error) {

	prop := raml.ToProperty(propName, propInterface)

	f := field{
		Name:     commons.NormalizeIdentifier(prop.Name),
		Required: prop.Required,
	}

	if prop.IsEnum() {
		// see if T actually has this property. if not, it's inherited, and we want to only define an Enum once
		// first, get the name of the type that actually defines this property, in the inheritance chain:
		typeDefiningProp := func() string {
			for _, typeMap := range typeHierarchy {
				for typeName, typeVal := range typeMap {
					for k, v := range typeVal.Properties {
						tempProp := raml.ToProperty(k, v)
						if tempProp.Name == prop.Name {
							return typeName
						}
					}
				}
			}
			return ""
		}()

		if typeDefiningProp != className {
			// this enum property isn't actually defined on T. it's inherited from a parent type
			// thus, we'll use the parent type's enum

			f.Enum = newEnum(typeDefiningProp, prop, false)
		}
		if f.Enum == nil {
			f.Enum = newEnum(className, prop, false)
		}
		f.Type = f.Enum.Name
		f.addImport(apiDef, "."+f.Type, f.Type)
	} else {
		f.setType(apiDef, prop.TypeString(), prop.Items.Type)
		if f.Type == "" {
			return f, fmt.Errorf("unsupported type:%v", prop.Type)
		}
	}

	f.DataType, f.HasChildProperties = buildDataType(apiDef, f, childProperties)

	// I don't really understand why we need childRequired and mainRequired here.
	// it is from the original code written by @razor-1
	// TODO: remove it if possible

	// see if there are different required properties for this instance of a type vs. the type's main declaration
	mainRequired := make([]string, 0)
	childRequired := make([]string, 0)
	if mainType, ok := types[f.Type]; ok {
		switch thisProp := propInterface.(type) {
		case map[interface{}]interface{}:
			if myChildProperties, ok := thisProp["properties"].(map[interface{}]interface{}); ok {
				for _, typeProp := range ChildProperties(mainType.Properties) {
					if typeProp.Required {
						mainRequired = append(mainRequired, typeProp.Name)
					}
				}

				myChildPropertyMap := make(map[string]interface{})
				for k, v := range myChildProperties {
					if childPropName, ok := k.(string); ok {
						myChildPropertyMap[childPropName] = v
					}
				}

				for _, myProp := range ChildProperties(myChildPropertyMap) {
					if myProp.Required {
						childRequired = append(childRequired, myProp.Name)
					}
				}
			}
		}
	}
	if len(childRequired) > len(mainRequired) {
		// some properties were made required and we need to validate them
		// sort the lists so we can get only the fields that are required on this child
		sort.Strings(childRequired)
		sort.Strings(mainRequired)

		f.RequiredChildProperties = childRequired[len(mainRequired):]
	}

	return f, nil
}

func buildDataType(apiDef *raml.APIDefinition, f field, childProperties []objectProperty) (string, bool) {
	/*
		build a string for the 'datatype' key of an objmap for this property
		a complete objmap looks like:
		{'attrname': {'datatype': [type], 'required': bool}}

		the type values in the 'datatype' list can be any type, but if they are a dict, it's in objmap format

		there can be many levels of nesting, but here is an example of one:
		{
			'sipSessionId': {
				'datatype': [
					{
						'local': {'datatype': [str], 'required': False},
						'remote': {'datatype': [str], 'required': False},
					},
				],
				'required': False
			}
		}
	*/

	if len(f.UnionTypes) > 0 {
		dataType := strings.Join(f.UnionTypes, ", ")
		return dataType, false
	}
	if f.Type != "dict" || len(childProperties) == 0 {
		return f.Type, false
	}

	// @TODO: I think the part will never be reached. The type will only be dict if
	// it doesnt have properties which means no children either. Sarah
	// we have a dict with child properties of type 'object'. build the datatype string
	// fmt.Println("childprops for", f.Type, childProperties)
	var datatypes []string
	for _, objProp := range childProperties {
		// fmt.Println("have childprop", objProp)
		reqstr := "True"
		if !objProp.required {
			reqstr = "False"
		}
		childField := field{
			Name: objProp.name,
		}
		childField.setType(apiDef, objProp.datatype, "")
		thisDatatype := childField.Type
		if len(objProp.childProperties) > 0 {
			thisDatatype, _ = buildDataType(apiDef, childField, objProp.childProperties)
		}
		thisProp := fmt.Sprintf("'%s': {'datatype': [%s], 'required': %s}", objProp.name, thisDatatype, reqstr)
		datatypes = append(datatypes, thisProp)
	}

	return strings.Join(datatypes, ", "), true
}

func (pf *field) addImport(apiDef *raml.APIDefinition, module, name string) {
	_, userDefined := apiDef.Types[name]
	if !userDefined && commons.IsBuiltinType(name) && !datetimeVariant(name) {
		// datetime variant is not builtin type in python
		return
	}
	imp := pyimport{
		Module: module,
		Name:   name,
	}
	pf.imports = append(pf.imports, imp)
}

// convert from raml Type to python type
func (pf *field) setType(apiDef *raml.APIDefinition, t, items string) {
	pt := toPythonType(t)
	if pt != nil {
		pf.Type = pt.name
		if pt.importName != "" {
			pf.addImport(apiDef, pt.importModule, pt.importName)
		}
		return // type already set, no need to go down
	}

	ramlType := raml.Type{
		Type:  t,
		Items: items,
	}
	// other types that need some processing
	switch {
	case ramlType.IsBidimensiArray(): // bidimensional array
		log.Info("validator has no support for bidimensional array, ignore it")
	case ramlType.IsArray(): // array
		pf.IsList = true
		pf.setType(apiDef, ramlType.ArrayType(), "")
	case strings.HasSuffix(t, "{}"): // map
		log.Info("validator has no support for map, ignore it")
	case ramlType.IsUnion():
		// send the list of union types to the template
		unionTypes, _ := ramlType.Union()
		for _, typename := range unionTypes {
			if v, ok := typeMap[typename]; ok {
				switch typename {
				case "datetime":
					pf.addImport(apiDef, "datetime", "datetime")
				case "uuid":
					pf.addImport(apiDef, "uuid", "UUID")
				case "string":
					pf.addImport(apiDef, "six", "string_types")
				default:
					pf.addImport(apiDef, "."+typename, typename)
				}
				pf.UnionTypes = append(pf.UnionTypes, v)
			}

			pf.Type = t
		}
	case strings.Index(t, ".") > 1:
		pf.Type = t[strings.Index(t, ".")+1:]
	default:
		pf.Type = t
		pf.addImport(apiDef, "."+t, t)
	}
}

var (
	dateTimeVariantMap = map[string]struct{}{
		"date-only":     struct{}{},
		"time-only":     struct{}{},
		"datetime-only": struct{}{},
		"datetime":      struct{}{},
	}
)

// check if a type name is variant of datetime
func datetimeVariant(name string) bool {
	_, ok := dateTimeVariantMap[name]
	return ok
}
