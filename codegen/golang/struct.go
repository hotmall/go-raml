package golang

import (
	"path/filepath"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/types"
	"github.com/hotmall/go-raml/raml"
)

const (
	structTemplateLocation         = "./templates/golang/struct.tmpl"
	inputValidatorTemplateLocation = "./templates/golang/struct_input_validator.tmpl"
	inputValidatorFileResult       = "struct_input_validator.go"
)

// StructDef defines a struct
type structDef struct {
	T           raml.Type           // raml.Type of this struct
	Name        string              // struct's name
	Description []string            // structs description
	PackageName string              // package name
	Fields      map[string]fieldDef // all struct's fields
	OneLineDef  string              // not empty if this struct can be defined in one line
	Enum        *enum

	Validators []string
}

// true if this struct is not an alias of `interface{}`
func (sd structDef) NotBareInterface() bool {
	return !strings.HasSuffix(sd.OneLineDef, " interface{}")
}

// create new struct definition
func newStructDef(apiDef *raml.APIDefinition, name, packageName, description string, properties map[string]interface{}) structDef {
	// generate struct's fields from type properties
	fields := make(map[string]fieldDef)
	for k, v := range properties {
		prop := raml.ToProperty(k, v)
		fields[prop.Name] = newFieldDef(apiDef, name, prop, packageName)
	}
	return structDef{
		Name:        strings.Title(commons.NormalizeIdentifier(name)),
		PackageName: packageName,
		Fields:      fields,
		Description: commons.ParseDescription(description),
	}
}

// create struct definition from RAML Type node
func newStructDefFromType(apiDef *raml.APIDefinition, t raml.Type, sName, packageName string) structDef {
	sd := newStructDef(apiDef, sName, packageName, t.Description, t.Properties)
	sd.T = t

	// handle advanced type on raml1.0
	sd.handleAdvancedType(apiDef)

	return sd
}

// generate Go struct
func (sd structDef) generate(dir string) error {
	// generate enums
	for _, f := range sd.Fields {
		if f.Enum != nil {
			if err := f.Enum.generate(dir); err != nil {
				return err
			}
		}
	}
	if sd.Enum != nil {
		return sd.Enum.generate(dir)
	}
	fileName := filepath.Join(dir, sd.Name+".go")
	return commons.GenerateFile(sd, structTemplateLocation, "struct_template", fileName, true)
}

func generateStructs(types map[string]raml.Type, dir string) error {
	apiDef := raml.APIDefinition{
		Types: types,
	}
	return generateAllStructs(&apiDef, dir)
}

func generateAllStructs(apiDef *raml.APIDefinition, dir string) error {
	pkgName := typePackage
	dir = filepath.Join(dir, typeDir)
	for _, t := range types.AllTypes(apiDef, pkgName) {
		switch tip := t.Type.(type) {
		case string:
			createGenerateStruct(apiDef, tip, dir, pkgName)
		case types.TypeInBody:
			// TODO:
			// need to change this Req/Resp body name
			// to simply use methodName
			// it will break compatibility!
			name := commons.NormalizeURITitle(tip.Endpoint.Addr)
			name = name[len(tip.Endpoint.ResourceName()):] + strings.Title(strings.ToLower(tip.Endpoint.Verb))
			if tip.ReqResp == types.HTTPRequest {
				name += "ReqBody"
			} else {
				name += "RespBody"
			}
			sd := newStructDef(apiDef, name, pkgName, tip.Description, tip.Properties)
			if err := sd.generate(dir); err != nil {
				return err
			}
		case raml.Type:
			sd := newStructDefFromType(apiDef, tip, t.Name, pkgName)
			if err := sd.generate(dir); err != nil {
				return err
			}
		}
	}
	return nil
}

// Imports returns all packages that
// need to be imported by this struct
func (sd structDef) Imports() []string {
	ip := map[string]struct{}{}

	if sd.needFmt() {
		ip[`"fmt"`] = struct{}{}
	}
	if sd.OneLineDef == "" {
		ip[`"gopkg.in/validator.v2"`] = struct{}{}
	}

	// libraries
	for _, fd := range sd.Fields {
		if fd.fieldType == "json.RawMessage" {
			ip[`"encoding/json"`] = struct{}{}
		} else if lib := libImportPath(globRootImportPath, fd.fieldType, globLibRootURLs); lib != "" {
			ip[lib] = struct{}{}
		}
	}
	return commons.MapToSortedStrings(ip)
}

