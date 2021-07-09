package path_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/stretchr/testify/assert"
)

func Test_Parameter(t *testing.T) {
	// ARRANGE ================================================================
	name := "test"
	paramLoc := path.PathParameterLocation
	required := true
	prop := property.NewBoolean("aaa")

	// ACT ====================================================================
	param := path.NewParameter(paramLoc, name, required, prop)

	// ASSERT =================================================================
	assert.Equal(t, paramLoc, param.Location())
	assert.Equal(t, name, param.Name())
	assert.Equal(t, required, param.Required())
	assert.Equal(t, prop, param.Schema())
}
