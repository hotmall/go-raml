package codegen

import (
	"github.com/hotmall/go-raml/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {
	Convey("server generator", t, func() {
		targetdir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)
		Convey("simple Go server", func() {
			s := Server{
				RAMLFile:       "./fixtures/server/user_api/api.raml",
				Kind:           "",
				Dir:            targetdir,
				PackageName:    "main",
				Lang:           "go",
				APIDocsDir:     "apidocs",
				RootImportPath: "examples.com/ramlcode",
				WithMain:       true,
			}
			err := s.Generate()
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/server/user_api/"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"main.go", "main.txt"},
				{"routes.go", "routes.txt"},
				{"users_if.go", "users_if.txt"},
				{filepath.Join("handlers", "users", "users_api.go"), "users_api.txt"},
				{"helloworld_if.go", "helloworld_if.txt"},
				{filepath.Join("handlers", "helloworld", "helloworld_api.go"), "helloworld_api.txt"},
				// goraml package
				{"goraml/datetime.go", "goraml/datetime.txt"},
				{"goraml/date_only.go", "goraml/date_only.txt"},
				{"goraml/datetime_only.go", "goraml/datetime_only.txt"},
				{"goraml/datetime_rfc2616.go", "goraml/datetime_rfc2616.txt"},
				{"goraml/time_only.go", "goraml/time_only.txt"},
				{"goraml/struct_input_validator.go", "goraml/struct_input_validator.txt"},
			}
			for _, check := range checks {
				s, err := utils.TestLoadFile(filepath.Join(targetdir, check.Result))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, check.Expected))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Reset(func() {
			os.RemoveAll(targetdir)
		})
	})
}