// handle advance type type into structField
// example:
//   Mammal:
//     type: Animal
//     properties:
//       name:
//         type: string
// the additional fieldDef would be Animal composition
func (sd *structDef) handleAdvancedType(apiDef *raml.APIDefinition) {
	if sd.T.Type == nil {
		sd.T.Type = "object"
	}

	parents, isMultipleInherit := sd.T.MultipleInheritance()
	parent, isSingleInherit := sd.T.SingleInheritance()

	switch {
	case isMultipleInherit:
		sd.addMultipleInheritance(apiDef, parents)
	case sd.T.IsUnion():
		sd.buildUnion()
	case sd.T.IsArray():
		sd.buildArray()
	case strings.ToLower(sd.T.TypeString()) == "object": // plain type
		return
	case sd.T.IsEnum():
		sd.buildEnum()
	case sd.T.IsAlias():
		sd.buildTypeAlias(apiDef)
	case isSingleInherit:
		sd.addSingleInheritance(apiDef, parent)
	}
}

// add single inheritance
// inheritance is implemented as composition
// spec : http://docs.raml.org/specs/1.0/#raml-10-spec-inheritance-and-specialization
func (sd *structDef) addSingleInheritance(apiDef *raml.APIDefinition, strType string) {
	// check if parent is user defined type
	if _, ok := apiDef.Types[strType]; ok {
		strType = strings.Title(strType)
	}

	fd := fieldDef{
		Name:          strType,
		IsComposition: true,
	}
	sd.Fields[strType] = fd

}

// construct multiple inheritance to Go type
// example :
// Anggora:
//	 type: [ Animal , Cat ]
//	 properties:
//		color:
//			type: string
// The additional fielddef would be a composition of Animal & Cat
// http://docs.raml.org/specs/1.0/#raml-10-spec-multiple-inheritance
func (sd *structDef) addMultipleInheritance(apiDef *raml.APIDefinition, parents []string) {
	for _, s := range parents {
		fieldType := strings.TrimSpace(s)

		// check if parent is user defined type
		if _, ok := apiDef.Types[fieldType]; ok {
			fieldType = strings.Title(fieldType)
		}

		fd := fieldDef{
			Name:          fieldType,
			IsComposition: true,
		}

		sd.Fields[fd.Name] = fd
	}
}

// buildEnum based on http://docs.raml.org/specs/1.0/#raml-10-spec-enums
// example result  `type TypeName []data_type`
func (sd *structDef) buildEnum() {
	sd.Enum = newEnumFromStruct(sd)
}

// build array type
// spec http://docs.raml.org/specs/1.0/#raml-10-spec-array-types
// example result  `type TypeName []something`
func (sd *structDef) buildArray() {
	sd.buildOneLine(convertToGoType(sd.T.Type.(string), ""))
}

// build union type
// union type is implemented as empty struct
func (sd *structDef) buildUnion() {
}

func (sd *structDef) buildTypeAlias(apiDef *raml.APIDefinition) {
	strType := sd.T.TypeString()
	// check if parent is user defined type
	if _, ok := apiDef.Types[strType]; ok {
		strType = strings.Title(strType)
	}

	sd.buildOneLine(convertToGoType(strType, ""))
}

func (sd *structDef) buildOneLine(tipe string) {
	sd.OneLineDef = "type " + sd.Name + " " + tipe
}

// generate input validator helper file
func generateInputValidator(packageName, dir string) error {
	var ctx = struct {
		PackageName string
	}{
		PackageName: packageName,
	}
	fileName := filepath.Join(dir, inputValidatorFileResult)
	return commons.GenerateFile(ctx, inputValidatorTemplateLocation, "struct_input_validator_template", fileName, true)
}

// true if this struct need to import 'fmt' package
// It is required by validation code,
// because validation error will need `fmt` to build error message
func (sd structDef) needFmt() bool {
	// array type min items and max items
	if sd.T.MinItems > 0 || sd.T.MaxItems > 0 {
		return true
	}

	// unique items
	for _, f := range sd.Fields {
		if f.UniqueItems {
			return true
		}
	}
	return false
}

func multipleInheritanceNewName(parents []string) string {
	return strings.Join(parents, "")
}

func unionNewName(tip string) string {
	tipes := strings.Split(tip, "|")
	for i, v := range tipes {
		tipes[i] = strings.TrimSpace(v)
	}
	return "Union" + strings.Join(tipes, "")
}

// create struct and generate it if possible.
// return:
// - newType Name if we try to generate it
// - nil if no error happened during generation
func createGenerateStruct(apiDef *raml.APIDefinition, tip, dir, pkgName string) (string, error) {
	rt := raml.Type{
		Type: tip,
	}
	if parents, isMultiple := rt.MultipleInheritance(); isMultiple {
		sd := newStructDef(apiDef, multipleInheritanceNewName(parents), pkgName, "", map[string]interface{}{})
		sd.addMultipleInheritance(apiDef, parents)
		return sd.Name, sd.generate(dir)
	}

	if rt.IsUnion() {
		t := raml.Type{
			Type: tip,
		}
		sd := newStructDef(apiDef, unionNewName(tip), pkgName, "", map[string]interface{}{})
		sd.T = t
		sd.buildUnion()
		sd.generate(dir)
	}
	return "", nil
}
