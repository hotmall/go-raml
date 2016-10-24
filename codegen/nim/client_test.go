package nim

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jumpscale/go-raml/raml"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateClient(t *testing.T) {
	Convey("generate client from raml", t, func() {
		var apiDef raml.APIDefinition
		err := raml.ParseFile("../fixtures/client_resources/client.raml", &apiDef)
		So(err, ShouldBeNil)

		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		client := Client{
			APIDef: &apiDef,
			Dir:    targetDir,
		}
		err = client.Generate()
		So(err, ShouldBeNil)

		rootFixture := "./fixtures/resource/client"
		checks := []struct {
			Result   string
			Expected string
		}{
			{"client.nim", "client.nim"},
			{"Users_service.nim", "Users_service.nim"},
		}

		for _, check := range checks {
			s, err := testLoadFile(filepath.Join(targetDir, check.Result))
			So(err, ShouldBeNil)

			tmpl, err := testLoadFile(filepath.Join(rootFixture, check.Expected))
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)
		}

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})
}
