package capnp

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hotmall/go-raml/raml"
	"github.com/hotmall/go-raml/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateCapnpSchema(t *testing.T) {
	Convey("generate capnp schema from raml", t, func() {
		var apiDef raml.APIDefinition
		err := raml.ParseFile("./fixtures/struct.raml", &apiDef)
		So(err, ShouldBeNil)

		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("Schema for Python & Nim", func() {
			err = GenerateCapnp(&apiDef, targetDir, "nim", "")
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/struct/vanilla"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"Animal.capnp", "Animal.capnp"},
				{"Cage.capnp", "Cage.capnp"},
				{"Admin.capnp", "Admin.capnp"},
				{"EnumAdminClearanceLevel.capnp", "EnumAdminClearanceLevel.capnp"},
			}

			for _, check := range checks {
				s, err := utils.TestLoadFileRemoveID(filepath.Join(targetDir, check.Result))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFileRemoveID(filepath.Join(rootFixture, check.Expected))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Convey("Schema for Go", func() {
			err = GenerateCapnp(&apiDef, targetDir, "go", "main")
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/struct/golang"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"Animal.capnp", "Animal.capnp"},
				{"Cage.capnp", "Cage.capnp"},
				{"Admin.capnp", "Admin.capnp"},
				{"EnumAdminClearanceLevel.capnp", "EnumAdminClearanceLevel.capnp"},
			}

			for _, check := range checks {
				s, err := utils.TestLoadFileRemoveID(filepath.Join(targetDir, check.Result))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFileRemoveID(filepath.Join(rootFixture, check.Expected))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})
}
