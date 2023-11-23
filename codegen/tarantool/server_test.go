package tarantool

import (
	"github.com/hotmall/go-raml/raml"
	"github.com/hotmall/go-raml/utils"

	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {
	Convey("server generator", t, func() {
		targetDir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)
		Convey("Tarantool server", func() {
			apiDef := new(raml.APIDefinition)
			err = raml.ParseFile("../fixtures/server/user_api/api.raml", apiDef)
			So(err, ShouldBeNil)

			s := NewServer(apiDef, "apidocs", targetDir)
			s.generateResources()
			err := s.generateMain()
			So(err, ShouldBeNil)
			err = s.generateResources()
			So(err, ShouldBeNil)

			rootFixture := "./fixtures/server/"
			checks := []string{
				"main.lua",
				"users_api.lua",
				"handlers/handlers.lua",
				"handlers/users_get_handler.lua",
				"handlers/users_post_handler.lua",
				"handlers/usersuser_id_get_handler.lua",
				"handlers/usersuser_id_delete_handler.lua",
				"handlers/get_user_address_by_id_handler.lua",
			}
			for _, check := range checks {
				s, err := utils.TestLoadFileRemoveID(filepath.Join(targetDir, check))
				So(err, ShouldBeNil)

				tmpl, err := utils.TestLoadFileRemoveID(filepath.Join(rootFixture, check))
				So(err, ShouldBeNil)

				So(s, ShouldEqual, tmpl)
			}

		})

		Reset(func() {
			os.RemoveAll(targetDir)
		})
	})
}
