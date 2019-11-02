package gorestful

import (
	"fmt"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/resource"
	"github.com/Jumpscale/go-raml/codegen/security"
	"github.com/Jumpscale/go-raml/raml"
)

// defines go method base object
type method struct {
	resource.Method
	ResourcePath  string
	reqBody       string // type of the request body
	PackageName   string
	resps         []respBody
	URIParameters map[string]raml.NamedParameter
}

func (m method) ReqBody() string {
	return formatReqRespBody(m.reqBody)
}

func (m method) RespBody() string {
	return m.firstSuccessRespBodyType()
}

// TODO : move it to codegen/resource
type respBody struct {
	Code     int
	respType string
}

func (rb respBody) Type() string {
	return formatReqRespBody(rb.respType)
}

func formatReqRespBody(tip string) string {
	// if builtin type, no need for further processing
	if isBuiltinOrGoramlType(tip) {
		return tip
	}

	// normalize the identifier name
	normType := commons.NormalizeIdentifierWithLib(tip, globAPIDef)

	// convert type name (last elem) to Title case
	elems := strings.Split(normType, ".")
	lastIdx := len(elems) - 1
	elems[lastIdx] = strings.Title(elems[lastIdx])

	return addPackage(typePackage, strings.Join(elems, "."))
}

// add package name to the request and response body
func addPackage(pkgName, typeStr string) string {
	if typeStr == "" {
		return ""
	}

	// we can't use raml package to get basic type and decide whether it is array
	// or not because in this phase the type is already Go type, not RAML type
	switch {
	case strings.HasPrefix(typeStr, "[][]"): // bidimensi array
		return "[][]" + addPackage(pkgName, strings.TrimPrefix(typeStr, "[][]"))
	case strings.HasPrefix(typeStr, "[]"): // array array
		return "[]" + addPackage(pkgName, strings.TrimPrefix(typeStr, "[]"))
	default:
		return aliasPackage(fmt.Sprintf("%v.%v", pkgName, typeStr))
	}
}

func aliasPackage(typeStr string) string {
	if typeStr == "" || strings.Index(typeStr, ".") < 0 {
		return typeStr
	}
	pkgs := strings.Split(typeStr, ".")
	if len(pkgs) != 3 {
		return typeStr
	}
	return fmt.Sprintf("%v_%v.%v", pkgs[1], pkgs[0], pkgs[2])
}

func newMethod(resMeth resource.Method) *method {
	var resps []respBody

	methodName := commons.DisplayNameToFuncName(resMeth.DisplayName)

	// creates response body
	for code, resp := range resMeth.Responses {
		resp := respBody{
			Code:     commons.AtoiOrPanic(string(code)),
			respType: setBodyName(resp.Bodies, methodName, commons.RespBodySuffix),
		}
		if resp.respType != "" {
			resps = append(resps, resp)
		}
	}

	// copy uri parameters
	URIParameters := make(map[string]raml.NamedParameter)
	copyURIParameters(URIParameters, resMeth.Resource)

	return &method{
		Method:        resMeth,
		ResourcePath:  commons.ParamizingURI(resMeth.Endpoint, "+"),
		reqBody:       setBodyName(resMeth.Bodies, methodName, commons.ReqBodySuffix),
		resps:         resps,
		URIParameters: URIParameters,
	}
}

func (m method) HasRespBody() bool {
	return len(m.RespBodyTypes()) > 0
}

// RespBodyTypes returns all possible type of response body
func (m method) RespBodyTypes() []respBody {
	return m.resps
}

// FailedRespBodyTypes return all response body that considered a failed response
// i.e. non 2xx status code
func (m method) FailedRespBodyTypes() (resps []respBody) {
	for _, resp := range m.RespBodyTypes() {
		if resp.Code < 200 || resp.Code >= 300 {
			resps = append(resps, resp)
		}
	}
	return
}

// SuccessRespBodyTypes returns all response body that considered as success
// i.e. 2xx status code
func (m method) SuccessRespBodyTypes() (resps []respBody) {
	for _, resp := range m.RespBodyTypes() {
		if resp.Code >= 200 && resp.Code < 300 {
			resps = append(resps, resp)
		}
	}
	return
}

func (m method) firstSuccessRespBodyType() string {
	resps := m.SuccessRespBodyTypes()
	if len(resps) == 0 {
		return ""
	}
	return resps[0].Type()
}

func (m method) FirstSuccessRespStatus() int {
	resps := m.SuccessRespBodyTypes()
	if len(resps) == 0 {
		return 204
	}
	return resps[0].Code
}

func (m method) firstSuccessRespBodyRawType() string {
	resps := m.SuccessRespBodyTypes()
	if len(resps) == 0 {
		return ""
	}
	return resps[0].respType
}

// returns true if need to import goraml generated types
func (m method) needImportGoramlTypes() bool {
	pkgPrefix := typePackage + "."
	needImport := func(typeStr string) bool {
		switch {
		case strings.HasPrefix(typeStr, "[][]"+pkgPrefix):
			return true
		case strings.HasPrefix(typeStr, "[]"+pkgPrefix):
			return true
		case strings.HasPrefix(typeStr, pkgPrefix):
			return true
		default:
			return false
		}
	}
	if needImport(m.ReqBody()) {
		return true
	}
	for _, resp := range m.RespBodyTypes() {
		if needImport(resp.Type()) {
			return true
		}
	}
	return false
}

// get oauth2 middleware handler from a security scheme
func getOauth2MwrHandler(ss raml.DefinitionChoice) (string, error) {
	// construct security scopes
	quotedScopes, err := security.GetQuotedScopes(ss)
	if err != nil {
		return "", err
	}
	scopesArgs := strings.Join(quotedScopes, ", ")

	// middleware name
	// need to handle case where it reside in different package
	var packageName string
	name := ss.Name

	if splitted := strings.Split(name, "."); len(splitted) == 2 {
		packageName = splitted[0]
		name = splitted[1]
	}
	mwr := fmt.Sprintf(`NewOauth2%vMiddleware([]string{%v}).Handler`, name, scopesArgs)
	if packageName != "" {
		mwr = packageName + "." + mwr
	}
	return mwr, nil
}

// setBodyName set name of method's request/response body.
//
// Rules:
//	- use bodies.Type if not empty and not `object`
//	- use bodies.ApplicationJSON.Type if not empty and not `object`
//	- use prefix+suffix if:
//		- not meet previous rules
//		- previous rules produces JSON string
func setBodyName(bodies raml.Bodies, prefix, suffix string) string {
	var tipe string
	prefix = commons.NormalizeURITitle(prefix)

	if len(bodies.Type) > 0 && bodies.Type != "object" {
		tipe = convertToGoType(bodies.Type, "")
	} else if bodies.ApplicationJSON != nil {
		if bodies.ApplicationJSON.TypeString() != "" && bodies.ApplicationJSON.TypeString() != "object" {
			tipe = convertToGoType(bodies.ApplicationJSON.TypeString(), "")
		} else {
			tipe = prefix + suffix
		}
	}

	if commons.IsJSONString(tipe) {
		tipe = prefix + suffix
	}

	return tipe
}
