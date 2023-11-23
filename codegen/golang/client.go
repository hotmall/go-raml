package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/raml"
)

// Client represents a Golang client
type Client struct {
	apiDef         *raml.APIDefinition
	Name           string
	BaseURI        string
	libraries      map[string]*raml.Library
	PackageName    string
	RootImportPath string
	Services       map[string]*ClientService
	TargetDir      string
	libsRootURLs   []string
}

// NewClient creates a new Golang client
func NewClient(apiDef *raml.APIDefinition, packageName, rootImportPath, targetDir string,
	libsRootURLs []string) (Client, error) {

	// rootImportPath only needed if we use libraries
	if rootImportPath == "" && len(apiDef.Libraries) > 0 {
		return Client{}, fmt.Errorf("--import-path can't be empty when we use libraries")
	}

	// TODO : get rid of this global variable
	rootImportPath = setRootImportPath(rootImportPath, targetDir)
	globRootImportPath = rootImportPath
	globAPIDef = apiDef
	globLibRootURLs = libsRootURLs

	// creates client services objects
	services := map[string]*ClientService{}
	for endpoint, res := range apiDef.Resources {
		rd := resource.New(apiDef, &res, commons.NormalizeURITitle(endpoint), true)
		services[endpoint] = newClientService(endpoint, packageName, &rd)
	}

	// creates client object
	client := Client{
		apiDef:         apiDef,
		Name:           commons.NormalizeIdentifier(commons.NormalizeURI(apiDef.Title)),
		BaseURI:        apiDef.BaseURI,
		libraries:      apiDef.Libraries,
		PackageName:    packageName,
		RootImportPath: rootImportPath,
		Services:       services,
		TargetDir:      targetDir,
		libsRootURLs:   libsRootURLs,
	}

	if strings.Index(client.BaseURI, "{version}") > 0 {
		client.BaseURI = strings.Replace(client.BaseURI, "{version}", apiDef.Version, -1)
	}
	return client, nil
}

// Generate generates all Go client files
func (gc Client) Generate() error {
	if err := commons.CheckDuplicatedTitleTypes(gc.apiDef); err != nil {
		return err
	}
	// helper package
	gh := goramlHelper{
		packageName: "goraml",
		packageDir:  "goraml",
	}
	if err := gh.generate(gc.TargetDir); err != nil {
		return err
	}

	// generate struct
	if err := generateAllStructs(gc.apiDef, gc.TargetDir); err != nil {
		return err
	}

	// libraries
	if err := generateLibraries(gc.libraries, gc.TargetDir, gc.libsRootURLs); err != nil {
		return err
	}

	if err := gc.generateHelperFile(gc.TargetDir); err != nil {
		return err
	}

	if err := gc.generateSecurity(gc.TargetDir); err != nil {
		return err
	}

	if err := gc.generateServices(gc.TargetDir); err != nil {
		return err
	}
	return gc.generateClientFile(gc.TargetDir)
}

// generate Go client helper
func (gc *Client) generateHelperFile(dir string) error {
	fileName := filepath.Join(dir, "/client_utils.go")
	return commons.GenerateFile(gc, "./templates/golang/client_utils_go.tmpl", "client_utils_go", fileName, true)
}

func (gc *Client) generateServices(dir string) error {
	for _, s := range gc.Services {
		//sort.Sort(resource.ByEndpoint(s.Methods))
		if err := commons.GenerateFile(s, "./templates/golang/client_service_go.tmpl", "client_service_go", s.filename(dir), true); err != nil {
			return err
		}
	}
	return nil
}

// generate security related files
// it currently only supports itsyou.online oauth2
func (gc *Client) generateSecurity(dir string) error {
	for name, ss := range gc.apiDef.SecuritySchemes {
		if v, ok := ss.Settings["accessTokenUri"]; ok {
			ctx := map[string]string{
				"ClientName":     gc.Name,
				"AccessTokenURI": fmt.Sprintf("%v", v),
				"PackageName":    gc.PackageName,
			}
			filename := filepath.Join(dir, "oauth2_client_"+name+".go")
			if err := commons.GenerateFile(ctx, "./templates/golang/oauth2_client_go.tmpl", "oauth2_client_go", filename, true); err != nil {
				return err
			}
		}
	}
	return nil
}

// generate Go client lib file
func (gc *Client) generateClientFile(dir string) error {
	fileName := filepath.Join(dir, "/client_"+strings.ToLower(gc.Name)+".go")
	return commons.GenerateFile(gc, "./templates/golang/client_go.tmpl", "client_go", fileName, true)
}
