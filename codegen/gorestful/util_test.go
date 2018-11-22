package gorestful

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetImportPath(t *testing.T) {
	Convey("TestSetImportPath", t, func() {
		oriGoPath := os.Getenv("GOPATH")
		Convey("users api", func() {
			fakeGopath := "/gopath"
			os.Setenv("GOPATH", fakeGopath)
			So(setRootImportPath("import.com/a", "target"), ShouldEqual, "import.com/a")
			So(setRootImportPath("", "/gopath/src/johndoe.com/cool"), ShouldEqual, "johndoe.com/cool")
		})

		Reset(func() {
			os.Setenv("GOPATH", oriGoPath)
		})
	})
}
