package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_ObjectReference(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewObjectReference("someName", "#/definitions/v1Permission")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "V1Permission", varType)
	assert.Equal(t, "null", nullVal)
}
