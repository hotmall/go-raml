package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hotmall/go-raml/raml"
)

func TestHasCatchAllInRootRoute(t *testing.T) {
	apiDef := &raml.APIDefinition{}

	err := raml.ParseFile("../fixtures/catch_all_recursive_in_root.raml", apiDef)
	assert.NoError(t, err)

	assert.True(t, HasCatchAllInRootRoute(apiDef))
}
