package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/property"
	"github.com/stretchr/testify/assert"
)

func Test_Integer(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewInteger("someName", "")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "int", varType)
}
