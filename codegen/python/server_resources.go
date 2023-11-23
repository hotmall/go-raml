package python

import (
	"sort"

	"github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/raml"
)

func getServerResourcesDefs(apiDef *raml.APIDefinition) []resource.Resource {
	var rds []resource.Resource

	rs := apiDef.Resources

	// sort the keys, so we have resource sorted by keys.
	// the generated code actually don't need it to be sorted.
	// but test fixture need it
	var keys []string
	for k := range rs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// create resource def
	for _, endpoint := range keys {
		res := rs[endpoint]
		rd := resource.New(apiDef, &res, endpoint, false)
		rds = append(rds, rd)
	}
	return rds
}
