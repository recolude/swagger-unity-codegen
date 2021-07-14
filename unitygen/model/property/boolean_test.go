package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_Boolean(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewBoolean("some name")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	classVar := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "some name", name)
	assert.Equal(t, "bool", varType)
	assert.Equal(t, "false", nullVal)
	assert.Equal(t, `	[JsonProperty("some name")]
	public bool SomeName { get; private set; }
`, classVar)
}
