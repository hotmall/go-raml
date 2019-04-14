package gorestful

import (
	"path/filepath"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/resource"
)

// ClientService represents a root endpoint of an API
type ClientService struct {
	resource.ClientService
	PackageName string
	BaseURI     string
	Methods     []clientMethod
}

func newClientService(rootEndpoint, packageName, baseURI string, rd *resource.Resource) *ClientService {
	var methods []clientMethod

	for _, rm := range rd.Methods {
		clientMeth := newClientMethod(rm)
		methods = append(methods, clientMeth)
	}

	cs := resource.NewClientService(rootEndpoint, rd.DisplayName)
	cs.EndpointName = strings.Title(cs.EndpointName)
	return &ClientService{
		ClientService: cs,
		PackageName:   packageName,
		BaseURI:       baseURI,
		Methods:       methods,
	}
}

func (cs ClientService) filename(dir string) string {
	name := filepath.Join(dir, cs.FilenameNoExt)
	return name + ".go"
}

// NeedImportJSON returns true if this service need
// to import encoding/json
func (cs ClientService) NeedImportJSON() bool {
	for _, cm := range cs.Methods {
		if cm.needImportEncodingJSON() {
			return true
		}
	}
	return false
}

// LibImportPaths returns all imported lib
func (cs ClientService) libImportPaths() map[string]struct{} {
	ip := map[string]struct{}{}

	var needImportGoraml bool
	// methods
	for _, gm := range cs.Methods {
		for lib := range gm.libImported(globRootImportPath) {
			ip[lib] = struct{}{}
		}
		if !needImportGoraml {
			needImportGoraml = gm.needImportGoraml()
		}
	}
	if needImportGoraml {
		ip[`"`+joinImportPath(globRootImportPath, "goraml")+`"`] = struct{}{}
	}
	return ip
}

func (cs ClientService) Imports() []string {
	imports := cs.libImportPaths()
	for _, m := range cs.Methods {
		if m.needImportGoramlTypes() {
			imports[`"`+globRootImportPath+"/"+typePackage+`"`] = struct{}{}
			break
		}
	}
	return commons.MapToSortedStrings(imports)
}
