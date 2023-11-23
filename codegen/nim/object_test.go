package nim

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hotmall/go-raml/raml"
	"github.com/hotmall/go-raml/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateObjectFromRaml(t *testing.T) {
	Convey("generate object", t, func() {
		var apiDef raml.APIDefinition

		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("From raml", func() {
			err = raml.ParseFile("../fixtures/struct/struct.raml", &apiDef)
			So(err, ShouldBeNil)

			err = generateAllObjects(&apiDef, targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/object/"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"EnumCity.nim", "EnumCity.nim"},
				{"EnumEnumCityEnum_homeNum.nim", "EnumEnumCityEnum_homeNum.nim"},
				{"EnumEnumCityEnum_parks.nim", "EnumEnumCityEnum_parks.nim"},
				{"animal.nim", "animal.nim"},
				{"Cage.nim", "Cage.nim"},
				{"Cat.nim", "Cat.nim"},
				{"ArrayOfCats.nim", "ArrayOfCats.nim"},
				{"BidimensionalArrayOfCats.nim", "BidimensionalArrayOfCats.nim"},
				{"EnumString.nim", "EnumString.nim"}, // object is enum type
				{"Catanimal.nim", "Catanimal.nim"},
				{"MultipleInheritance.nim", "MultipleInheritance.nim"},
				{"petshop.nim", "petshop.nim"},
				{"NumberFormat.nim", "NumberFormat.nim"},
			}

			for _, check := range checks {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, check.Result))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, check.Expected))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Convey("From raml with JSON", func() {
			err = raml.ParseFile("../fixtures/struct/json/api.raml", &apiDef)
			So(err, ShouldBeNil)

			err = generateAllObjects(&apiDef, targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/object/json"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"PersonInclude.nim", "PersonInclude.nim"},
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
func TestGenerateObjectMethodBody(t *testing.T) {
	Convey("generate object from method body", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("From raml", func() {
			var apiDef raml.APIDefinition
			err := raml.ParseFile("../fixtures/struct/struct.raml", &apiDef)
			So(err, ShouldBeNil)

			err = generateAllObjects(&apiDef, targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/object/"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"usersPostReqBody.nim", "usersPostReqBody.nim"},
				{"usersidGetRespBody.nim", "usersidGetRespBody.nim"},
			}

			for _, check := range checks {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, check.Result))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, check.Expected))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Convey("From raml with included JSON", func() {
			var apiDef raml.APIDefinition
			err := raml.ParseFile("../fixtures/struct/json/api.raml", &apiDef)
			So(err, ShouldBeNil)

			err = generateAllObjects(&apiDef, targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/object/json"
			checks := []struct {
				Result   string
				Expected string
			}{
				{"personPostReqBody.nim", "personPostReqBody.nim"},
				{"personPostReqBody.nim", "personPostReqBody.nim"},
				{"personGetRespBody.nim", "personGetRespBody.nim"},
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
