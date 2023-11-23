package nim

import (
	"path/filepath"

	"github.com/hotmall/go-raml/codegen/commons"
)

type clientService struct {
	resource
}

func newClientService(r resource) clientService {
	return clientService{
		resource: r,
	}
}

func (cs *clientService) generate(dir string) error {
	filename := filepath.Join(dir, cs.Name+"_service.nim")
	return commons.GenerateFile(cs, "./templates/nim/client_service_nim.tmpl", "client_service_nim", filename, true)
}

func (cs *clientService) Imports() []string {
	return cs.resource.Imports()
}

func (cs *clientService) ClientName() string {
	return clientName(cs.APIDef)
}
