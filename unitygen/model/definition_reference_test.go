package model_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/stretchr/testify/assert"
)

func TestDefinitionReference(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := model.NewDefinitionReference("#/definitions/SomeEnum")

	// ********************************** ACT *********************************
	varType := ref.ToVariableType()
	name := ref.Name()
	converter := ref.JsonConverter()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "", converter)
	assert.Equal(t, "SomeEnum", varType)
	assert.Equal(t, "#/definitions/SomeEnum", name)
	assert.PanicsWithError(t, "unimplemented", func() {
		ref.ToCSharp()
	})
}
