package tarantool

import (
	"strings"

	"sort"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/raml"
	"github.com/pinzolo/casee"
)

type response struct {
	Code     int
	RespType string
}

// Method is a tarantool code representation of a method
type method struct {
	resource.Method
	ReqBody   string
	Responses []response
}

// Resource is tarantool code representation of a resource
type tarantoolResource struct {
	resource.Resource
	Methods []*method
}

// Handler returns the name of the function that handles requests for this method
func (m *method) Handler() string {
	var name string
	if len(m.DisplayName) > 0 {
		name = commons.DisplayNameToFuncName(m.DisplayName)
	} else {
		name = commons.NormalizeIdentifier(commons.NormalizeURI(m.Endpoint)) + m.Verb()
	}
	return casee.ToSnakeCase(name + "Handler")
}

func formatUri(uri string) string {
	return strings.Replace(strings.Replace(uri, "{", ":", -1), "}", "", -1)
}

// URI returns the tarantool URI of the method Endpoint
func (m *method) URI() string {
	return formatUri(m.Endpoint)

}

func getServerResourcesDefs(apiDef *raml.APIDefinition) []tarantoolResource {
	var trs []tarantoolResource

	// create tarantoolResource
	for endpoint, res := range apiDef.Resources {
		rs := resource.New(apiDef, &res, endpoint, false)
		tr := tarantoolResource{Resource: rs}

		// generate methods
		for _, rm := range rs.Methods {
			var resps []response
			// creates response body
			for code, resp := range rm.Responses {
				resp := response{
					Code:     commons.AtoiOrPanic(string(code)),
					RespType: setBodyName(resp.Bodies, rm.Endpoint+rm.VerbTitle(), commons.RespBodySuffix),
				}
				if resp.RespType != "" {
					resps = append(resps, resp)
				}
			}
			sort.SliceStable(resps, func(i, j int) bool {
				return resps[i].Code < resps[j].Code
			})

			normalizedEndpoint := commons.NormalizeURITitle(rm.Endpoint)
			method := &method{
				Method:    rm,
				ReqBody:   setBodyName(rm.Bodies, normalizedEndpoint+rm.VerbTitle(), commons.ReqBodySuffix),
				Responses: resps,
			}
			tr.Methods = append(tr.Methods, method)
		}

		trs = append(trs, tr)
	}

	sort.SliceStable(trs, func(i, j int) bool {
		return trs[i].Name < trs[j].Name
	})
	return trs
}

// setBodyName set name of method's request/response body.
//
// Rules:
//  - use bodies.Type if not empty and not `object`
//  - use bodies.ApplicationJSON.Type if not empty and not `object`
//  - use prefix+suffix if:
//      - not meet previous rules
//      - previous rules produces JSON string
func setBodyName(bodies raml.Bodies, prefix, suffix string) string {
	var bodyName string
	prefix = commons.NormalizeURITitle(prefix)

	if len(bodies.Type) > 0 && bodies.Type != "object" {
		bodyName = bodies.Type
	} else if bodies.ApplicationJSON != nil {
		if bodies.ApplicationJSON.TypeString() != "" && bodies.ApplicationJSON.TypeString() != "object" {
			bodyName = bodies.ApplicationJSON.TypeString()
		} else {
			bodyName = prefix + suffix
		}
	}

	if commons.IsJSONString(bodyName) {
		bodyName = prefix + suffix
	}

	return casee.ToPascalCase(bodyName)

}
