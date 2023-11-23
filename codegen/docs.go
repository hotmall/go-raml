package codegen

import (
	"fmt"
	"github.com/hotmall/go-raml/raml"
)

const (
	formatMarkdown = "markdown"
)

// GenerateDocs generate markdown docs from RAML specs
func GenerateDocs(apiDef *raml.APIDefinition, format, output string) error {
	switch format {
	case formatMarkdown:
		md := &markdownDocs{
			api:    apiDef,
			output: output,
		}
		return md.generate()
	default:
		return fmt.Errorf("unknown format '%s'", format)
	}
}
