package python

import (
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/codegen/security"
	"github.com/hotmall/go-raml/raml"
	"github.com/pinzolo/casee"
)

// python server method
type serverMethod struct {
	*method
	MiddlewaresArr []middleware
	SecuredBy      []raml.DefinitionChoice
}

// setup sets all needed variables
func (sm *serverMethod) setup(apiDef *raml.APIDefinition, r *raml.Resource, rd *resource.Resource, resourceParams []string) error {
	sm.MethodName = commons.SnackCaseServerMethodName(sm.DisplayName, sm.Verb(), r)

	if commons.HasJSONBody(&(sm.Bodies)) {
		// TODO : make it to call proper func
		sm.reqBody = casee.ToPascalCase(sm.MethodName + commons.ReqBodySuffix)
	}

	sm.Params = strings.Join(resourceParams, ", ")
	sm.Endpoint = strings.Replace(sm.Endpoint, "{", "<", -1)
	sm.Endpoint = strings.Replace(sm.Endpoint, "}", ">", -1)

	// security middlewares
	for _, v := range sm.SecuredBy {
		if !security.ValidateScheme(v.Name, apiDef) {
			continue
		}
		// oauth2 middleware
		m, err := newPythonOauth2Middleware(v)
		if err != nil {
			log.Errorf("error creating middleware for method.err = %v", err)
			return err
		}
		sm.MiddlewaresArr = append(sm.MiddlewaresArr, m)
	}
	return nil
}

func newServerMethod(apiDef *raml.APIDefinition, rd *resource.Resource, rm resource.Method, kind string) serverMethod {
	meth := newMethod(rm)
	sm := serverMethod{
		method:    meth,
		SecuredBy: security.GetMethodSecuredBy(apiDef, rm.Resource, rm.Method),
	}
	sm.reqBody = setServerReqBodyName(meth.reqBody)

	params := resource.GetResourceParams(rm.Resource)
	if kind == serverKindSanic {
		params = append([]string{"request"}, params...)
	}
	sm.setup(apiDef, rm.Resource, rd, params)
	return sm

}

func (sm *serverMethod) RespBody() string {
	return sm.firstSuccessRespBodyType()
}

// HasReqValidator returns true if this method
// has request body validator:
// if request body is not builtin type
func (sm serverMethod) HasReqValidator() bool {
	jsObj, ok := jsObjects[sm.reqBody]
	if !ok {
		return false
	}

	return !commons.IsBuiltinType(jsObj.T.TypeString())
}

func (sm serverMethod) Route() string {
	if !sm.IsCatchAllRoute() {
		return sm.escapedEndpoint()
	}
	return strings.Replace(sm.Endpoint, catchAllRouteSuffix, catchAllRouteSuffixFlask, 1)
}

func (sm serverMethod) RouteCatchAll() string {
	return strings.TrimSuffix(sm.Endpoint, catchAllRouteSuffix)
}

func (sm serverMethod) IsCatchAllRoute() bool {
	return strings.HasSuffix(sm.Endpoint, catchAllRouteSuffix)
}

func (sm serverMethod) ParamsSanicWithCatchAll() string {
	if !sm.IsCatchAllRoute() {
		return sm.Params
	}
	return strings.Replace(sm.Params, "path", `path=""`, 1)
}

func setServerReqBodyName(bodyName string) string {
	if !commons.IsArrayType(bodyName) {
		return bodyName
	}
	return jsArrayName(bodyName)
}
