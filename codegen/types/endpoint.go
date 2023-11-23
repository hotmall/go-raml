package types

import (
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
	"github.com/pinzolo/casee"
)

const (
	// HTTPRequest defines an HTTP request
	HTTPRequest = iota

	// HTTPResponse defines an HTTP response
	HTTPResponse
)

// Endpoint represents an endpoint defined in the RAML specs
type Endpoint struct {
	Addr     string // complete endpoint address
	Method   *raml.Method
	Resource *raml.Resource
	Verb     string
}

// ResourceName returns resource name of an endpoint
func (ep Endpoint) ResourceName() string {
	name := ep.Addr
	splt := strings.Split(name, "/")
	if len(splt) > 0 {
		name = splt[0]
	}
	return strings.TrimSuffix(strings.TrimPrefix(name, "/"), "/")
}

func getAllEndpoints(apiDef *raml.APIDefinition) map[string][]Endpoint {
	endpoints := map[string][]Endpoint{}
	for _, r := range apiDef.Resources {
		getEndpointsOfResource("", &r, endpoints)
	}
	return endpoints
}

func getEndpointsOfResource(parentPath string, r *raml.Resource, endpoints map[string][]Endpoint) {
	var methods = []struct {
		Name   string
		Method *raml.Method
	}{
		{Name: "Get", Method: r.Get},
		{"Post", r.Post},
		{"Head", r.Head},
		{"Put", r.Put},
		{"Delete", r.Delete},
		{"Patch", r.Patch},
		{"Options", r.Options},
	}

	for _, m := range methods {
		if m.Method == nil {
			continue
		}
		endp := Endpoint{
			Addr:     parentPath + r.URI,
			Method:   m.Method,
			Verb:     strings.ToUpper(m.Name),
			Resource: r,
		}
		if _, exists := endpoints[endp.Addr]; exists {
			endpoints[endp.Addr] = append(endpoints[endp.Addr], endp)
		} else {
			endpoints[endp.Addr] = []Endpoint{endp}
		}
	}

	for _, v := range r.Nested {
		getEndpointsOfResource(parentPath+r.URI, v, endpoints)
	}
}

// PascalCaseTypeName generates pascalcase type name from snackcase method name
func PascalCaseTypeName(tip TypeInBody) string {
	methodName := commons.SnackCaseServerMethodName(tip.Endpoint.Method.DisplayName, tip.Endpoint.Verb, tip.Endpoint.Resource)
	suffix := commons.ReqBodySuffix
	if tip.ReqResp == HTTPResponse {
		suffix = commons.RespBodySuffix
	}
	return casee.ToPascalCase(methodName + suffix)
}
