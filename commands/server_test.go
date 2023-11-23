package commands

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hotmall/go-raml/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServerGeneration(t *testing.T) {
	Convey("test command server generation", t, func() {
		targetdir, err := ioutil.TempDir("", "test_server_command")
		So(err, ShouldBeNil)
		Convey("Test run server command using go language", func() {

			cmd := ServerCommand{
				Language:    "go",
				Dir:         targetdir,
				RamlFile:    "../codegen/fixtures/server/user_api/api.raml",
				ImportPath:  "examples.com/ramlcode",
				PackageName: "main",
			}
			err := cmd.Execute()
			So(err, ShouldBeNil)

			// check users api implementation
			s, err := utils.TestLoadFile(filepath.Join(targetdir, "handlers", "users", "users_api.go"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile("../codegen/fixtures/server/user_api/users_api.txt")
			So(err, ShouldBeNil)
			So(s, ShouldEqual, tmpl)

			// check user interface
			s, err = utils.TestLoadFile(filepath.Join(targetdir, "users_if.go"))
			So(err, ShouldBeNil)

			tmpl, err = utils.TestLoadFile("../codegen/fixtures/server/user_api/users_if.txt")
			So(err, ShouldBeNil)
			So(s, ShouldEqual, tmpl)

			// check main file
			s, err = utils.TestLoadFile(filepath.Join(targetdir, "main.go"))
			So(err, ShouldBeNil)

			tmpl, err = utils.TestLoadFile("../codegen/fixtures/server/user_api/main.txt")
			So(err, ShouldBeNil)
			So(s, ShouldEqual, tmpl)
		})

		Reset(func() {
			//cleanup
			os.RemoveAll(targetdir)
		})
	})
}

func TestServerNoMainGeneration(t *testing.T) {
	Convey("test command server generation without a main", t, func() {
		targetdir, err := ioutil.TempDir("", "test_server_command")
		So(err, ShouldBeNil)
		Convey("Test run server command without a main", func() {

			cmd := ServerCommand{
				Dir:              targetdir,
				RamlFile:         "../codegen/fixtures/server/user_api/api.raml",
				PackageName:      "main",
				Language:         "go",
				ImportPath:       "examples.com/ramlcode",
				NoMainGeneration: true,
			}
			err := cmd.Execute()
			So(err, ShouldBeNil)

			// check main fil
			if _, err := os.Stat(filepath.Join(targetdir, "main.go")); err == nil {
				So(errors.New("main.go file exists"), ShouldBeNil)
			}
		})

		Reset(func() {
			//cleanup
			os.RemoveAll(targetdir)
		})
	})
}
