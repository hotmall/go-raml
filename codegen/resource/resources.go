package resource

import (
	"regexp"
	"sort"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
)

var (
	//reResource      = regexp.MustCompile(`({{1}[\w\s]+}{1})`)
	reResourceParam = regexp.MustCompile(`\{(.*?)\}`)
)

// Resource is Go code representation of a resource
type Resource struct {
	APIDef   *raml.APIDefinition
	Name     string   // resource name
	Endpoint string   // root endpoint
	Methods  []Method // all methods of this resource
	*raml.Resource
}

// New creates a resource definition
func New(apiDef *raml.APIDefinition, r *raml.Resource, endpoint string, sortMethod bool) Resource {
	name := func() string {
		if r.DisplayName == "" {
			return commons.NormalizeIdentifier(strings.Title(commons.NormalizeURI(endpoint)))
		}
		return strings.Title(commons.DisplayNameToFuncName(r.DisplayName))
	}()

	res := Resource{
		Endpoint: endpoint,
		APIDef:   apiDef,
		Name:     name,
		Resource: r,
	}

	res.generateMethods(r)
	if sortMethod {
		sort.Sort(byEndpoint(res.Methods))
	}
	return res
}

// add a method to resource definition
func (rd *Resource) addMethod(r *raml.Resource, m *raml.Method, methodName string) {
	if m == nil {
		return
	}
	rd.Methods = append(rd.Methods, newMethod(r, rd, m, methodName))
}

// GenerateMethods generates all methods of a resource recursively
func (rd *Resource) generateMethods(r *raml.Resource) {
	rd.addMethod(r, r.Get, "Get")
	rd.addMethod(r, r.Post, "Post")
	rd.addMethod(r, r.Put, "Put")
	rd.addMethod(r, r.Patch, "Patch")
	rd.addMethod(r, r.Delete, "Delete")
	rd.addMethod(r, r.Options, "Options")
	rd.addMethod(r, r.Head, "Head")

	for _, v := range r.Nested {
		rd.generateMethods(v)
	}
}

// _getResourceParams is the recursive function of getResourceParams
func _getResourceParams(r *raml.Resource, params []string) []string {
	if r == nil {
		return params
	}

	matches := reResourceParam.FindAllString(r.URI, -1)
	for _, match := range matches {
		matchEscaped := commons.NormalizeIdentifier(match[1 : len(match)-1])
		params = append(params, matchEscaped)
	}

	return _getResourceParams(r.Parent, params)
}

// GetResourceParams get all params of a resource
// examples:
// /users  							  : no params
// /users/{userId}					  : params 1 = userId
// /users/{userId}/address/{addressId : params 1= userId, param 2= addressId
func GetResourceParams(r *raml.Resource) []string {
	params := []string{}
	return _getResourceParams(r, params)
}
