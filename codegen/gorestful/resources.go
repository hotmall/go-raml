package gorestful

import (
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
)

// const (
// 	resourceIfTemplate  = "./templates/golang/server_resources_interface.tmpl" // resource interface template
// 	resourceAPITemplate = "./templates/golang/server_resources_api.tmpl"       // resource API template
// )

const (
	serverAPIDir = "handlers" // dir on which we put our server API implementation
	serviceDir   = "service"  // dir on which we put our service
	delegateDir  = "delegate" // dir on which we put our delegate
	resourceDir  = "resource" // dir on which we put our resource
	messageDir   = "message"  // dir on which we put our message
)

type goResource struct {
	*resource.Resource
	Methods     []serverMethod
	PackageName string
	BaseURI     string
}

func newGoResource(rd *resource.Resource, packageName string) *goResource {
	var methods []serverMethod
	for _, rm := range rd.Methods {
		methods = append(methods, newServerMethod(rm, rd.APIDef, rd))
	}

	baseURI := rd.APIDef.BaseURI
	if strings.Index(baseURI, "{version}") > 0 {
		baseURI = strings.Replace(baseURI, "{version}", rd.APIDef.Version, -1)
	}

	return &goResource{
		Resource:    rd,
		Methods:     methods,
		PackageName: packageName,
		BaseURI:     baseURI,
	}
}

// generate interface file of a resource
// func (gr *goResource) generateInterfaceFile(directory string) error {
// 	filename := directory + "/" + strings.ToLower(gr.Name) + "_if.go"
// 	return commons.GenerateFile(gr, resourceIfTemplate, "resource_if_template", filename, true)
// }

// generate API implementation in one file per method mode
// func (gr *goResource) generateAPIImplementations(dir string) error {
// 	// generate the main API impl file, which only contains struct
// 	mainCtx := map[string]interface{}{
// 		"PackageName": gr.PackageName,
// 		"Name":        gr.Name,
// 		"Endpoint":    gr.Endpoint,
// 	}

// 	mainFile := filepath.Join(dir, strings.ToLower(gr.Name)+"_api")
// 	if err := commons.GenerateFile(mainCtx, "./templates/golang/server_resource_api_main_go.tmpl",
// 		"server_resource_api_main_go", mainFile+".go", false); err != nil {
// 		return err
// 	}

// 	// generate per methods file
// 	for _, sm := range gr.Methods {
// 		ctx := map[string]interface{}{
// 			"Method":      sm,
// 			"APIName":     gr.Name,
// 			"PackageName": gr.PackageName,
// 		}
// 		filename := mainFile + "_" + sm.MethodName + ".go"
// 		if err := commons.GenerateFile(ctx, "./templates/golang/server_resource_api_impl_go.tmpl",
// 			"server_resource_api_impl_go", filename, false); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// generate Go representation of server's resource.
// A resource have two kind of files:
//   - interface file
//     contains API interface and routing code
//     always regenerated
//   - API implementation
//     implementation of the API interface.
//     Don't generate if the file already exist
// func (gr *goResource) generate(r *raml.Resource, URI, dir string, libRootURLs []string) error {
// 	if err := gr.generateInterfaceFile(dir); err != nil {
// 		return err
// 	}

// 	apiDir := filepath.Join(dir, serverAPIDir, gr.PackageName)
// 	return gr.generateAPIImplementations(apiDir)
// }

// InterfaceImportPaths returns all packages imported by
// this resource interface file
func (gr goResource) InterfaceImportPaths() []string {
	ip := map[string]struct{}{
		"net/http":               {},
		"github.com/gorilla/mux": {},
	}

	for _, gm := range gr.Methods {
		// if has middleware, we need to import middleware helper library
		if len(gm.Middlewares) > 0 {
			ip["github.com/justinas/alice"] = struct{}{}
		}
		for _, sb := range gm.SecuredBy {
			if lib := libImportPath(globRootImportPath, sb.Name, globLibRootURLs); lib != "" {
				ip[lib] = struct{}{}
			}
		}
	}
	// return sorted array for predictable order
	// we need it for unit test to always return same order
	return commons.MapToSortedStrings(ip)
}

// APIImportPaths returns all packages that need to be imported
// by the API implementation
func (gr goResource) APILibImportPaths() []string {
	ip := map[string]struct{}{}

	// methods
	for _, gm := range gr.Methods {
		for _, v := range gm.Imports() {
			ip[v] = struct{}{}
		}
	}

	// return sorted array for predictable order
	// we need it for unit test to always return same order
	return commons.MapToSortedStrings(ip)
}
