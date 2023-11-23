package python

import (
	"path/filepath"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
)

type sanicRouteView struct {
	Name     string
	Endpoint string
	Methods  []serverMethod
}

func newSanicRouteView(endpoint string) sanicRouteView {
	// TODO : need better way to generate name
	name := strings.TrimSuffix(endpoint, "/")
	name = strings.TrimPrefix(name, "/")
	name = strings.TrimSuffix(name, ">")
	name = strings.Replace(name, "/<", "_by", -1)
	name = strings.Replace(name, ">/", "_", -1)
	name = strings.Replace(name, "/", "_", -1)
	name = strings.Replace(name, ">", "", -1)
	name = commons.NormalizeIdentifier(name)

	return sanicRouteView{
		Name:     name,
		Endpoint: endpoint,
	}
}

func (srv sanicRouteView) Route() string {
	if len(srv.Methods) == 0 {
		return srv.Endpoint
	}
	// first method is enough to generate route
	// because all methods has same route
	return srv.Methods[0].Route()
}

func (srv sanicRouteView) RouteCatchAll() string {
	return srv.Methods[0].RouteCatchAll()
}

func (srv sanicRouteView) IsCatchAllRoute() bool {
	if len(srv.Methods) == 0 {
		return false
	}
	return srv.Methods[0].IsCatchAllRoute()
}

func (s SanicServer) generateResources(dir string) error {
	if err := generateResources(s.ResourcesDef, s.Template, dir); err != nil {
		return err
	}
	if err := s.generateResourcesAPIIface(dir); err != nil {
		return err
	}
	return nil
}

func (s SanicServer) generateResourcesAPIIface(dir string) error {
	// routing & middlewares
	for _, pr := range s.ResourcesDef {
		// create sanicRouteView objects
		srvMap := map[string]sanicRouteView{}
		for _, pm := range pr.Methods {
			srv, ok := srvMap[pm.Endpoint]
			if !ok {
				srv = newSanicRouteView(pm.Endpoint)
			}

			srv.Methods = append(srv.Methods, pm)
			srvMap[pm.Endpoint] = srv
		}

		// generate
		filename := filepath.Join(dir, strings.ToLower(pr.Name)+"_if.py")
		ctx := map[string]interface{}{
			"Name":      strings.ToLower(pr.Name),
			"RouteView": srvMap,
			"PR":        pr,
		}
		if err := commons.GenerateFile(ctx, "./templates/python/server_resource_if_sanic.tmpl", "server_resource_if_sanic", filename, true); err != nil {
			return err
		}
	}

	return nil
}
