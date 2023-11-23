package python

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
)

// defines a python client lib method
type clientMethod struct {
	*method
	ResourcePath string
	PRArgs       string // python requests's args
	PRCall       string // the way we call python request
}

func newClientMethod(rm resource.Method) clientMethod {
	meth := newMethod(rm)
	cm := clientMethod{
		method:       meth,
		ResourcePath: commons.ParamizingURI(rm.Endpoint, "+"),
	}
	cm.setup()
	return cm
}

func (pcm *clientMethod) setup() {
	prArgs := []string{"uri"}  // prArgs are arguments we supply to python request
	params := []string{"self"} // params are method signature params

	// for method with request body, we add `data` argument
	if !pcm.Bodies.IsEmpty() || pcm.Verb() == "PUT" || pcm.Verb() == "POST" || pcm.Verb() == "PATCH" {
		params = append(params, "data")
		prArgs = append(prArgs, "data")
	} else {
		prArgs = append(prArgs, "None")
	}

	// construct prArgs string from the array
	prArgs = append(prArgs, "headers", "query_params", "content_type")
	pcm.PRArgs = strings.Join(prArgs, ", ")

	// construct method signature
	params = append(params, resource.GetResourceParams(pcm.Resource)...)
	params = append(params, pcm.queryParams()...)
	params = append(params, "headers=None", "query_params=None", `content_type="application/json"`)
	pcm.Params = strings.Join(params, ", ")

	// python request call
	// we encapsulate the call to put, post, and patch.
	// To be able to accept plain string or dict.
	// if it is a dict, we encode it to json
	if pcm.Verb() == "GET" || pcm.Verb() == "PUT" || pcm.Verb() == "POST" || pcm.Verb() == "PATCH" || pcm.Verb() == "DELETE" {
		pcm.PRCall = fmt.Sprintf("self.client.%v", strings.ToLower(pcm.Verb()))
	} else {
		pcm.PRCall = fmt.Sprintf("self.client.session.%v", strings.ToLower(pcm.Verb()))
	}

	if len(pcm.DisplayName) > 0 {
		pcm.MethodName = commons.DisplayNameToFuncName(pcm.DisplayName)
	} else {
		pcm.MethodName = commons.NormalizeIdentifier(snakeCaseResourceURI(pcm.Resource) + "_" + strings.ToLower(pcm.Verb()))
	}
}

func (pcm *clientMethod) queryParams() (params []string) {
	var (
		paramsWithDefault []string
		paramsNoDefault   []string
	)

	for name, qp := range pcm.QueryParameters {
		if qp.Default == nil {
			paramsNoDefault = append(paramsNoDefault, name)
		} else {
			paramsWithDefault = append(paramsWithDefault, fmt.Sprintf("%s=%v", name, qp.Default))
		}
	}

	if len(paramsNoDefault) > 0 {
		sort.Strings(paramsNoDefault)
		params = append(params, paramsNoDefault...)
	}

	if len(paramsWithDefault) > 0 {
		sort.Strings(paramsWithDefault)
		params = append(params, paramsWithDefault...)
	}

	return
}

func (pcm clientMethod) imports() []string {
	var imports []string

	// import from response body
	for _, rb := range pcm.resps {
		basicType := rb.BasicType()
		imports = append(imports, "from ."+basicType+" import "+basicType)
	}
	return imports
}

// Returns true if this client method has response body
func (pcm clientMethod) HasRespBody() bool {
	return len(pcm.resps) > 0
}

func (pcm clientMethod) Route() string {
	if pcm.ResourcePath == "" {
		return ""
	}

	route := "+ " + pcm.ResourcePath
	if !pcm.IsCatchAllRoute() {
		return route
	}

	return strings.Replace(route, `/"`, `"`, 1)

}
