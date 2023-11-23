package gorestful

import (
	"path/filepath"

	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/codegen/date"
)

// date generator
type dateGen struct {
	PackageName string
}

func (dg dateGen) generate(dir string) error {
	dates := []struct {
		Type     string
		Format   string
		FileName string
	}{
		{"date-only", "", "date_only.go"},
		{"time-only", "", "time_only.go"},
		{"datetime-only", "", "datetime_only.go"},
		{"datetime", "RFC3339", "datetime.go"},
		{"datetime", "RFC2616", "datetime_rfc2616.go"},
	}
	for _, d := range dates {
		b, err := date.Get(d.Type, d.Format)
		if err != nil {
			return err
		}
		ctx := map[string]interface{}{
			"PackageName": dg.PackageName,
			"Content":     string(b),
		}

		err = commons.GenerateFile(ctx, "./templates/golang/date.tmpl", "date", filepath.Join(dir, d.FileName), true)
		if err != nil {
			return err
		}
	}
	return nil

}
