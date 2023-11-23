package types

import (
	"testing"

	"github.com/hotmall/go-raml/raml"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateStructBodyFromRaml(t *testing.T) {
	Convey("generate struct body from raml", t, func() {
		apiDef := new(raml.APIDefinition)

		Convey("simple raml", func() {
			err := raml.ParseFile("./fixtures/api.raml", apiDef)
			So(err, ShouldBeNil)

			expected := []string{
				"/users:POST:body:0",
				"/users/{id}:GET:body:200",
				"[Cat,animal]",
				"Cat | animal",
				"animal",
				"Cat",
				"EnumCity",
			}
			tts := AllTypes(apiDef, "main")

			So(len(tts), ShouldEqual, len(expected))
			for name := range tts {
				So(expected, ShouldContain, name)
				t.Logf("name=%v\n", name)
			}
		})
	})
}
