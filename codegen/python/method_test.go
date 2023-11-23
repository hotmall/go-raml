package python

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hotmall/go-raml/raml"
	"github.com/hotmall/go-raml/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMethod(t *testing.T) {
	Convey("server method with display name", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("resource with request body", func() {
			apiDef := new(raml.APIDefinition)
			err := raml.ParseFile("../fixtures/server_resources/display_name/api.raml", apiDef)
			So(err, ShouldBeNil)

			fs := NewFlaskServer(apiDef, "apidocs", targetDir, true, nil, false, "TestApp")

			err = fs.Generate()
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/method/flask/display_name"
			files := []string{
				"coolness_api.py",
			}

			for _, f := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, f))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}
		})

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})
}

func TestServerMethodWithComplexBody(t *testing.T) {
	Convey("TestServerMethodWithComplexBody", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/body.raml", apiDef)
		So(err, ShouldBeNil)

		fs := NewFlaskServer(apiDef, "apidocs", targetDir, true, nil, false, "TestApp")

		err = fs.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/server/complex_body"
		files := []string{
			"arrays_api.py",
		}

		for _, f := range files {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, f))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})
}

func TestServerMethodWithSpecialChars(t *testing.T) {
	Convey("TestServerMethodWithSpecialChars", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/special_chars.raml", apiDef)
		So(err, ShouldBeNil)

		fs := NewFlaskServer(apiDef, "apidocs", targetDir, true, nil, false, "TestApp")

		err = fs.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/special_chars/server"
		files := []string{
			"uri_api.py",
			"escape_type_api.py",
			"handlers/escape_type_postHandler.py",
			"handlers/__init__.py",
			"handlers/uri_byUsers_id_getHandler.py",
			"handlers/schema/User2_0_schema.json",
		}

		for _, f := range files {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, f))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

}

func TestServerMethodWithCatchAllRecursiveURL(t *testing.T) {
	Convey("Flask ", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/catch_all_recursive_url.raml", apiDef)
		So(err, ShouldBeNil)

		fs := NewFlaskServer(apiDef, "apidocs", targetDir, true, nil, false, "TestApp")

		err = fs.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/catch_all_recursive_url/server/flask"
		files := []string{
			"tree_api.py",
		}

		for _, f := range files {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, f))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

	Convey("Sanic ", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/catch_all_recursive_url.raml", apiDef)
		So(err, ShouldBeNil)

		fs := NewSanicServer(apiDef, "apidocs", targetDir, true, nil)

		err = fs.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/catch_all_recursive_url/server/sanic"
		files := []string{
			"tree_if.py",
		}

		for _, f := range files {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, f))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

}

func TestServerMethodWithInRootCatchAllRecursiveURL(t *testing.T) {
	Convey("Flask ", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/catch_all_recursive_in_root.raml", apiDef)
		So(err, ShouldBeNil)

		fs := NewFlaskServer(apiDef, "apidocs", targetDir, true, nil, false, "TestApp")

		err = fs.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/catch_all_recursive_url/server/flask-in-root"
		files := []string{
			"path_api.py",
			"app.py",
		}

		for _, f := range files {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, f))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

	Convey("Sanic ", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/catch_all_recursive_in_root.raml", apiDef)
		So(err, ShouldBeNil)

		fs := NewSanicServer(apiDef, "apidocs", targetDir, true, nil)

		err = fs.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/catch_all_recursive_url/server/sanic-in-root"
		files := []string{
			"path_if.py",
			"app.py",
		}

		for _, f := range files {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, f))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

}
