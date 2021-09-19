package gorestful

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/raml"
)

type paramDef fieldDef

func newParamDef(param *raml.NamedParameter) paramDef {
	var (
		paramType = param.Type                      // the field type
		basicType = commons.GetBasicType(paramType) // basic type of the field type
	)

	paramType = convertToGoType(paramType, basicType)

	pd := paramDef{
		Name:      formatFieldName(param.Name),
		fieldType: paramType,
		IsOmitted: !param.Required,
		Default:   param.Default,
	}

	pd.buildValidators(param)

	return pd
}

func (pd paramDef) Type() string {
	// doesn't have "." -> doesnt import from other package
	if !strings.Contains(pd.fieldType, ".") {
		return pd.fieldType
	}

	elems := strings.Split(pd.fieldType, ".")

	// import goraml or json package
	if elems[0] == "goraml" || elems[0] == "json" {
		return pd.fieldType
	}

	return fmt.Sprintf("%v_%v.%v", elems[0], typePackage, elems[1])
}

func (pd *paramDef) buildValidators(p *raml.NamedParameter) {
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
		// `\` need to be escaped with 2 backslashes `\\`
		pattern := strings.ReplaceAll(*p.Pattern, "\\", "\\\\")
		// Commas need to be escaped with 2 backslashes `\\`.
		pattern = strings.ReplaceAll(pattern, ",", "\\\\,")
		addVal(fmt.Sprintf("regexp=%v", pattern))
	}

	// Number
	if p.Minimum != nil {
		addVal(fmt.Sprintf("min=%v", *p.Minimum))
	}

	if p.Maximum != nil {
		addVal(fmt.Sprintf("max=%v", *p.Maximum))
	}

	// Required
	if !pd.IsOmitted && pd.fieldType != "bool" {
		addVal("nonzero")
	}

	pd.Validators = strings.Join(validators, ",")
}

func copyURIParameters(context map[string]raml.NamedParameter, resource *raml.Resource) {

	for name, parameter := range resource.URIParameters {
		context[name] = parameter
	}

	if resource.Parent != nil {
		copyURIParameters(context, resource.Parent)
	}
}

func copyHeaders(headers map[string]raml.NamedParameter, m *method) {
	if m == nil {
		return
	}

	for key, value := range m.Headers {
		name := fmt.Sprintf("%s", key)
		headers[name] = raml.NamedParameter{
			Name:        value.Name,
			DisplayName: value.DisplayName,
			Type:        value.Type,
			Pattern:     value.Pattern,
			MinLength:   value.MinLength,
			MaxLength:   value.MaxLength,
			Minimum:     value.Minimum,
			Maximum:     value.Maximum,
			Repeat:      value.Repeat,
			Required:    value.Required,
		}
	}
}

type goContext struct {
	Name        string              // context's name
	PackageName string              // package name
	Fields      map[string]paramDef // all context's fields
}

func newGoContext(m *method) *goContext {
	parameters := make(map[string]raml.NamedParameter)

	copyURIParameters(parameters, m.Resource)

	for name, param := range m.QueryParameters {
		parameters[name] = param
	}

	copyHeaders(parameters, m)

	fields := make(map[string]paramDef)

	for name, param := range parameters {
		param.Name = name
		fields[name] = newParamDef(&param)
	}

	methodName := m.Name
	if len(m.DisplayName) > 0 {
		methodName = strings.Title(commons.DisplayNameToFuncName(m.DisplayName))
	}

	return &goContext{
		Name:        methodName,
		PackageName: "types",
		Fields:      fields,
	}
}

func (gc goContext) generate(dir string) error {
	filename := filepath.Join(dir+"/types/", gc.Name+"Context.go")
	if err := commons.GenerateFile(gc, "./templates/golang/gorestful_context.tmpl",
		"gorestful_context_template", filename, true); err != nil {
		return err
	}
	return nil
}

func (gc goContext) generateClient(dir string) error {
	filename := filepath.Join(dir+"/types/", gc.Name+"Context.go")
	if err := commons.GenerateFile(gc, "./templates/golang/requests_client_context.tmpl",
		"requests_client_context", filename, true); err != nil {
		return err
	}
	return nil
}
