package model_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func TestDefinitionWrapperReference(t *testing.T) {
	// ******************************** ARRANGE *******************************
	refs := []model.Definition{
		model.NewStringEnum("cool cat", []string{"Mortimer"}),
		model.NewObject("some obj", []model.Property{
			property.NewBoolean("is cool"),
		}),
	}
	wrap := model.NewDefinitionWrapper(nil)

	// ***************************** ACT / ASSERT *****************************
	for _, ref := range refs {
		wrap.UpdateDefinition(ref)
		varType := wrap.ToVariableType()
		name := wrap.Name()
		converter := wrap.JsonConverter()
		cSharp := wrap.ToCSharp()

		assert.Equal(t, ref.JsonConverter(), converter)
		assert.Equal(t, ref.ToVariableType(), varType)
		assert.Equal(t, ref.Name(), name)
		assert.Equal(t, ref.ToCSharp(), cSharp)
	}
}
