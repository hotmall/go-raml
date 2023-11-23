package nim

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
)

// Client represents a Nim client
type Client struct {
	APIDef *raml.APIDefinition
	Dir    string
}

// NewClient creates a new Nim client
func NewClient(apiDef *raml.APIDefinition, dir string) Client {
	return Client{
		APIDef: apiDef,
		Dir:    dir,
	}
}

// Generate generates all Nim client files
func (c *Client) Generate() error {
	rs := getAllResources(c.APIDef, false)

	if err := generateAllObjects(c.APIDef, c.Dir); err != nil {
		return err
	}
	// services files
	if err := c.generateServices(rs); err != nil {
		return err
	}

	if err := c.generateSecurity(); err != nil {
		return err
	}
	// main client file
	if err := c.generateMain(); err != nil {
		return err
	}
	return nil
}

func (c *Client) generateMain() error {
	filename := filepath.Join(c.Dir, clientName(c.APIDef)+".nim")
	return commons.GenerateFile(c, "./templates/nim/client_nim.tmpl", "client_nim", filename, true)
}

func (c *Client) generateServices(rs []resource) error {
	for _, r := range rs {
		cs := newClientService(r)
		if err := cs.generate(c.Dir); err != nil {
			return err
		}
	}
	return nil
}

// generate security related files
// it currently only supports itsyou.online oauth2
func (c *Client) generateSecurity() error {
	for name, ss := range c.APIDef.SecuritySchemes {
		if v, ok := ss.Settings["accessTokenUri"]; ok {
			ctx := map[string]string{
				"ClientName": clientName(c.APIDef),
				"BaseURI":    fmt.Sprintf("%v", v),
			}
			filename := filepath.Join(c.Dir, "oauth2_client_"+name+".nim")
			if err := commons.GenerateFile(ctx, "./templates/nim/oauth2_client_nim.tmpl", "oauth2_client_nim", filename, true); err != nil {
				return err
			}
		}
	}
	return nil
}

// returns client name of an API definition
func clientName(apiDef *raml.APIDefinition) string {
	splt := strings.Split(apiDef.Title, " ")
	return "client_" + strings.ToLower(splt[0])
}
