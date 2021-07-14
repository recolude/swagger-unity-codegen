package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_Integer(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewInteger("Some name", "")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	cSharp := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "Some name", name)
	assert.Equal(t, "int", varType)
	assert.Equal(t, "0", nullVal)
	assert.Equal(t, `	[JsonProperty("Some name")]
	public int SomeName { get; private set; }
`, cSharp)
}
