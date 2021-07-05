package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_Integer(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewInteger("someName", "")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "int", varType)
	assert.Equal(t, "0", nullVal)
}
