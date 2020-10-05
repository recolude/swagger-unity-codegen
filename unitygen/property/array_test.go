package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/property"
	"github.com/stretchr/testify/assert"
)

func Test_Array(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewArray("someName", property.NewInteger("anything", ""))

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "int[]", varType)
	assert.Equal(t, "null", nullVal)
}
