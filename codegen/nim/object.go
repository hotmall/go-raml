package nim

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/types"
	"github.com/hotmall/go-raml/raml"
)

var (
	// saves all generated objects
	objectsRegister map[string]adt
)

func init() {
	// register of all objects.
	// needed to generate `import`
	objectsRegister = map[string]adt{}
}

type adt interface {
	NimType() string
	Name() string
	FieldsMap() map[string]field
}

// object represents a Nim object
type object struct {
	name        string
	Description []string
	Fields      map[string]field
	T           raml.Type
	OneLineDef  string
	Parents     []string
	Enum        *enum
}

func (o object) NimType() string {
	return "object"
}

func (o object) Name() string {
	return o.name
}

func (o object) FieldsMap() map[string]field {
	return o.Fields
}

func generateAllObjects(apiDef *raml.APIDefinition, dir string) error {
	objs := []object{}

	// array of tip that need to be generated in the end of this
	// process. because it needs other object to be registered first
	delayedMI := []string{} // delayed multiple inheritance

	addObj := func(obj object) {
		objs = append(objs, obj)
		registerObject(obj)
	}

	for name, t := range types.AllTypes(apiDef, "") {
		switch tip := t.Type.(type) {
		case string:
			rt := raml.Type{
				Type: tip,
			}
			// we currently only handle multiple inheritance.
			// TODO we also need to handle union, but we still don't have
			// union support
			if rt.IsMultipleInheritance() {
				delayedMI = append(delayedMI, tip)
			}
		case types.TypeInBody:
			suffix := commons.RespBodySuffix
			if tip.ReqResp == types.HTTPRequest {
				suffix = commons.ReqBodySuffix
			}
			verb := strings.Title(strings.ToLower(tip.Endpoint.Verb))
			bodyName := setBodyName(tip.Body(), tip.Endpoint.Addr+verb, suffix)
			obj := newObject(bodyName, "", tip.Properties)
			addObj(obj)

		case raml.Type:
			obj := newObjectFromType(tip, t.Name)
			obj.name = name

			for _, f := range obj.Fields {
				if f.Enum != nil {
					registerObject(f.Enum)
				}
			}
			addObj(obj)
		}
	}

	for _, tip := range delayedMI {
		rt := raml.Type{
			Type: tip,
		}
		if parents, isMult := rt.MultipleInheritance(); isMult {
			obj := newObject(multipleInheritanceNewName(parents), "",
				map[string]interface{}{})
			obj.inherit(parents)
			addObj(obj)
		}
	}

	for _, obj := range objs {
		if err := obj.generate(dir); err != nil {
			return err
		}
	}
	return nil
}

// create new object from an RAML type
func newObjectFromType(t raml.Type, name string) object {
	obj := newObject(name, t.Description, t.Properties)
	obj.T = t
	obj.handleAdvancedType()
	return obj
}

func newObject(name, description string, properties map[string]interface{}) object {
	// generate fields from type properties
	fields := make(map[string]field)

	for k, v := range properties {
		prop := raml.ToProperty(k, v)
		fd := newField(name, prop)
		fields[fd.Name] = fd
	}

	return object{
		name:        name,
		Fields:      fields,
		Description: commons.ParseDescription(description),
	}
}

// generate nim object representation
func (o *object) generate(dir string) error {
	o.handleAdvancedType()

	// generate enums
	for _, f := range o.Fields {
		if f.Enum != nil {
			if err := f.Enum.generate(dir); err != nil {
				return err
			}
		}
	}

	if o.Enum != nil {
		return o.Enum.generate(dir)
	}
	filename := filepath.Join(dir, o.Name()+".nim")
	if err := commons.GenerateFile(o, "./templates/nim/object_nim.tmpl", "object_nim", filename, true); err != nil {
		return err
	}
	return nil
}

// get array of all imported modules:
// - `times` module if it use Time
// - other generated modules that used as field
// - parent object
func (o object) Imports() []string {
	ip := map[string]struct{}{}

	for _, f := range o.Fields {
		if name, ok := objectRegistered(f.Type); ok {
			ip[name] = struct{}{}
		}
		if f.Type == "Time" {
			ip["times"] = struct{}{}
		}
	}

	for _, p := range o.Parents {
		if name, ok := objectRegistered(p); ok {
			ip[name] = struct{}{}
		}
	}
	// del reference to our self
	if _, ok := ip[o.Name()]; ok {
		delete(ip, o.Name())
	}
	return commons.MapToSortedStrings(ip)
}

// handle RAML advanced data type
func (o *object) handleAdvancedType() {
	strType := o.T.TypeString()

	parents, isMultipleInheritance := o.T.MultipleInheritance()
	switch {
	case isMultipleInheritance: //multiple inheritance
		o.inherit(parents)
	case o.T.IsEnum():
		o.makeEnum()
	case strType == "object" || strType == "": // plain type
	case o.T.IsArray():
		o.makeArray(strType)
	}
}

func (o *object) makeEnum() {
	o.Enum = newEnumFromObject(o)
}

func (o *object) makeArray(t string) {
	o.Parents = append(o.Parents, toNimType(t, ""))
	o.buildOneLine(toNimType(t, ""))
}

func (o *object) buildOneLine(tipe string) {
	o.OneLineDef = fmt.Sprintf("%v* = %v", o.Name(), tipe)
}

// Multiple inheritance is currently not supported by Nim.
// see https://nim-lang.org/docs/tut2.html
// Golang also doesn't support it but we use composition.
// we can't use composition in Nim because it works differently
// than Go. We can't use composed type directly in Nim.
// So, we inherit the properties.
func (o *object) inherit(parents []string) {
	// save current fields
	oriFields := o.Fields

	// inherit from parents
	for _, parent := range parents {
		obj, ok := objectsRegister[parent]
		if !ok {
			fmt.Printf("parent %v not exist in object register\n", parent)
			continue
		}
		for name, f := range obj.FieldsMap() {
			o.Fields[name] = f
		}
	}
	// restore original fields
	for name, f := range oriFields {
		o.Fields[name] = f
	}
}

func registerObject(a adt) {
	name := strings.TrimSpace(a.Name())
	objectsRegister[name] = a
}

// regex to find string inside a square bracket
var reSquareBracket = regexp.MustCompile(`\[(.*?)\]`)

// check if an object named `objName` has been generated
func objectRegistered(objName string) (string, bool) {
	objName = extractTypeName(objName)
	_, ok := objectsRegister[objName]
	return objName, ok
}

// extract type name from potentially complex type declaration
func extractTypeName(tip string) string {
	// trim all spaces
	tip = strings.TrimSpace(tip)

	// detect array
	if strings.Index(tip, "]") < 0 {
		return tip
	}

	// find type inside square bracket
	inBracket := reSquareBracket.FindString(tip)
	if inBracket == "" {
		return tip
	}
	last := strings.LastIndex(inBracket, "[")
	inBracket = inBracket[last+1:]
	return strings.TrimSuffix(inBracket, "]")
}

func multipleInheritanceNewName(parents []string) string {
	return strings.Join(parents, "")
}
