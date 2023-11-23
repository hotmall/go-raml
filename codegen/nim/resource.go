package nim

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	cr "github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/raml"
)

type resource struct {
	cr.Resource
	Methods []method
}

func newResource(name string, apiDef *raml.APIDefinition, isServer bool) resource {
	res := apiDef.Resources[name]

	rd := cr.New(apiDef, &res, name, false)

	nimRes := resource{
		Resource: rd,
	}

	for _, rm := range rd.Methods {
		var meth method
		if isServer {
			meth = newServerMethod(apiDef, rm)
		} else {
			meth = newClientMethod(rm)
		}
		nimRes.Methods = append(nimRes.Methods, meth)
	}
	return nimRes
}

// get array of all imported modules
func (r *resource) Imports() []string {
	ip := map[string]struct{}{}

	for _, m := range r.Methods {
		var names []string
		if m.ReqBody != "" {
			names = append(names, m.ReqBody)
		}
		if m.RespBody != "" {
			names = append(names, m.RespBody)
		}
		for _, name := range names {
			if tipe, ok := objectRegistered(name); ok {
				ip[tipe] = struct{}{}
			}
		}
	}
	return commons.MapToSortedStrings(ip)
}

// generate server resource API implementation
func (r *resource) generate(dir string) error {
	filename := filepath.Join(dir, r.apiName()+".nim")
	return commons.GenerateFile(r, "./templates/nim/server_resources_api_nim.tmpl", "server_resources_api_nim", filename, true)
}

// returns server's API name
func (r *resource) apiName() string {
	return strings.ToLower(r.Name) + "_api"
}

// NeedJWT returns true if this resource need JWT Library
func (r *resource) NeedJWT() bool {
	for _, m := range r.Methods {
		if m.Secured() {
			return true
		}
	}
	return false
}
func getAllResources(apiDef *raml.APIDefinition, isServer bool) []resource {
	rs := []resource{}

	// sort the keys, so we have resource sorted by keys.
	// the generated code actually don't need it to be sorted.
	// but test fixture need it
	var keys []string
	for k := range apiDef.Resources {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		rs = append(rs, newResource(k, apiDef, isServer))
	}
	return rs
}
