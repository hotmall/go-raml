package gorestful

import (
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/security"
	"github.com/Jumpscale/go-raml/raml"
)

type goSecurity struct {
	*security.Security
}

// generate Go representation of a security scheme
// it implemented as struct based middleware
func (gs *goSecurity) generate(dir string) error {
	fileName := path.Join(dir, "oauth2_"+gs.Name+"_middleware.go")
	return commons.GenerateFile(gs, "./templates/golang/oauth2_middleware.tmpl", "oauth2_middleware", fileName, true)
}

func generateSecurity(schemes map[string]raml.SecurityScheme, dir, packageName string) error {
	var err error

	// generate oauth2 middleware
	for k, ss := range schemes {
		if ss.Type != security.Oauth2 {
			continue
		}

		sd := security.New(ss, k, packageName)

		gss := goSecurity{Security: &sd}
		err = gss.generate(dir)

		if err != nil {
			log.Errorf("generateSecurity() failed to generate %v, err=%v", k, err)
			return err
		}
	}
	return nil
}
