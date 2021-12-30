package python

import (
	"path/filepath"

	"github.com/Jumpscale/go-raml/codegen/commons"
)

type classInit struct {
	Types map[string][]string
}

func (ci *classInit) append(fileName string, typeNames []string) {
	ci.Types[fileName] = make([]string, 0)
	ci.Types[fileName] = append(ci.Types[fileName], typeNames...)
}

func (ci *classInit) generate(dir string) error {
	template := "./templates/python/class_init.tmpl"
	templateName := "class_init"
	fileName := filepath.Join(dir, "__init__.py")
	return commons.GenerateFile(ci, template, templateName, fileName, true)
}
