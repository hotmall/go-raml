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
	}

	pd.buildValidators(param)

	return pd
}

func (pd paramDef) Type() string {
	// doesn't have "." -> doesnt import from other package
	if strings.Index(pd.fieldType, ".") < 0 {
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
		addVal(fmt.Sprintf("regexp=%v", *p.Pattern))
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

type goContext struct {
	Name        string              // context's name
	PackageName string              // package name
	Fields      map[string]paramDef // all context's fields
}

func newGoContext(method *serverMethod) *goContext {
	parameters := make(map[string]raml.NamedParameter)

	copyURIParameters(parameters, method.Resource)

	for name, param := range method.QueryParameters {
		parameters[name] = param
	}

	/*for name, header := range method.Headers {
		parameters[(string)name] := (raml.NamedParameter)header
	}*/

	fields := make(map[string]paramDef)

	for name, param := range parameters {
		param.Name = name
		fields[name] = newParamDef(&param)
	}

	methodName := method.Name
	if len(method.DisplayName) > 0 {
		methodName = strings.Title(commons.DisplayNameToFuncName(method.DisplayName))
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
