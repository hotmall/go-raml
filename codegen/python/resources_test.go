package python

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jumpscale/go-raml/raml"
	"github.com/Jumpscale/go-raml/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPythonResource(t *testing.T) {
	Convey("resource generator", t, func() {
		targetdir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("resource with request body", func() {
			apiDef := new(raml.APIDefinition)
			err := raml.ParseFile("../fixtures/server_resources/deliveries.raml", apiDef)
			So(err, ShouldBeNil)

			fs := NewFlaskServer(apiDef, "apidocs", targetdir, true, nil, false, "TestApp")

			err = fs.generateResources(targetdir)
			So(err, ShouldBeNil)

			// check  api implementation
			s, err := utils.TestLoadFile(filepath.Join(targetdir, "deliveries_api.py"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile("../fixtures/server_resources/deliveries_api.py")
			So(err, ShouldBeNil)
			So(s, ShouldEqual, tmpl)
		})

		Reset(func() {
			os.RemoveAll(targetdir)
		})
	})
}
