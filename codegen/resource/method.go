package resource

import (
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
)

// Method defines base Method struct
type Method struct {
	*raml.Method
	MethodName   string
	Endpoint     string
	verb         string
	Resource     *raml.Resource // resource object of this method
	Params       string         // methods params
	FuncComments []string
}

// Verb returns HTTP Verb of this method in all uppercase format
func (m Method) Verb() string {
	return m.verb
}

// VerbTitle returns HTTP Verb of this method in title format
func (m Method) VerbTitle() string {
	return strings.Title(strings.ToLower(m.verb))
}

func newMethod(r *raml.Resource, rd *Resource, m *raml.Method, methodName string) Method {
	method := Method{
		Method:       m,
		Endpoint:     r.FullURI(),
		verb:         strings.ToUpper(methodName),
		Resource:     r,
		FuncComments: commons.ParseDescription(m.Description),
	}

	return method
}

// IsCatchAllRoute returns true if this method
// use go-raml Catch-All route
func (m Method) IsCatchAllRoute() bool {
	return strings.HasSuffix(m.Endpoint, CatchAllRoute)
}

// byEndpoint implements sort interface to sort methods
// based on its endpoint address
type byEndpoint []Method

func (b byEndpoint) Len() int      { return len(b) }
func (b byEndpoint) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byEndpoint) Less(i, j int) bool {
	return b[i].Endpoint+b[i].verb < b[j].Endpoint+b[j].verb
}
