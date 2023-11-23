package python

import (
	"regexp"
	"sort"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/raml"
)

var (
	reEndpoint = regexp.MustCompile(`[^/<A-Za-z0-9>]+`)
)

const (
	catchAllRouteSuffix      = "<path:*>"
	catchAllRouteSuffixFlask = "<path:path>"
)

type method struct {
	resource.Method
	reqBody string
	resps   []respBody
}

func (m method) ReqBody() string {
	return commons.NormalizeIdentifierWithLib(m.reqBody, globAPIDef)
}

func (m method) escapedEndpoint() string {
	return reEndpoint.ReplaceAllString(m.Endpoint, "_")
}

type respBody struct {
	Code     int
	respType string
}

// byStatusCode implements sorter interface which sort resp body by status code
type byStatusCode []respBody

func (b byStatusCode) Len() int      { return len(b) }
func (b byStatusCode) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byStatusCode) Less(i, j int) bool {
	return b[i].Code < b[j].Code
}

func (rb *respBody) IsArray() bool {
	t := raml.Type{Type: rb.respType}
	return t.IsArray() || t.IsBidimensiArray()
}

func (rb *respBody) BasicType() string {
	return commons.NormalizeIdentifierWithLib(commons.GetBasicType(rb.respType), globAPIDef)
}

func newMethod(rm resource.Method) *method {
	var resps []respBody
	// creates response body
	for code, resp := range rm.Responses {
		resp := respBody{
			Code:     commons.AtoiOrPanic(string(code)),
			respType: setBodyName(resp.Bodies, rm.Endpoint+rm.VerbTitle(), commons.RespBodySuffix),
		}
		if resp.respType != "" {
			resps = append(resps, resp)
		}
	}
	sort.Sort(byStatusCode(resps))

	normalizedEndpoint := commons.NormalizeURITitle(rm.Endpoint)
	return &method{
		Method:  rm,
		reqBody: setBodyName(rm.Bodies, normalizedEndpoint+rm.VerbTitle(), commons.ReqBodySuffix),
		resps:   resps,
	}
}

// SuccessRespBodyTypes returns all possible type of response body
func (m method) SuccessRespBodyTypes() (resps []respBody) {
	for _, resp := range m.resps {
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
	return resps[0].respType
}

// create snake case function name from a resource URI
func snakeCaseResourceURI(r *raml.Resource) string {
	return _snakeCaseResourceURI(r, "")
}

func _snakeCaseResourceURI(r *raml.Resource, completeURI string) string {
	if r == nil {
		return completeURI
	}
	var snake string
	if len(r.URI) > 0 {
		uri := commons.NormalizeURI(r.URI)
		if len(uri) > 0 {
			if r.Parent != nil { // not root resource, need to add "_"
				snake = "_"
			}

			if strings.HasPrefix(r.URI, "/{") {
				snake += "by" + strings.ToUpper(uri[:1])
			} else {
				snake += strings.ToLower(uri[:1])
			}

			if len(uri) > 1 { // append with the rest of uri
				snake += uri[1:]
			}
		}
	}
	return _snakeCaseResourceURI(r.Parent, snake+completeURI)
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
	var tipe string
	prefix = commons.NormalizeURITitle(prefix)

	if len(bodies.Type) > 0 && bodies.Type != "object" {
		tipe = bodies.Type
	} else if bodies.ApplicationJSON != nil {
		if bodies.ApplicationJSON.TypeString() != "" && bodies.ApplicationJSON.TypeString() != "object" {
			tipe = bodies.ApplicationJSON.TypeString()
		} else {
			tipe = prefix + suffix
		}
	}

	if commons.IsJSONString(tipe) {
		tipe = prefix + suffix
	}

	return tipe

}
