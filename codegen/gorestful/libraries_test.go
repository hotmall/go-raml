package gorestful

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jumpscale/go-raml/raml"
	"github.com/Jumpscale/go-raml/utils"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testLibRootURLs = []string{
		"https://raw.githubusercontent.com/Jumpscale/go-raml/master/codegen/fixtures/libraries",
		"https://raw.githubusercontent.com/Jumpscale/go-raml/libraries-in-file/codegen/fixtures/libraries/",
	}
)

func TestLibrary(t *testing.T) {
	Convey("Library usage in server", t, func() {
		var apiDef raml.APIDefinition

		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		err = raml.ParseFile("../fixtures/libraries/api.raml", &apiDef)
		So(err, ShouldBeNil)

		server := NewServer("gorestful", &apiDef, "main", "apidocs", "examples.com/ramlcode", true,
			targetDir, testLibRootURLs)
		err = server.Generate()
		So(err, ShouldBeNil)

		rootFixture := "../fixtures/libraries/go_server"
		checks := []struct {
			Result   string
			Expected string
		}{
			{filepath.Join(typeDir, "Place.go"), "Place.txt"},
			{filepath.Join(serverAPIDir, "dirs", "dirs_api.go"), "dirs_api.txt"},
			{filepath.Join(serverAPIDir, "configs", "configs_api.go"), "configs_api.txt"},
			{filepath.Join(serverAPIDir, "configs", "configs_api_Put.go"), "configs_api_Put.txt"},
			{filepath.Join(serverAPIDir, "configs", "configs_api_Post.go"), "configs_api_Post.txt"},
			{filepath.Join(serverAPIDir, "configs", "configs_api_Get.go"), "configs_api_Get.txt"},
		}

		for _, check := range checks {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, check.Result))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, check.Expected))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

	Convey("Library usage in client", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/libraries/api.raml", apiDef)
		So(err, ShouldBeNil)

		client, err := NewClient(apiDef, "requests", "theclient", "examples.com/theclient", targetDir, testLibRootURLs)
		So(err, ShouldBeNil)

		err = client.Generate()
		So(err, ShouldBeNil)

		rootFixture := "../fixtures/libraries/go_client"
		checks := []struct {
			Result   string
			Expected string
		}{
			{filepath.Join(typeDir, "Place.go"), "Place.txt"},
			{"client_exampleapi.go", "client_exampleapi.txt"},
			{"client_utils.go", "client_utils.txt"},
			{"dirs_service.go", "dirs_service.txt"},
			{"configs_service.go", "configs_service.txt"},
		}

		for _, check := range checks {
			s, err := utils.TestLoadFile(filepath.Join(targetDir, check.Result))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, check.Expected))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

	Convey("raml-examples", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("server", func() {

			var apiDef raml.APIDefinition

			err = raml.ParseFile("../fixtures/raml-examples/libraries/api.raml", &apiDef)
			So(err, ShouldBeNil)

			server := NewServer("gorestful", &apiDef, "main", "apidocs", "examples.com/libro", true, targetDir, nil)
			err = server.Generate()
			So(err, ShouldBeNil)

			rootFixture := "../fixtures/libraries/raml-examples/go_server"
			checks := []struct {
				Result   string
				Expected string
			}{
				{filepath.Join(serverAPIDir, "person", "person_api.go"), "person_api.txt"},
				{filepath.Join("types_lib", typeDir, "Person.go"), "types_lib/Person.txt"},
			}

			for _, check := range checks {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, check.Result))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, check.Expected))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}
		})

		Convey("client", func() {
			var apiDef raml.APIDefinition

			err = raml.ParseFile("../fixtures/raml-examples/libraries/api.raml", &apiDef)
			So(err, ShouldBeNil)

			client, err := NewClient(&apiDef, "requests", "client", "examples.com/libro", targetDir, nil)
			So(err, ShouldBeNil)

			err = client.Generate()
			So(err, ShouldBeNil)

			rootFixture := "../fixtures/libraries/raml-examples/go_client"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"person_service.go", "person_service.txt"},
				{filepath.Join("types_lib", typeDir, "Person.go"), "types_lib/Person.txt"},
			}

			for _, check := range checks {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, check.Result))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, check.Expected))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}
		})

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})

}

func TestAliasLibTypeImportPath(t *testing.T) {
	Convey("TestAliasLibTypeImportPath", t, func() {
		tests := []struct {
			path    string
			aliased string
		}{
			{"a.com/libraries/libname/types", `libname_types "a.com/libraries/libname/types"`},
			{"a.com/libname/types", `libname_types "a.com/libname/types"`},
		}

		for _, test := range tests {
			aliased := aliasLibTypeImportPath(test.path)
			So(aliased, ShouldEqual, test.aliased)
		}
	})
}
