package gorestful

import (
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/resource"
	"github.com/Jumpscale/go-raml/codegen/security"
	"github.com/Jumpscale/go-raml/raml"
)

const (
	muxCatchAllRoute = "{path:.*}"
)

type serverMethod struct {
	*method
	Middlewares string
	SecuredBy   []raml.DefinitionChoice
}

func newServerMethod(resMeth resource.Method, apiDef *raml.APIDefinition, rd *resource.Resource) serverMethod {
	goMeth := newMethod(resMeth)

	meth := serverMethod{
		method:    goMeth,
		SecuredBy: security.GetMethodSecuredBy(apiDef, resMeth.Resource, resMeth.Method),
	}
	meth.setup(apiDef, resMeth.Resource, rd, resMeth.VerbTitle())
	return meth
}

func serverMethodName(endpoint, displayName, verb, resName string) string {
	if len(displayName) > 0 {
		return commons.DisplayNameToFuncName(displayName)
	}
	name := commons.NormalizeIdentifier(commons.NormalizeURI(endpoint))
	return name[len(resName):] + verb
}

func (gm serverMethod) Route() string {
	return strings.Replace(gm.Endpoint, resource.CatchAllRoute, muxCatchAllRoute, 1)
}

// setup go server method, initializes all needed variables
func (gm *serverMethod) setup(apiDef *raml.APIDefinition, r *raml.Resource, rd *resource.Resource, methodName string) error {
	gm.MethodName = strings.Title(serverMethodName(gm.Endpoint, gm.DisplayName, methodName, rd.Name))

	// setting middlewares
	middlewares := []string{}

	// security middlewares
	for _, v := range gm.SecuredBy {
		if !security.ValidateScheme(v.Name, apiDef) {
			continue
		}
		// oauth2 middleware
		m, err := getOauth2MwrHandler(v)
		if err != nil {
			return err
		}
		middlewares = append(middlewares, m)
	}

	gm.Middlewares = strings.Join(middlewares, ", ")

	return nil
}

// return all libs imported by this method
func (gm serverMethod) libImported(rootImportPath string) map[string]struct{} {
	libs := map[string]struct{}{}

	// req body
	if lib := libImportPath(rootImportPath, gm.reqBody, globLibRootURLs); lib != "" {
		libs[lib] = struct{}{}
	}

	// resp body
	if lib := libImportPath(rootImportPath, gm.firstSuccessRespBodyRawType(), globLibRootURLs); lib != "" {
		libs[lib] = struct{}{}
	}
	return libs
}

// Imports return all packages needed
// by this method
func (gm serverMethod) Imports() []string {
	ip := map[string]struct{}{}
	if gm.needImportGoramlTypes() {
		ip[`"`+globRootImportPath+"/types"+`"`] = struct{}{}
	}
	for lib := range gm.libImported(globRootImportPath) {
		ip[lib] = struct{}{}
	}
	return commons.MapToSortedStrings(ip)
}

// true if req body need validation code
func (gm serverMethod) ReqBodyNeedValidation() bool {
	// we can't use t.GetBuiltinType here because
	// the reqBody type is already in Go type
	getBuiltinType := func() string {
		switch {
		case strings.HasPrefix(gm.reqBody, "[][]"): // bidimensional array
			return strings.TrimPrefix(gm.reqBody, "[][]")
		case strings.HasPrefix(gm.reqBody, "[]"): // array
			return strings.TrimPrefix(gm.reqBody, "[]")
		default:
			return gm.reqBody
		}
	}
	t := raml.Type{
		Type: getBuiltinType(),
	}

	return !t.IsBuiltin()
}
