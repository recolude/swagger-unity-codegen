package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_Objec(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewObject("someName", model.NewObject("someName", []model.Property{
		property.NewBoolean("mybool"),
	}))

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	csharp := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "SomeName", varType)
	assert.Equal(t, "null", nullVal)
	assert.Equal(t, `	[System.Serializable]
public class SomeName {

	public bool mybool;

}
	public SomeName someName;`, csharp)
}
