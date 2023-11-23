package raml

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResourceTypeInheritance(t *testing.T) {
	apiDef := new(APIDefinition)
	err := ParseFile("./samples/resource_types.raml", apiDef)
	Convey("resource type & traits inheritance", t, func() {
		So(err, ShouldBeNil)

		Convey("checking users", func() {
			r := apiDef.Resources["/Users"]

			So(r.URI, ShouldEqual, "/Users")
			So(r.Description, ShouldEqual, "The collection of Users")

			So(r.Get, ShouldNotBeNil)
			So(r.Get.Description, ShouldEqual, "Get all Users, optionally filtered")
			So(r.Get.DisplayName, ShouldEqual, "ListAllUsers")
			So(r.Get.Responses["200"].Bodies.Type, ShouldEqual, "Users")

			So(r.Post, ShouldNotBeNil)
			So(r.Post.Description, ShouldEqual, "Create a new User")
			So(r.Post.Bodies.ApplicationJSON.Type, ShouldEqual, "User")
			So(r.Post.Responses["200"].Bodies.Type, ShouldEqual, "User")
		})

		Convey("checking queues - optional method", func() {
			r := apiDef.Resources["/queues"]
			So(r, ShouldNotBeNil)

			So(r.Get, ShouldNotBeNil)
			So(r.Get.Description, ShouldEqual, "Get all queues")

			So(r.Post, ShouldBeNil)
		})

		Convey("checking corps - header - resourcePath - request body", func() {
			r := apiDef.Resources["/corps"]
			So(r, ShouldNotBeNil)

			So(r.Post, ShouldNotBeNil)

			props := r.Post.Bodies.ApplicationJSON.Properties
			So(ToProperty("name", props["name"]).Type, ShouldEqual, "string")
			So(ToProperty("age", props["age"]).Type, ShouldEqual, "int")
			So(r.Post.Headers["X-Chargeback"].Required, ShouldBeTrue)

			mem := r.Nested["/{id}"]
			So(mem, ShouldNotBeNil)
			So(mem.Get.Description, ShouldEqual, "get /corps/{id}") // check resourcePath parsing

			// check resourcePathName parsing
			respCode := HTTPCode("200")
			So(mem.Get.Responses, ShouldContainKey, respCode)
			So(mem.Get.Responses[respCode].Bodies.Type, ShouldEqual, "corps")
		})

		Convey("books - query parameters", func() {
			r := apiDef.Resources["/books"]
			So(r, ShouldNotBeNil)

			qps := r.Get.QueryParameters
			So(qps["title"].Description, ShouldEqual, "Return books that have their title matching the given value")
			So(qps["digest_all_fields"].Description, ShouldEqual,
				"If no values match the value given for title, use digest_all_fields instead")

			// collection merging
			// test disabled because of issue: https://github.com/hotmall/go-raml/issues/99
			//So(qps["platform"].Enum, ShouldContain, "mac")
			//So(qps["platform"].Enum, ShouldContain, "unix")
			//So(qps["platform"].Enum, ShouldContain, "win")
		})

		Convey("query parameters traits", func() {
			r := apiDef.Resources["/books"]
			So(r, ShouldNotBeNil)

			So(apiDef.Traits, ShouldContainKey, "paged")
			So(r.Get, ShouldNotBeNil)

			qps := r.Get.QueryParameters
			numPages := qps["numPages"]
			So(numPages.Description, ShouldEqual, "The number of pages to return, not to exceed 10")
			So(numPages.Type, ShouldEqual, "integer")
			So(*numPages.Minimum, ShouldEqual, 1)
			So(numPages.Required, ShouldEqual, true)

			So(qps["access_token"].Description, ShouldEqual, "A valid access_token is required")

		})

		Convey("request body traits", func() {
			r := apiDef.Resources["/servers"]
			So(r, ShouldNotBeNil)

			props := r.Post.Bodies.ApplicationJSON.Properties

			So(props, ShouldContainKey, "name")
			So(props, ShouldContainKey, "address?")
			So(props, ShouldNotContainKey, "location?")
			So(props, ShouldNotContainKey, "location")
		})
		Convey("resource types can use traits", func() {
			So(apiDef.ResourceTypes, ShouldContainKey, "file")

			file := apiDef.ResourceTypes["file"]
			So(file.Put, ShouldNotBeNil)
			So(file.Put.Headers, ShouldContainKey, HTTPHeader("drm-key"))
			So(file.Put.Headers["drm-key"].Required, ShouldBeTrue)
		})
	})
}
