package python

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/hotmall/go-raml/raml"
	"github.com/hotmall/go-raml/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClientBasic(t *testing.T) {
	Convey("Python client", t, func() {
		apiDef := new(raml.APIDefinition)
		err := raml.ParseFile("./fixtures/client/client.raml", apiDef)
		So(err, ShouldBeNil)

		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("requests client", func() {
			log.Info("requests client")
			client := NewClient(apiDef, "", false)
			err = client.Generate(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/client/requests_client"
			// cek with generated with fixtures
			files := []string{
				"http_client.py",
				"__init__.py",
				"client_utils.py",
				"users_service.py",
			}

			for _, file := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, file))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, file))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

			// make sure these files are exists
			filesExist := []string{
				"Address.py",
				"City.py",
				"GetUsersReqBody.py",
			}
			for _, f := range filesExist {
				_, err := os.Stat(filepath.Join(targetDir, f))
				So(err, ShouldBeNil)
			}

		})

		Convey("requests gevent client", func() {
			log.Info("requests gevent client")
			client := NewClient(apiDef, "gevent-requests", false)
			err = client.Generate(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/client/requests_client"

			s, err := utils.TestLoadFile(filepath.Join(targetDir, "__init__.py"))
			So(err, ShouldBeNil)

			tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, "gevent_client.py"))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		})

		Convey("requests client with unmarshall response", func() {
			log.Info("requests client with unmarshall_response")
			client := NewClient(apiDef, "", true)
			err = client.Generate(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/client/requests_client/unmarshall_response"
			// cek with generated with fixtures
			files := []string{
				"__init__.py",
				"users_service.py",
				"unmarshall_error.py",
				"unhandled_api_error.py",
			}

			for _, file := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, file))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, file))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}
		})

		Convey("aiohttp client", func() {
			log.Info("aiohttp client")
			client := NewClient(apiDef, clientNameAiohttp, false)
			err = client.Generate(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/client/aiohttp_client"
			// cek with generated with fixtures
			files := []string{
				"http_client.py",
				"__init__.py",
				"client_utils.py",
				"users_service.py",
			}

			for _, f := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, f))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, f))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}
		})

		Convey("aiohttp client with unmarshall response", func() {
			log.Info("aiohttp client with unmarshall response")
			client := NewClient(apiDef, clientNameAiohttp, true)
			err = client.Generate(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/client/aiohttp_client/unmarshall_response"
			// cek with generated with fixtures
			files := []string{
				"__init__.py",
				"users_service.py",
				"unmarshall_error.py",
				"unhandled_api_error.py",
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

func TestClientMultislash(t *testing.T) {
	Convey("Python client with multislash root endpoint", t, func() {
		apiDef := new(raml.APIDefinition)
		err := raml.ParseFile("../fixtures/client_resources/multislash.raml", apiDef)
		So(err, ShouldBeNil)

		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("requests client", func() {
			log.Info("requests client")
			client := NewClient(apiDef, "", false)
			err = client.Generate(targetDir)
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/client/multislash/"
			// cek with generated with fixtures
			files := []string{
				"__init__.py",
				"animalsid_service.py",
			}

			for _, file := range files {
				s, err := utils.TestLoadFile(filepath.Join(targetDir, file))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFile(filepath.Join(rootFixture, file))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})
}
