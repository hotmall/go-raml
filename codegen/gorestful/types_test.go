package gorestful

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTypeConversion(t *testing.T) {
	Convey("Test Type Conversion", t, func() {
		globGoramlPkgDir = "goraml"
		Convey("Type conversion", func() {
			So(convertToGoType("string", ""), ShouldEqual, "string")
			So(convertToGoType("number", ""), ShouldEqual, "float64")
			So(convertToGoType("integer", ""), ShouldEqual, "int")
			So(convertToGoType("boolean", ""), ShouldEqual, "bool")
			So(convertToGoType("file", ""), ShouldEqual, "string")
			So(convertToGoType("date-only", ""), ShouldEqual, "goraml.DateOnly")
			So(convertToGoType("time-only", ""), ShouldEqual, "goraml.TimeOnly")
			So(convertToGoType("Object", ""), ShouldEqual, "Object")
			So(convertToGoType("string[]", ""), ShouldEqual, "[]string")
			So(convertToGoType("string[][]", ""), ShouldEqual, "[][]string")
			//So(convertToGoType("string | Person"), ShouldEqual, "interface{}")
			//So(convertToGoType("(string | Person)[]"), ShouldEqual, "[]interface{}")
		})
	})
}
