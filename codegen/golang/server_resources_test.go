package golang

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hotmall/go-raml/raml"
	"github.com/hotmall/go-raml/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestResource(t *testing.T) {
	Convey("resource generator", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)

		Convey("interface of simple resource", func() {
			err := raml.ParseFile("../fixtures/server_resources/deliveries.raml", apiDef)
			So(err, ShouldBeNil)

			gs := NewServer(apiDef, "main", "apidocs", "examples.com/libro", true, targetDir, nil)
			_, err = gs.generateServerResources(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/server_resources/simple"
			files := []string{
				"deliveries_if",
			}
			for _, f := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, f+".go"))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f+".txt"))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Convey("simple resource with one api file per method", func() {
			err := raml.ParseFile("../fixtures/server_resources/deliveries.raml", apiDef)
			So(err, ShouldBeNil)

			gs := NewServer(apiDef, "main", "apidocs", "examples.com/libro", true, targetDir, nil)
			_, err = gs.generateServerResources(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/server_resources/one_file_per_method_simple"
			files := []string{
				"deliveries_api",
				"deliveries_api_Get",
				"deliveries_api_Post",
				"deliveries_api_DeliveryIdDelete",
				"deliveries_api_DeliveryIdPatch",
				"deliveries_api_GetDeliveriesByDeliveryID",
			}
			for _, f := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, serverAPIDir, "deliveries", f+".go"))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f+".txt"))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Convey("big raml to check interface consistency", func() {
			err := raml.ParseFile("../fixtures/server_resources/grid/api.raml", apiDef)
			So(err, ShouldBeNil)

			gs := NewServer(apiDef, "main", "apidocs", "examples.com/libro", true, targetDir, nil)
			_, err = gs.generateServerResources(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/server_resources/grid/"
			files := []string{
				"nodes_if",
			}
			for _, f := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, f+".go"))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f+".txt"))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}
		})

		Convey("resource with request body", func() {
			err := raml.ParseFile("../fixtures/server_resources/usergroups.raml", apiDef)
			So(err, ShouldBeNil)

			gs := NewServer(apiDef, "main", "apidocs", "examples.com/libro", true, targetDir, nil)
			_, err = gs.generateServerResources(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/server_resources/with_req_body"
			files := []string{
				"users_api_CreateUsers",
				"users_api",
			}
			for _, f := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, serverAPIDir, "users", f+".go"))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f+".txt"))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}
		})

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})
}
