package capnp

import (
	"fmt"

	"github.com/hotmall/go-raml/codegen/commons"
)

// schema for Go need additional data needed by Go compiler
// go compilation : capnp compile -I$GOPATH/src/zombiezen.com/go/capnproto2/std   -ogo *.capnp

func (s *Struct) goImports() []string {
	return []string{`using Go = import "/go.capnp"`}
}

func (s *Struct) goAnnotations() []string {
	pkg := []string{fmt.Sprintf(`$Go.package("%v")`, s.pkg)}
	annosMap := map[string]struct{}{}
	// import enums
	for _, f := range s.Fields {
		if f.Enum != nil {
			annosMap[fmt.Sprintf(`$Go.import("%v")`, f.Enum.pkg)] = struct{}{}
		}
	}
	// import our own package
	annosMap[fmt.Sprintf(`$Go.import("%v")`, s.pkg)] = struct{}{}
	annos := commons.MapToSortedStrings(annosMap)
	return append(pkg, annos...)
}
