package golang

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hotmall/go-raml/codegen/resource"
	"github.com/hotmall/go-raml/raml"
	"github.com/hotmall/go-raml/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServerMethodWithSpecialChars(t *testing.T) {
	Convey("TestServerMethodWithSpecialChars", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/special_chars.raml", apiDef)
		So(err, ShouldBeNil)

		gs := NewServer(apiDef, "main", "apidocs", "examples.com/libro", true, targetDir, nil)
		_, err = gs.generateServerResources(targetDir)
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/special_chars/server"
		files := []string{
			filepath.Join(serverAPIDir, "escape_type", "escape_type_api_Post"),
			filepath.Join(serverAPIDir, "uri", "uri_api_Users_idGet"),
			"uri_if",
		}

		for _, f := range files {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, f+".go"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f+".txt"))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

}

func TestCatchAllRoute(t *testing.T) {
	Convey("TestServerMethodWithSpecialChars", t, func() {
		sm := serverMethod{
			method: &method{},
		}
		sm.Endpoint = "/users/" + resource.CatchAllRoute

		So(sm.Route(), ShouldEqual, "/users/"+muxCatchAllRoute)
	})
}

func TestServerMethodWithCatchAllRoute(t *testing.T) {
	Convey("TestServerMethodWithCatchAllRoute", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/catch_all_recursive_url.raml", apiDef)
		So(err, ShouldBeNil)

		gs := NewServer(apiDef, "main", "apidocs", "examples.com/libro", true, targetDir, nil)
		_, err = gs.generateServerResources(targetDir)
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/catch_all_recursive_url/server"
		files := []string{
			"tree_if",
		}

		for _, f := range files {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, f+".go"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f+".txt"))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

}

func TestServerMethodWithCatchAllRouteInRoot(t *testing.T) {
	Convey("TestServerMethodWithCatchAllRouteInRoot", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/catch_all_recursive_in_root.raml", apiDef)
		So(err, ShouldBeNil)

		gs := NewServer(apiDef, "main", "apidocs", "examples.com/libro", true, targetDir, nil)
		err = gs.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/catch_all_recursive_url/server-in-root"
		files := []string{
			"path_if",
			"main",
		}

		for _, f := range files {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, f+".go"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f+".txt"))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

}
