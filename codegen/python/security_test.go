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

func TestOauth2Middleware(t *testing.T) {
	Convey("oauth2 middleware", t, func() {

		targetdir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("middleware generation test", func() {
			apiDef := new(raml.APIDefinition)
			err := raml.ParseFile("../fixtures/security/dropbox.raml", apiDef)
			So(err, ShouldBeNil)

			err = generateServerSecurity(apiDef.SecuritySchemes, templates(serverKindFlask), targetdir)
			So(err, ShouldBeNil)

			// oauth 2 in dropbox
			s, err := utils.TestLoadFile(filepath.Join(targetdir, "oauth2_Dropbox.py"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile("./fixtures/security/oauth2_Dropbox.py")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			// oauth 2 facebook
			s, err = utils.TestLoadFile(filepath.Join(targetdir, "oauth2_Facebook.py"))
			So(err, ShouldBeNil)

			tmpl, err = utils.TestLoadFile("./fixtures/security/oauth2_Facebook.py")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		})

		Convey("routes generation", func() {
			apiDef := new(raml.APIDefinition)
			err := raml.ParseFile("../fixtures/security/dropbox.raml", apiDef)
			So(err, ShouldBeNil)

			fs := NewFlaskServer(apiDef, "apidocs", targetdir, true, nil, false, "TestApp")
			err = fs.generateResources(targetdir)
			So(err, ShouldBeNil)

			// check route
			s, err := utils.TestLoadFile(filepath.Join(targetdir, "deliveries_api.py"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile("./fixtures/security/deliveries_api.py")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		})

		Convey("flask security classes", func() {
			apiDef := new(raml.APIDefinition)
			err := raml.ParseFile("./fixtures/client/security/client.raml", apiDef)
			So(err, ShouldBeNil)

			c := NewClient(apiDef, clientNameRequests, false)
			err = c.generateSecurity(targetdir)
			So(err, ShouldBeNil)

			files := []string{
				"oauth2_client_itsyouonline.py",
				"basicauth_client_basic.py",
				"passthrough_client_passthrough.py",
			}
			for _, file := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetdir, file))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join("./fixtures/client/security/flask/", file))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Convey("aiohttp security classes", func() {
			apiDef := new(raml.APIDefinition)
			err := raml.ParseFile("./fixtures/client/security/client.raml", apiDef)
			So(err, ShouldBeNil)

			c := NewClient(apiDef, clientNameAiohttp, false)
			err = c.generateSecurity(targetdir)
			So(err, ShouldBeNil)

			files := []string{
				"oauth2_client_itsyouonline.py",
				"basicauth_client_basic.py",
				"passthrough_client_passthrough.py",
			}
			for _, file := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetdir, file))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join("./fixtures/client/security/aiohttp/", file))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})
		Reset(func() {
			os.RemoveAll(targetdir)
		})
	})
}
