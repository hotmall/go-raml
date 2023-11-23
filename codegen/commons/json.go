package commons

import (
	"strings"

	"github.com/hotmall/go-raml/raml"
)

// HasJSONBody checks if this raml.Bodies has JSON body that need to be generated it's struct.
// rules:
//	- not nil application/json
//	- has properties or has tipe in JSON string
func HasJSONBody(body *raml.Bodies) bool {
	return body.ApplicationJSON != nil && (len(body.ApplicationJSON.Properties) > 0 || IsJSONString(body.ApplicationJSON.TypeString()))
}

// IsJSONString returns true if a string is a JSON string
func IsJSONString(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}
