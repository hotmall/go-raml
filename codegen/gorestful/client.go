package gorestful

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/resource"
	"github.com/Jumpscale/go-raml/raml"
)

// Client represents a Golang client
type Client struct {
	apiDef         *raml.APIDefinition
	Name           string
	BaseURI        string
	Kind           string
	libraries      map[string]*raml.Library
	PackageName    string
	RootImportPath string
	Services       map[string]*ClientService
	TargetDir      string
	libsRootURLs   []string
}

// NewClient creates a new Golang client
func NewClient(apiDef *raml.APIDefinition, kind, packageName, rootImportPath, targetDir string,
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

	baseURI := apiDef.BaseURI
	if strings.Index(baseURI, "{version}") > 0 {
		baseURI = strings.Replace(baseURI, "{version}", apiDef.Version, -1)
	}

	// creates client services objects
	services := map[string]*ClientService{}
	for endpoint, res := range apiDef.Resources {
		rd := resource.New(apiDef, &res, commons.NormalizeURITitle(endpoint), true)
		services[endpoint] = newClientService(endpoint, packageName, baseURI, &rd)
	}

	// creates client object
	client := Client{
		apiDef:         apiDef,
		Name:           commons.NormalizeIdentifier(commons.NormalizeURI(apiDef.Title)),
		BaseURI:        baseURI,
		Kind:           kind,
		libraries:      apiDef.Libraries,
		PackageName:    packageName,
		RootImportPath: rootImportPath,
		Services:       services,
		TargetDir:      targetDir,
		libsRootURLs:   libsRootURLs,
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
		command:     "client",
		kind:        gc.Kind,
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

	if gc.Kind == "grequests" {
		if err := gc.generateHelperFile(gc.TargetDir); err != nil {
			return err
		}
	}

	if err := gc.generateSecurity(gc.TargetDir); err != nil {
		return err
	}

	// if err := gc.generateServices(gc.TargetDir); err != nil {
	// 	return err
	// }
	// return gc.generateClientFile(gc.TargetDir)
	return gc.generateServices(gc.TargetDir)
}

// generate Go client helper
func (gc *Client) generateHelperFile(dir string) error {
	fileName := filepath.Join(dir, "/client_utils.go")
	return commons.GenerateFile(gc, "./templates/golang/grequests_client_utils.tmpl", "grequests_client_utils", fileName, true)
}

func (gc *Client) generateServices(dir string) error {
	switch gc.Kind {
	case "grequests":
		for _, s := range gc.Services {
			if err := commons.GenerateFile(s, "./templates/golang/grequests_client_api.tmpl", "grequests_client_api", s.filename(dir), true); err != nil {
				return err
			}
		}
	case "requests":
		for _, s := range gc.Services {
			for _, m := range s.Methods {
				c := newGoContext(m.method)
				if err := c.generateClient(dir); err != nil {
					return err
				}
			}
			if err := commons.GenerateFile(s, "./templates/golang/requests_client_api.tmpl", "requests_client_api", s.filename(dir), true); err != nil {
				return err
			}
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
	fileName := filepath.Join(dir, strings.ToLower(gc.Name)+"_api.go")
	return commons.GenerateFile(gc, "./templates/golang/grequests_client_api.tmpl", "grequests_client_api", fileName, true)
}
