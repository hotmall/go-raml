package python

import (
	"path/filepath"
	"strings"
	//"github.com/hotmall/go-raml/codegen/resource"
)

func (fs FlaskServer) generateResources(dir string) error {
	for _, pr := range fs.ResourcesDef {

		//pr := newResourceFromDef(rdi.(resource.Resource), fs.APIDef, newServerMethodFlask)

		filename := filepath.Join(dir, strings.ToLower(pr.Name)+"_api.py")

		if err := pr.generate(filename, "./templates/python/server_resource_api_flask.tmpl", "server_resource_api_flask", dir); err != nil {
			return err
		}
	}
	return nil
}
