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

func TestOauth2Middleware(t *testing.T) {
	Convey("oauth2 middleware", t, func() {

		targetdir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("middleware generation test", func() {
			apiDef := new(raml.APIDefinition)
			err := raml.ParseFile("../fixtures/security/dropbox.raml", apiDef)
			So(err, ShouldBeNil)

			err = generateSecurity(apiDef.SecuritySchemes, targetdir, "main")
			So(err, ShouldBeNil)

			// oauth 2 facebook
			s, err := utils.TestLoadFile(filepath.Join(targetdir, "oauth2_Facebook_middleware.go"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile("../fixtures/security/oauth2_Facebook_middleware.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			// oauth 2 dropbox
			s, err = utils.TestLoadFile(filepath.Join(targetdir, "oauth2_Dropbox_middleware.go"))
			So(err, ShouldBeNil)

			tmpl, err = utils.TestLoadFile("../fixtures/security/oauth2_Dropbox_middleware.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

		})

		Convey("Go routes generation", func() {
			apiDef := new(raml.APIDefinition)
			err := raml.ParseFile("../fixtures/security/dropbox.raml", apiDef)
			So(err, ShouldBeNil)

			gs := NewServer("gorestful", apiDef, "main", "", "examples.com/goraml", true, targetdir, nil)
			_, err = gs.generateServerResources(targetdir)
			So(err, ShouldBeNil)

			// check route
			s, err := utils.TestLoadFile(filepath.Join(targetdir, "deliveries_if.go"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile("../fixtures/security/deliveries_if.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		})

		Convey("With included .raml file", func() {
			apiDef := new(raml.APIDefinition)
			err := raml.ParseFile("../fixtures/security/dropbox_with_include.raml", apiDef)
			So(err, ShouldBeNil)

			err = generateSecurity(apiDef.SecuritySchemes, targetdir, "main")
			So(err, ShouldBeNil)

			// oauth 2 middleware
			s, err := utils.TestLoadFile(filepath.Join(targetdir, "oauth2_DropboxIncluded_middleware.go"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile("../fixtures/security/oauth2_DropboxIncluded_middleware.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

		})

		Reset(func() {
			os.RemoveAll(targetdir)
		})
	})
}
