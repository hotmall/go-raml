package python

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/resource"
	"github.com/Jumpscale/go-raml/raml"
)

type pythonResource struct {
	*resource.Resource
	MiddlewaresArr []middleware
	Methods        []serverMethod
	BaseURI        string
}

func generateResources(resources []pythonResource, templates serverTemplate, dir string) error {
	allMethods := make([]serverMethod, 0)
	handlersPath := filepath.Join(dir, handlersDir)

	for _, pr := range resources {
		filename := filepath.Join(dir, strings.ToLower(pr.Name)+"_api.py")
		if err := pr.generate(filename, templates.apiFile, templates.apiName, dir); err != nil {
			return err
		}

		// Generate resource handlers
		for _, method := range pr.Methods {
			allMethods = append(allMethods, method)

			fileName := fmt.Sprintf("%vHandler.py", method.MethodName)
			if err := commons.GenerateFile(method, templates.handlerFile, templates.handlerName, filepath.Join(handlersPath, fileName), false); err != nil {
				return err
			}
		}
	}

	methods := map[string]interface{}{
		"Methods": allMethods,
	}

	// Generate handlers init file
	if err := commons.GenerateFile(methods, "./templates/python/server_handlers_init.tmpl", "server_handlers_init", filepath.Join(handlersPath, "__init__.py"), true); err != nil {
		return err
	}

	return nil
}

func (pr *pythonResource) addMiddleware(mwr middleware) {
	// check if already exist
	for _, v := range pr.MiddlewaresArr {
		if v.Name == mwr.Name {
			return
		}
	}
	pr.MiddlewaresArr = append(pr.MiddlewaresArr, mwr)
}

func newResource(rd resource.Resource, apiDef *raml.APIDefinition, kind string) pythonResource {
	baseURI := rd.APIDef.BaseURI
	if strings.Index(baseURI, "{version}") > 0 {
		baseURI = strings.Replace(baseURI, "{version}", rd.APIDef.Version, -1)
	}
	r := pythonResource{
		Resource: &rd,
		BaseURI:  baseURI,
	}
	// generate methods
	for _, rm := range rd.Methods {
		sm := newServerMethod(apiDef, &rd, rm, kind)
		r.Methods = append(r.Methods, sm)
	}

	r.setMiddlewares()
	return r
}

// set middlewares to import
func (pr *pythonResource) setMiddlewares() {
	for _, pm := range pr.Methods {
		for _, m := range pm.MiddlewaresArr {
			pr.addMiddleware(m)
		}
	}
}

// generate flask representation of an RAML resource
// It has one file : an API route and implementation
func (pr *pythonResource) generate(fileName, tmplFile, tmplName, dir string) error {
	return commons.GenerateFile(pr, tmplFile, tmplName, fileName, true)
}
