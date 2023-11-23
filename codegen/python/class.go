package python

import (
	"path/filepath"
	"strings"

	"github.com/hotmall/go-raml/codegen/capnp"
	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/types"
	"github.com/hotmall/go-raml/raml"
	"github.com/pinzolo/casee"
)

// class defines a python class
type class struct {
	T           raml.Type
	name        string
	Description []string
	Fields      map[string]field
	Enum        *enum
	Capnp       bool
	AliasOf     string
	imports     map[string]struct{}
}

func (pc class) Name() string {
	return commons.NormalizeIdentifier(pc.name)
}

type objectProperty struct {
	name            string
	required        bool
	datatype        string
	childProperties []objectProperty
}

// create a python class representations
func newClass(apiDef *raml.APIDefinition, T raml.Type, name string, description string, properties map[string]interface{}, capnp bool) class {
	pc := class{
		name:        name,
		Description: commons.ParseDescription(description),
		Fields:      map[string]field{},
		T:           T,
		Capnp:       capnp,
		imports:     make(map[string]struct{}),
	}

	if pc.T.IsAlias() {
		return pc
	}

	// initialize the fields

	types := globAPIDef.Types

	typeHierarchy := getTypeHierarchy(name, T, types)
	ramlTypes := make([]raml.Type, 0)
	for _, v := range typeHierarchy {
		for _, iv := range v {
			ramlTypes = append(ramlTypes, iv)
		}
	}
	mergedProps := getTypeProperties(ramlTypes)
	for k, v := range properties {
		prop := raml.ToProperty(k, v)
		mergedProps[k] = prop
	}

	for propName, propInterface := range mergedProps {
		op := objectProperties(propName, propInterface)
		field, err := newField(name, apiDef, T, propName, propInterface, types, op, typeHierarchy)
		if err != nil {
			continue
		}

		pc.Fields[field.Name] = field
	}

	return pc
}

func newClassFromType(apiDef *raml.APIDefinition, T raml.Type, name string, capnp bool) class {
	pc := newClass(apiDef, T, name, T.Description, T.Properties, capnp)
	pc.T = T
	pc.handleAdvancedType()
	return pc
}

func objectProperties(name string, p interface{}) []objectProperty {
	props := make([]objectProperty, 0)

	switch prop := p.(type) {
	case map[interface{}]interface{}:
		ramlProp := raml.ToProperty(name, p)
		if ramlProp.Type == "object" {
			for k, v := range prop {
				switch k {
				case "properties":
					for propName, childProp := range v.(map[interface{}]interface{}) {
						rProp := raml.ToProperty(propName.(string), childProp)
						objprop := objectProperty{
							name:     rProp.Name,
							required: rProp.Required,
							datatype: rProp.TypeString(),
						}
						if rProp.Type == "object" {
							objprop.childProperties = append(objprop.childProperties, objectProperties(propName.(string), childProp)...)
						}
						props = append(props, objprop)
					}
				}
			}
		}
	}

	return props
}

// ChildProperties returns child properties of a property
func ChildProperties(Properties map[string]interface{}) []raml.Property {
	props := make([]raml.Property, 0)

	for propName, propInterface := range Properties {
		props = append(props, raml.ToProperty(propName, propInterface))
	}

	return props
}

func getTypeHierarchy(name string, T raml.Type, types map[string]raml.Type) []map[string]raml.Type {
	typelist := []map[string]raml.Type{map[string]raml.Type{name: T}}

	for _, parent := range T.Parents() {
		parentType, inherited := types[parent]
		if inherited {
			for _, pt := range getTypeHierarchy(parent, parentType, types) {
				typelist = append(typelist, pt)
			}
		}
	}

	return typelist
}

func getTypeProperties(typelist []raml.Type) map[string]raml.Property {
	// get a list of the types in the inheritance chain for T
	// walk it from the top down and add the properties
	properties := map[string]raml.Property{}
	for i := len(typelist) - 1; i >= 0; i-- {
		for k, v := range typelist[i].Properties {
			prop := raml.ToProperty(k, v)
			// we convert it to property here
			// because we need the proper name, sometimes the name has "?" suffix
			//   which mean optional properties
			properties[prop.Name] = prop
		}
	}

	return properties
}

