package model_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/stretchr/testify/assert"
)

func TestObjectReference(t *testing.T) {
	// ******************************** ARRANGE *******************************
	enum := model.NewObjectReference("#/definitions/SomeEnum")

	// ********************************** ACT *********************************
	varType := enum.ToVariableType()
	name := enum.Name()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "SomeEnum", varType)
	assert.Equal(t, "#/definitions/SomeEnum", name)
	assert.PanicsWithError(t, "unimplemented", func() {
		enum.ToCSharp()
	})
}
