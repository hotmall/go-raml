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

func TestClientBasic(t *testing.T) {
	Convey("Test client", t, func() {
		var apiDef raml.APIDefinition
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("users api", func() {
			err = raml.ParseFile("../fixtures/client_resources/client.raml", &apiDef)
			So(err, ShouldBeNil)

			client, err := NewClient(&apiDef, "theclient", "examples.com/libro", targetDir, nil)
			So(err, ShouldBeNil)

			err = client.Generate()
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/client_resources"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"users_service.go", "users_service.txt"},
				{"client_structapitest.go", "client_structapitest.txt"},
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

func TestClientMultiSlashEndpoint(t *testing.T) {
	Convey("Test client with multislash endpoint", t, func() {
		var apiDef raml.APIDefinition
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("users api", func() {
			err = raml.ParseFile("../fixtures/client_resources/multislash.raml", &apiDef)
			So(err, ShouldBeNil)

			client, err := NewClient(&apiDef, "theclient", "examples.com/libro", targetDir, nil)
			So(err, ShouldBeNil)

			err = client.Generate()
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/client_resources/multislash"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"client_goramldir.go", "client_goramldir.txt"},
				{"animalsid_service.go", "animalsid_service.txt"},
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
