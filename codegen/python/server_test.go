package python

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/Jumpscale/go-raml/raml"
	"github.com/Jumpscale/go-raml/utils"
)

func TestServer(t *testing.T) {
	Convey("server generator", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/congo/api.raml", apiDef)
		So(err, ShouldBeNil)
		rootFixture := "../fixtures/congo/python_server"
		files := []string{
			"drones_api.py",
			"deliveries_api.py",
			"app.py",
			"handlers/__init__.py",
			"handlers/drones_postHandler.py",
		}

		Convey("Congo python server", func() {
			server := NewFlaskServer(apiDef, "apidocs", targetDir, true, nil, false, "TestApp")
			err = server.Generate()
			So(err, ShouldBeNil)
			validateFiles(files, targetDir, rootFixture)
			// test that this file exist
			filesExist := []string{
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
			for _, f := range filesExist {
				_, err := os.Stat(filepath.Join(targetDir, f))
				So(err, ShouldBeNil)
			}
		})

		Convey("Congo gevent python server", func() {
			files = append(files, "server.py")
			server := NewFlaskServer(apiDef, "apidocs", targetDir, true, nil, true, "TestApp")
			err = server.Generate()
			So(err, ShouldBeNil)
			validateFiles(files, targetDir, rootFixture)
		})

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})
}

func validateFiles(files []string, targetDir string, rootFixture string) {
	for _, f := range files {
		s, err := utils.TestLoadFile(filepath.Join(targetDir, f))
		So(err, ShouldBeNil)

		tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f))
		So(err, ShouldBeNil)

		So(s, ShouldEqual, tmpl)
	}
}
