package capnp

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"

	"github.com/pinzolo/casee"
)

type Struct struct {
	ID            string
	Name          string
	Fields        map[string]field
	Description   []string
	OrderedFields []field
	T             raml.Type
	lang          string
	pkg           string
	Enum          *enum
}

func NewStruct(t raml.Type, name, lang, pkg string) (Struct, error) {
	name = casee.ToPascalCase(name)

	// generate fields from type properties
	fields := make(map[string]field)
	fieldNames := []string{}

	for k, v := range t.Properties {
		fd := newField(name, raml.ToProperty(k, v), lang, pkg)
		fields[fd.Name] = fd
		fieldNames = append(fieldNames, fd.Name)
	}

	s := Struct{
		ID:          getID(),
		Name:        name,
		Fields:      fields,
		Description: commons.ParseDescription(t.Description),
		T:           t,
		pkg:         pkg,
		lang:        lang,
	}
	if s.T.Enum != nil {
		s.Enum = newEnumFromType(name, s.T, lang, pkg)
	}
	if err := s.checkValidCapnp(); err != nil {
		return s, err
	}

	s.orderFields(fieldNames)

	return s, nil
}

// Generate generates struct code
func (s *Struct) Generate(dir string) error {
	if s.Enum != nil {
		return s.Enum.generate(dir)
	}

	if err := s.generateEnums(dir); err != nil {
		return err
	}
	filename := filepath.Join(dir, s.Name+".capnp")
	return commons.GenerateFile(s, "./templates/capnp/struct_capnp.tmpl", "struct_capnp", filename, true)
}

// generate all enums contained in this struct
func (s *Struct) generateEnums(dir string) error {
	for _, f := range s.Fields {
		if f.Enum == nil {
			continue
		}
		if err := f.Enum.generate(dir); err != nil {
			return err
		}
	}
	return nil
}

func (s *Struct) Imports() []string {
	imports := s.importsNonBuiltin()
	switch s.lang {
	case "go":
		imports = append(imports, s.goImports()...)
	default:
	}
	sort.Strings(imports)
	return imports
}

func (s *Struct) importsNonBuiltin() []string {
	imports := []string{}
	for _, f := range s.Fields {
		// import non buitin types
		if typesRegistered(f.Type) {
			imports = append(imports, fmt.Sprintf(`using import "%v.capnp".%v`, f.Type, f.Type))
		}
		// import enum types
		if f.Enum != nil {
			imports = append(imports, fmt.Sprintf(`using import "%v.capnp".%v`, f.Enum.Name, f.Enum.Name))
		}
		// import non-builtin types used in List
		if typesRegistered(f.Items) {
			imports = append(imports, fmt.Sprintf(`using import "%v.capnp".%v`, f.Items, f.Items))
		}
	}
	return imports
}
func (s *Struct) Annotations() []string {
	switch s.lang {
	case "go":
		return s.goAnnotations()
	default:
		return []string{}
	}
}

func (s *Struct) checkValidCapnp() error {
	if !casee.IsPascalCase(s.Name) {
		return fmt.Errorf("invalid type name:%v. Type names must be PascalCase", s.Name)
	}

	for _, field := range s.Fields {
		if !casee.IsCamelCase(field.Name) {
			return fmt.Errorf("invalid decleration name:%v. Decleration names must be CamelCase", field.Name)
		}
	}

	return nil
}

func (s *Struct) orderFields(fieldNames []string) {
	sort.Strings(fieldNames)
	for index, name := range fieldNames {
		field := s.Fields[name]
		field.Num = index
		s.OrderedFields = append(s.OrderedFields, field)
	}
}
