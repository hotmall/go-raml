package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hotmall/go-raml/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClientGeneration(t *testing.T) {
	Convey("test command client generattion", t, func() {
		targetdir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("Test run client command using go language", func() {
			cmd := ClientCommand{
				Language:    "go",
				Dir:         targetdir,
				RamlFile:    "../codegen/fixtures/client_resources/client.raml",
				PackageName: "theclient",
				ImportPath:  "examples.com/client",
			}
			err := cmd.Execute()
			So(err, ShouldBeNil)

			s, err := utils.TestLoadFile(filepath.Join(targetdir, "client_structapitest.go"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile("../codegen/golang/fixtures/client_resources/client_structapitest.txt")
			So(err, ShouldBeNil)

			So(tmpl, ShouldEqual, s)
		})

		Reset(func() {
			//cleanup
			os.RemoveAll(targetdir)
		})
	})
}
