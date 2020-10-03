package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/property"
	"github.com/stretchr/testify/assert"
)

func Test_NumberDefaultsToFloat(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewNumber("someName", "")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "float", varType)
}

func Test_NumberInterpretsInt32(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewNumber("someName", "int32")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "int", varType)
}
