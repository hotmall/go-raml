package gorestful

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hotmall/go-raml/raml"
	"github.com/hotmall/go-raml/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClientMethodWithSpecialChars(t *testing.T) {
	Convey("TestClientMethodWithSpecialChars", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/special_chars.raml", apiDef)
		So(err, ShouldBeNil)

		client, err := NewClient(apiDef, "requests", "theclient", "examples.com/libro", targetDir, nil)
		So(err, ShouldBeNil)

		err = client.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/special_chars/client"
		files := []string{
			"escape_type_service",
			"uri_service",
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

func TestClientMethodCatchAllRecursiveURL(t *testing.T) {
	Convey("TestClientMethodCatchAllRecursiveURL", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		apiDef := new(raml.APIDefinition)
		err = raml.ParseFile("../fixtures/catch_all_recursive_url.raml", apiDef)
		So(err, ShouldBeNil)

		client, err := NewClient(apiDef, "requests", "theclient", "examples.com/libro", targetDir, nil)
		So(err, ShouldBeNil)

		err = client.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/method/catch_all_recursive_url/client"
		files := []string{
			"tree_service",
			"client_the_0_metadata",
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
