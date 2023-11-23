package gorestful

import (
	"sort"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
)

// generate Server's Go representation of RAML resources
func (s *Server) generateServerResources(dir string) ([]*goResource, error) {
	var serverResources []*goResource

	resources := s.apiDef.Resources

	// sort the keys, so we have resource sorted by keys.
	// the generated code actually don't need it to be sorted.
	// but test fixture need it
	var keys []string
	for k := range resources {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// create resource def
	var err error
	for _, endpoint := range keys {
		r := resources[endpoint]
		rd := resource.New(s.apiDef, &r, endpoint, true)
		pkgName := strings.ToLower(commons.NormalizeIdentifier(rd.Name))
		gr := newGoResource(&rd, pkgName)
		err = gr.generateService(dir)
		if err != nil {
			return nil, err
		}
		serverResources = append(serverResources, gr)
	}
	return serverResources, nil
}
