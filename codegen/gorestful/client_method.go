package gorestful

import (
	"fmt"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/resource"
	"github.com/Jumpscale/go-raml/raml"
)

type clientMethod struct {
	*method
}

func newClientMethod(resMeth resource.Method) clientMethod {
	goMeth := newMethod(resMeth)
	cm := clientMethod{
		method: goMeth,
	}
	cm.setup(resMeth.VerbTitle())
	return cm
}

func (gcm *clientMethod) setup(methodName string) {
	// build func/method params
	buildParams := func(r *raml.Resource, bodyType string) string {
		params := resource.GetResourceParams(r)

		if len(params) > 0 {
			// all params has string type
			params[len(params)-1] = params[len(params)-1] + " string"
		}

		// append request body type
		if len(bodyType) > 0 {
			params = append(params, "body "+bodyType)
		}

		// append header
		params = append(params, "headers,queryParams map[string]string")

		return strings.Join(params, ", ")
	}

	// method name
	name := commons.NormalizeURITitle(gcm.Endpoint)

	if len(gcm.DisplayName) > 0 {
		gcm.MethodName = commons.DisplayNameToFuncName(gcm.DisplayName)
	} else {
		gcm.MethodName = name + methodName
	}
	gcm.MethodName = commons.NormalizeIdentifier(strings.Title(gcm.MethodName))

	// method param
	gcm.Params = buildParams(gcm.Resource, gcm.ReqBody())
}

// return true if this method need to import encoding/json
func (gcm clientMethod) needImportEncodingJSON() bool {
	return len(gcm.SuccessRespBodyTypes()) > 0
}

func (gcm clientMethod) libImported(rootImportPath string) map[string]struct{} {
	libs := map[string]struct{}{}

	// req body
	if lib := libImportPath(rootImportPath, gcm.reqBody, globLibRootURLs); lib != "" {
		libs[lib] = struct{}{}
	}
	// resp body
	for _, resp := range gcm.RespBodyTypes() {
		if lib := libImportPath(rootImportPath, resp.respType, globLibRootURLs); lib != "" {
			libs[lib] = struct{}{}
		}

	}
	return libs
}

// ReturnTypes returns all types returned by this method
func (gcm clientMethod) ReturnTypes() string {
	var types []string
	for _, resp := range gcm.SuccessRespBodyTypes() {
		types = append(types, resp.Type())
	}
	types = append(types, []string{"*http.Response", "error"}...)

	return fmt.Sprintf("(%v)", strings.Join(types, ","))
}

func (gcm clientMethod) needImportGoraml() bool {
	return gcm.HasRespBody()
}

func (gcm clientMethod) Route() string {
	if gcm.ResourcePath == "" {
		return ""
	}

	route := "+" + gcm.ResourcePath

	if !gcm.IsCatchAllRoute() {
		return route
	}

	return strings.Replace(route, `/"`, `"`, 1)
}
