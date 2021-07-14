package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_Array(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewArray("my array", property.NewInteger("anything", ""))

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	cSharp := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "my array", name)
	assert.Equal(t, "int[]", varType)
	assert.Equal(t, "null", nullVal)
	assert.Equal(t, `	[JsonProperty("my array")]
	public int[] MyArray { get; private set; }
`, cSharp)
}
