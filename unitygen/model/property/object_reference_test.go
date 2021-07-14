package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_ObjectReference(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewObjectReference("my-permissions", "#/definitions/v1Permission")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	cSharp := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "my-permissions", name)
	assert.Equal(t, "V1Permission", varType)
	assert.Equal(t, "null", nullVal)
	assert.Equal(t, `	[JsonProperty("my-permissions")]
	public V1Permission MyPermissions { get; private set; }
`, cSharp)
}