// generate a python class file
func (pc *class) generate(dir, template, name string) ([]string, error) {
	fileName := filepath.Join(dir, pc.Name()+".py")

	if pc.AliasOf != "" {
		return []string{pc.Name()}, commons.GenerateFile(pc, "./templates/python/class_alias.tmpl",
			"class_alias", fileName, true)
	}

	// generate enums
	typeNames := make([]string, 0)
	for _, f := range pc.Fields {
		if f.Enum != nil {
			typeNames = append(typeNames, f.Enum.Name)
			if err := f.Enum.generate(dir); err != nil {
				return typeNames, err
			}
		}
	}

	if pc.Enum != nil {
		typeNames = append(typeNames, pc.Enum.Name)
		return typeNames, pc.Enum.generate(dir)
	}

	typeNames = append(typeNames, pc.Name())
	return typeNames, commons.GenerateFile(pc, template, name, fileName, true)
}

func (pc *class) handleAdvancedType() {
	if pc.T.Type == nil {
		pc.T.Type = "object"
	}
	switch {
	case pc.T.IsEnum():
		pc.Enum = newEnumFromClass(pc)
	case pc.T.IsAlias():
		pc.createTypeAlias()
	}
}

func (pc *class) createTypeAlias() {
	typeStr := pc.T.TypeString()

	pt := toPythonType(typeStr)
	if pt != nil { // from builtin type
		pc.AliasOf = pt.name
		if pt.importName != "" {
			pc.addImport(pt.importModule, pt.importName)
		}
		return
	}

	// from another type
	pc.AliasOf = typeStr
	pc.addImport(".", typeStr)
}

// return list of import statements
func (pc *class) Imports() []string {
	for _, field := range pc.Fields {
		for _, imp := range field.imports {
			// do not import ourself
			if imp.Name == pc.Name() {
				continue
			}
			pc.addImport(imp.Module, imp.Name)
		}
	}

	return commons.MapToSortedStrings(pc.imports)
}

func (pc *class) addImport(mod, name string) {
	importString := "from " + mod + " import " + name
	pc.imports[importString] = struct{}{}
}
func (pc class) CapnpName() string {
	return casee.ToPascalCase(pc.Name())
}

// generate all python classes from a RAML document
func GenerateAllClasses(apiDef *raml.APIDefinition, dir string, capnp bool) ([]string, error) {
	// array of tip that need to be generated in the end of this
	// process. because it needs other object to be registered first
	delayedMI := []string{} // delayed multiple inheritance
	template := "./templates/python/class_python.tmpl"
	templateName := "class_python"

	ci := classInit{
		Types: make(map[string][]string),
	}

	names := []string{}
	for name, t := range types.AllTypes(apiDef, "") {
		var errGen error
		var results []string
		switch tip := t.Type.(type) {
		case string:
			rt := raml.Type{
				Type: tip,
			}
			if rt.IsMultipleInheritance() {
				delayedMI = append(delayedMI, tip)
			}
		case types.TypeInBody:
			newTipName := types.PascalCaseTypeName(tip)
			pc := newClass(apiDef, raml.Type{Type: "object"}, newTipName, "", tip.Properties, capnp)
			results, errGen = pc.generate(dir, template, templateName)
		case raml.Type:
			if name == "UUID" {
				continue
			}
			pc := newClassFromType(apiDef, tip, name, capnp)
			results, errGen = pc.generate(dir, template, templateName)
		}

		if errGen != nil {
			return names, errGen
		}

		if len(results) > 0 {
			className := results[len(results)-1]
			ci.append(className, results)
		}
		names = append(names, results...)
	}

	for _, tip := range delayedMI {
		rt := raml.Type{
			Type: tip,
		}
		if parents, isMult := rt.MultipleInheritance(); isMult {
			pc := newClassFromType(apiDef, rt, strings.Join(parents, ""), capnp)
			results, err := pc.generate(dir, template, templateName)
			if err != nil {
				return names, err
			}
			if len(results) > 0 {
				className := results[len(results)-1]
				ci.append(className, results)
			}
			names = append(names, results...)
		}
	}

	// generate pytypes/__init__.py
	ci.generate(dir)
	return names, nil
}

// GeneratePythonCapnpClasses generates python classes from a raml definition along with function to load binaries from/to capnp
// and generates the needed capnp schemas
func GeneratePythonCapnpClasses(apiDef *raml.APIDefinition, dir string) error {
	// TODO : get rid of this global variables
	globAPIDef = apiDef

	if err := capnp.GenerateCapnp(apiDef, dir, "", ""); err != nil {
		return err
	}

	_, err := GenerateAllClasses(apiDef, dir, true)
	return err
}
