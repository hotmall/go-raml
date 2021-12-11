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

func TestGenerateStructFromBody(t *testing.T) {
	Convey("generate struct body from raml", t, func() {
		apiDef := new(raml.APIDefinition)

		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("simple body", func() {
			err = raml.ParseFile("../fixtures/struct/struct.raml", apiDef)
			So(err, ShouldBeNil)

			s := NewServer("gorestful", apiDef, "main", "", "examples.com", false, targetDir, nil)
			err := s.Generate()
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/struct/body"
			typeFiles := []string{
				"UsersIdGetRespBody",
				"UsersPostReqBody",
				"Catanimal",
				"UnionCatanimal",
			}

			for _, f := range typeFiles {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, typeDir, f+".go"))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f+".txt"))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

			apiFiles := []string{
				"users_api",
				"users_api_IdGet",
				"users_api_IdPut",
				"users_api_Post",
			}

			for _, f := range apiFiles {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, serverAPIDir, "users", f+".go"))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f+".txt"))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}
		})

		Convey("builtin type doesn't need validation code", func() {
			err = raml.ParseFile("../fixtures/struct/validation.raml", apiDef)
			So(err, ShouldBeNil)

			s := NewServer("gorestful", apiDef, "main", "", "examples.com", false, targetDir, nil)
			err := s.Generate()
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/struct/validation"
			files := []string{
				"builtin_api",
				"builtin_api_Morecomplextype",
				"builtin_api_Scalartype",
			}

			for _, f := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, serverAPIDir, "builtin", f+".go"))
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
