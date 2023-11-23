package python

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hotmall/go-raml/raml"
	"github.com/hotmall/go-raml/utils"
)

func TestSanicServer(t *testing.T) {
	Convey("sanic server generator", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("Hello world server", func() {
			apiDef := new(raml.APIDefinition)
			err = raml.ParseFile("../fixtures/raml-examples/helloworld/helloworld.raml", apiDef)
			So(err, ShouldBeNil)

			server := NewSanicServer(apiDef, "apidocs", targetDir, true, nil)
			err = server.Generate()
			So(err, ShouldBeNil)

			// check drones API implementation
			rootFixture := "./fixtures/sanic/raml-examples/helloworld"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"app.py", "app.py"},
				{"helloworld_api.py", "helloworld_api.py"},
				{"helloworld_if.py", "helloworld_if.py"},
			}

			for _, check := range checks {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, check.Result))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, check.Expected))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Convey("Congo", func() {
			apiDef := new(raml.APIDefinition)
			err = raml.ParseFile("../fixtures/congo/api.raml", apiDef)
			So(err, ShouldBeNil)

			server := NewSanicServer(apiDef, "apidocs", targetDir, true, nil)
			err = server.Generate()
			So(err, ShouldBeNil)

			// check drones API implementation
			rootFixture := "./fixtures/sanic/server/congo"
			files := []string{
				"app.py",
				"deliveries_api.py",
				"deliveries_if.py",
				"drones_api.py",
				"drones_if.py",
				"handlers/schema/User_schema.json",
				"handlers/__init__.py",
				"handlers/drones_postHandler.py",
			}

			for _, filename := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, filename))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, filename))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

			// test that this file exist
			files = []string{
				"types/User.py",
				"types/client_support.py",
				"handlers/deliveries_getHandler.py",
				"handlers/deliveries_postHandler.py",
				"handlers/deliveries_byDeliveryId_getHandler.py",
				"handlers/deliveries_byDeliveryId_patchHandler.py",
				"handlers/deliveries_byDeliveryId_deleteHandler.py",
				"handlers/drones_getHandler.py",
				"handlers/drones_postHandler.py",
				"handlers/drones_byDroneId_getHandler.py",
				"handlers/drones_byDroneId_patchHandler.py",
				"handlers/drones_byDroneId_deleteHandler.py",
				"handlers/drones_byDroneId_deliveries_getHandler.py",
			}
			for _, f := range files {
				_, err := os.Stat(filepath.Join(targetDir, f))
				So(err, ShouldBeNil)
			}

		})

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})
}
