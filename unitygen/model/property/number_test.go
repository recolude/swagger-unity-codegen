package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_NumberDefaultsToFloat(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewNumber("someName", "")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "float", varType)
	assert.Equal(t, "0f", nullVal)
}

func Test_NumberInterpretsInt32(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewNumber("someName", "int32")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "int", varType)
	assert.Equal(t, "0", nullVal)
}
