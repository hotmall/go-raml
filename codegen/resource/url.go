package resource

import (
	"github.com/hotmall/go-raml/raml"
)

const (
	// go-raml extension to express Catch-All URL params
	// for example /path/{path:*} is going to match:
	//  - /path/a    -> path = a
	//  - /path/a/b  -> path = a/b
	//  - /path/a/b/c -> path = a/b/c
	CatchAllRoute = "{path:*}"
)

// HasCatchAllInRootRoute returns true if the given
// raml api definition has catch all route as root endpoint
func HasCatchAllInRootRoute(apiDef *raml.APIDefinition) bool {
	rootCatchAllRoute := "/" + CatchAllRoute
	for endpoint := range apiDef.Resources {
		if endpoint == rootCatchAllRoute {
			return true
		}
	}
	return false
}
